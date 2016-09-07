package kubernetes

type API interface {
	Apply(manifestFile string) ([]byte, error)
	Get(manifestFile string) ([]byte, error)
}
