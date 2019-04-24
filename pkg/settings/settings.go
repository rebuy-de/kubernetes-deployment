package settings

import (
	"path"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

type Settings struct {
	Defaults Service  `yaml:"defaults"`
	Services Services `yaml:"services"`
}

func FromBytes(data []byte) (*Settings, error) {
	config := new(Settings)
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Read(client kubernetes.Interface) (*Settings, error) {
	const (
		namespace = "default"
		name      = "kubernetes-deployment"
	)

	cm, err := client.Core().ConfigMaps("default").Get("kubernetes-deployment", meta.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Could read ConfigMap '%s/%s'", namespace, name)
	}

	if len(cm.Data) != 1 {
		return nil, errors.Errorf("ConfigMap '%s/%s' needs to contain exactly one file", namespace, name)
	}

	var content string
	for _, data := range cm.Data {
		content = data
	}

	return FromBytes([]byte(content))
}

func (s *Settings) Service(project string) *Service {
	var (
		parts   = strings.SplitN(project, "/", 2)
		name    = parts[0]
		subpath = ""
	)

	if len(parts) > 1 {
		subpath = parts[1]
	}

	service := s.Services.Get(name)
	if service == nil {
		service = &Service{
			Location: gh.Location{
				Repo: name,
			},
		}
	}

	merged := new(Service)
	merged.Defaults(*service)
	merged.Defaults(s.Defaults)
	merged.Location.Path = path.Join(merged.Location.Path, subpath)
	merged.Clean(s.Defaults)

	return merged
}

func (s *Settings) Clean() {
	log.Debug("cleaning settings file")

	s.Defaults.Defaults(Defaults)

	for i := range s.Services {
		service := new(Service)
		service.Defaults(s.Services[i])
		service.Defaults(s.Defaults)
		service.Clean(s.Defaults)
		s.Services[i] = *service
	}
}
