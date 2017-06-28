package interceptors

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
)

type Multi struct {
	Interceptors []interface{}
}

func New(interceptors ...interface{}) *Multi {
	return &Multi{
		Interceptors: interceptors,
	}
}

func (m *Multi) ManifestApplied(obj runtime.Object) error {
	for _, i := range m.Interceptors {
		c, ok := i.(ManifestApplied)
		if !ok {
			continue
		}

		err := c.ManifestApplied(obj)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *Multi) AllManifestsApplied(objs []runtime.Object) error {
	for _, i := range m.Interceptors {
		c, ok := i.(AllManifestsApplied)
		if !ok {
			continue
		}

		err := c.AllManifestsApplied(objs)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *Multi) Close() error {
	var err error

	for _, i := range m.Interceptors {
		c, ok := i.(Closer)
		if !ok {
			continue
		}

		err = c.Close()
		if err != nil {
			log.Warn(err)
		}
	}

	return err
}
