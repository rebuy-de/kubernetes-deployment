package settings

import (
	"io/ioutil"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Contexts map[string]Service

type Settings struct {
	Defaults Service  `yaml:"defaults"`
	Services Services `yaml:"services"`
	Contexts Contexts `yaml:"contexts"`

	context string
}

func FromBytes(data []byte) (*Settings, error) {
	config := new(Settings)
	err := yaml.Unmarshal(data, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func Read(location string, client gh.Interface) (*Settings, error) {
	if strings.HasPrefix(location, "github.com") {
		return ReadFromGitHub(location, client)
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

func (s *Settings) CurrentContext() Service {
	return s.Contexts[s.context]
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
	merged.Defaults(s.CurrentContext())
	merged.Location.Path = path.Join(merged.Location.Path, subpath)
	merged.Clean(s.CurrentContext())

	return merged
}

func (s *Settings) Clean(contextName string) {
	s.context = contextName
	if s.context == "" {
		s.context = s.Defaults.Context
	}

	log.WithFields(log.Fields{
		"Context": s.context,
	}).Debug("cleaning settings file")

	s.Defaults.Defaults(Defaults)

	for name := range s.Contexts {
		context := s.Contexts[name]
		context.Defaults(s.Defaults)
		context.Context = name
		s.Contexts[name] = context
	}

	for i := range s.Services {
		service := new(Service)
		service.Defaults(s.Services[i])
		service.Defaults(s.CurrentContext())
		service.Clean(s.CurrentContext())
		s.Services[i] = *service
	}
}
