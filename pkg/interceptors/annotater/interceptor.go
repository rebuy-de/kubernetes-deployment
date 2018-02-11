package annotater

import (
	"reflect"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	log "github.com/sirupsen/logrus"
	v1beta1apps "k8s.io/api/apps/v1beta1"
	v1beta1extensions "k8s.io/api/extensions/v1beta1"
	v1meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

type Interceptor struct {
	branch gh.Branch
}

func New() *Interceptor {
	return &Interceptor{}
}

func (i *Interceptor) PostFetch(branch *gh.Branch) error {
	i.branch = *branch
	return nil
}

func (i *Interceptor) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	switch typed := obj.(type) {
	case *v1beta1extensions.Deployment:
		i.AddAnnotations(&typed.ObjectMeta)

	case *v1beta1apps.StatefulSet:
		i.AddAnnotations(&typed.ObjectMeta)

	default:
		log.WithFields(log.Fields{
			"type": reflect.TypeOf(obj),
		}).Debug("type doesn't support adding of annotations")
	}
	return obj, nil
}

func (i *Interceptor) AddAnnotations(meta *v1meta.ObjectMeta) {
	meta.SetAnnotations(map[string]string{
		"rebuy.com/kubernetes-deployment.commit-sha":      i.branch.SHA,
		"rebuy.com/kubernetes-deployment.commit-date":     i.branch.Date.String(),
		"rebuy.com/kubernetes-deployment.commit-author":   i.branch.Author,
		"rebuy.com/kubernetes-deployment.commit-message":  i.branch.Message,
		"rebuy.com/kubernetes-deployment.commit-location": i.branch.Location.String(),
	})
}
