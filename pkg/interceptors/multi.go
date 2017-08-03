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

func (m *Multi) Add(interceptors ...interface{}) {
	m.Interceptors = append(m.Interceptors, interceptors...)
}

func (m *Multi) PreManifestApply(obj runtime.Object) error {
	for _, i := range m.Interceptors {
		c, ok := i.(PreManifestApplier)
		if !ok {
			continue
		}

		err := c.PreManifestApply(obj)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *Multi) PostApply(objs []runtime.Object) error {
	for _, i := range m.Interceptors {
		c, ok := i.(PostApplier)
		if !ok {
			continue
		}

		err := c.PostApply(objs)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (m *Multi) PostManifestRender(obj runtime.Object) (runtime.Object, error) {
	var err error

	for _, i := range m.Interceptors {
		c, ok := i.(PostManifestRenderer)
		if !ok {
			continue
		}

		obj, err = c.PostManifestRender(obj)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return obj, nil
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
