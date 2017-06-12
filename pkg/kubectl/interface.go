package kubectl

import "io"

type Interface interface {
	Apply(stdin io.Reader) error
}
