package statuschecker

import (
	"regexp"
	"time"

	"github.com/google/go-github/github"
	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/cmdutil"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

type State int

const (
	Empty State = iota
	Success
	Ignored
	Pending
	Failure
	Error
)

type Interceptor struct {
	GitHub  gh.Interface
	Options Options
	Branch  *gh.Branch
}

func New(gitHub gh.Interface, options Options) *Interceptor {
	return &Interceptor{
		GitHub:  gitHub,
		Options: options,
	}
}

func (i *Interceptor) PostFetch(branch *gh.Branch) error {
	i.Branch = branch
	return nil
}

func (i *Interceptor) PreApply([]runtime.Object) error {
	age := time.Since(i.Branch.Date)
	if age < i.Options.InitialDelay {
		log.WithFields(log.Fields{
			"Delay": i.Options.InitialDelay - age,
		}).Debug("the commit is very young; waiting so the build system has time to register")
		time.Sleep(i.Options.InitialDelay - age)
	}

	worst, err := i.getWorstState()
	if err != nil {
		return errors.WithStack(err)
	}

	if worst <= Ignored {
		return nil
	}

	if worst >= Failure {
		log.Warn("aborting deployment, because a build failed")
		cmdutil.Exit(1)
	}

	log.Warn("delaying deployment, because there are pending builds")

	notification := time.NewTicker(i.Options.NotificationInterval)
	defer notification.Stop()
	go func() {
		for _ = range notification.C {
			log.Info("still waiting for pending builds")
		}
	}()

	for {
		worst, err := i.getWorstState()
		if err != nil {
			return errors.WithStack(err)
		}

		if worst <= Ignored {
			notification.Stop()
			log.Info("builds finished, continuing with deployment")
			return nil
		}

		if worst >= Failure {
			log.Warn("aborting deployment, because a build failed")
			cmdutil.Exit(1)
		}

		time.Sleep(i.Options.PullInterval)
	}
}

func (i *Interceptor) getWorstState() (State, error) {
	statuses, err := i.GitHub.GetStatuses(&i.Branch.Location)
	if err != nil {
		return Error, errors.WithStack(err)
	}

	worst := Empty

	for _, status := range statuses {
		state, err := i.getState(status)
		if err != nil {
			return Error, errors.WithStack(err)
		}

		if worst < state {
			worst = state
		}
	}

	return worst, nil
}

func (i *Interceptor) getState(status github.RepoStatus) (State, error) {
	logger := log.WithFields(log.Fields{
		"ID":          *status.ID,
		"URL":         *status.URL,
		"State":       *status.State,
		"TargetURL":   *status.TargetURL,
		"Description": *status.Description,
		"Context":     *status.Context,
	})

	ok, err := regexp.MatchString(i.Options.TargetURLRegex, *status.TargetURL)
	if err != nil {
		return Error, errors.Wrapf(err, "failed to execute regex %v", i.Options.TargetURLRegex)
	}

	if !ok {
		logger.WithFields(log.Fields{
			"Regex": i.Options.TargetURLRegex,
		}).Debugf("ignoring status, since target URL doesn't match regex")
		return Ignored, nil
	}

	ok, err = regexp.MatchString(i.Options.ContextRegex, *status.Context)
	if err != nil {
		return Error, errors.Wrapf(err, "failed to execute regex %v", i.Options.TargetURLRegex)
	}

	if !ok {
		logger.WithFields(log.Fields{
			"Regex": i.Options.ContextRegex,
		}).Debugf("ignoring status, since context doesn't match regex")
		return Ignored, nil
	}

	logger.Debugf("status is '%s'", status.GetState())

	switch status.GetState() {
	case "success":
		return Success, nil
	case "pending":
		return Pending, nil
	case "error":
		return Failure, nil
	case "failure":
		return Failure, nil
	default:
		return Error, errors.Errorf("Got unexpected state '%s'", status.GetState())
	}
}
