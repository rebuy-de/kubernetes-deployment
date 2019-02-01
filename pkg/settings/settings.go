package settings

import (
	"io/ioutil"
	"path"
	"regexp"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
)

// This matches the value of `.metadata.selfLink` if ConfigMaps. It is a bit
// weird for identifying that a ConfigMap is requested, but it is good enough
// for now.
var reKubeSelfLink = regexp.MustCompile(`^/api/v1/namespaces/([^/]+)/configmaps/([^/]+)$`)
var configMapFilename = `settings.yaml`

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

func Read(location string, ghClient gh.Interface, kubeClient kubernetes.Interface) (*Settings, error) {
	if strings.HasPrefix(location, "github.com") {
		return ReadFromGitHub(location, ghClient)
	} else if reKubeSelfLink.MatchString(location) {
		matches := reKubeSelfLink.FindStringSubmatch(location)
		return ReadFromConfigMap(matches[1], matches[2], kubeClient)
	} else {
		return ReadFromFile(location)
	}
}

func ReadFromFile(filename string) (*Settings, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not open file '%s'", filename)
	}

	return FromBytes(data)

}

func ReadFromConfigMap(namespace, name string, client kubernetes.Interface) (*Settings, error) {
	cm, err := client.Core().ConfigMaps(namespace).Get(name, meta.GetOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "Could read ConfigMap '%s/%s'", namespace, name)
	}

	content, exists := cm.Data[configMapFilename]

	if !exists {
		return nil, errors.Errorf("Configmap '%s/%s' does not have a file named '%s'",
			namespace, name, configMapFilename)
	}

	return FromBytes([]byte(content))
}

func ReadFromGitHub(filename string, client gh.Interface) (*Settings, error) {
	location, err := gh.NewLocation(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "parse GitHub location '%s'; use './' prefix to use a directory named 'github.com'", filename)
	}

	location.Defaults(gh.Location{
		Ref: "master",
	})

	file, err := client.GetFile(location)
	if err != nil {
		return nil, errors.Wrapf(err, "could not download file '%s'", location)
	}
	return FromBytes([]byte(file.Content))
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
