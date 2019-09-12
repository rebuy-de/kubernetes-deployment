package imagechecker

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ecr"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/rebuy-de/kubernetes-deployment/pkg/kubeutil"
)

type set map[string]*struct{}

func (s set) add(value string) { s[value] = nil }

type Interceptor struct {
	session *session.Session
	Options Options
}

func New(session *session.Session, options Options) *Interceptor {
	return &Interceptor{
		session: session,
		Options: options,
	}
}

func (i *Interceptor) PreApply(objects []runtime.Object) error {
	images := make(set)

	for _, obj := range objects {
		template := kubeutil.PodTemplateSpecAccessor(obj)
		if template == nil {
			continue
		}

		for _, container := range template.Spec.Containers {
			images.add(container.Image)
		}

		for _, container := range template.Spec.InitContainers {
			images.add(container.Image)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), i.Options.WaitTimeout)
	defer cancel()
	go notifier(ctx)

	for image := range images {
		err := i.wait(ctx, image)
		if err != nil {
			return err
		}
	}

	return nil
}

func notifier(ctx context.Context) {
	firstTry := true
	t := time.NewTicker(1 * time.Minute)
	defer t.Stop()

	for {
		select {
		case <-t.C:
			if firstTry {
				firstTry = false
				continue
			}
			logrus.Info("Deployment is waiting for image to be built.")
		case <-ctx.Done():
			return
		}
	}
}

func (i *Interceptor) wait(ctx context.Context, image string) error {
	t := time.NewTicker(i.Options.CheckInterval)
	defer t.Stop()

	logrus.WithFields(logrus.Fields{
		"image": image,
	}).Debugf("Checking for image availability.")

	for {
		select {
		case <-t.C:
			found, err := i.check(ctx, image)
			if err != nil {
				return errors.Wrapf(err, "failed to check image %s", image)
			}
			if found {
				return nil
			}
		case <-ctx.Done():
			return fmt.Errorf("Deployment aborted because image %s is still missing", image)
		}
	}
}

func (i *Interceptor) check(ctx context.Context, image string) (bool, error) {
	reECR, err := regexp.Compile(`^([0-9]+).dkr.ecr.([a-z0-9-]+).amazonaws.com/([^:]+):(.*)$`)
	if err != nil {
		return false, errors.WithStack(err)
	}

	m := reECR.FindStringSubmatch(image)
	if m != nil {
		return i.checkECR(ctx, m[1], m[2], m[3], m[4])
	}

	logrus.WithField("image", image).Debugf("Unknown image provider. Assuming it exists.")
	return true, nil

}

func (i *Interceptor) checkECR(ctx context.Context, account, region, repo, tag string) (bool, error) {
	l := logrus.WithFields(logrus.Fields{
		"account": account,
		"region":  region,
		"repo":    repo,
		"tag":     tag,
	})
	l.Debugf("Checking for image on ECR.")

	sess := i.session.Copy(&aws.Config{
		Region: aws.String(region),
	})
	creds := stscreds.NewCredentials(sess, fmt.Sprintf("arn:aws:iam::%s:role/ecr-cross-access", account))

	svc := ecr.New(sess, &aws.Config{Credentials: creds})
	images, err := svc.DescribeImagesWithContext(ctx, &ecr.DescribeImagesInput{
		RegistryId:     aws.String(account),
		RepositoryName: aws.String(repo),
		ImageIds: []*ecr.ImageIdentifier{
			&ecr.ImageIdentifier{
				ImageTag: aws.String(tag),
			},
		},
	})

	if err == nil {
		image := images.ImageDetails[0]
		l.WithFields(logrus.Fields{
			"digest":    aws.StringValue(image.ImageDigest),
			"size":      aws.Int64Value(image.ImageSizeInBytes),
			"tags":      aws.StringValueSlice(image.ImageTags),
			"pushed_at": aws.TimeValue(image.ImagePushedAt),
		}).Debugf("Found image.")
		return true, nil
	}

	aerr, ok := err.(awserr.Error)
	if !ok {
		l.Debugf("Unkown error: %v", err)
		return false, errors.WithStack(err)
	}

	switch aerr.Code() {
	case ecr.ErrCodeImageNotFoundException:
		// This means the repo was found, but the image tag wasn't.
		l.Debugf("Not available, yet.")
		return false, nil
	default:
		l.Debugf("Unexpected error: %v", err)
		return false, errors.WithStack(err)
	}
}
