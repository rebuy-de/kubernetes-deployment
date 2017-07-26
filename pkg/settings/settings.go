package settings

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

type Contexts map[string]Service

type Settings struct {
	Defaults Service  `yaml:"defaults"`
	Services Services `yaml:"services"`
	Contexts Contexts `yaml:"contexts"`
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

	data, err := client.GetFile(location)
	if err != nil {
		return nil, errors.Wrapf(err, "could not download file '%s'", location)
	}
	return FromBytes([]byte(data))
}

func (s *Settings) Service(name string) *Service {
	for _, service := range s.Services {
		if service.Name == name {
			return &service
		}
	}

	return nil
}

func (s *Settings) Clean(contextName string) {
	if contextName == "" {
		contextName = s.Defaults.Context
	}

	log.WithFields(log.Fields{
		"Context": contextName,
	}).Debug("cleaning settings file")

	s.Defaults.Defaults(Defaults)

	for name := range s.Contexts {
		context := s.Contexts[name]
		context.Defaults(s.Defaults)
		s.Contexts[name] = context
	}

	context := s.Contexts[contextName]

	for i := range s.Services {
		service := &s.Services[i]
		service.Defaults(context)

		if strings.TrimSpace(service.Name) == "" {
			nameParts := []string{}
			if service.Location.Owner != s.Defaults.Location.Owner {
				nameParts = append(nameParts, service.Location.Owner)
			}

			if service.Location.Repo != s.Defaults.Location.Repo {
				nameParts = append(nameParts, service.Location.Repo)
			}

			if service.Location.Path != s.Defaults.Location.Path {
				path := service.Location.Path
				path = strings.Trim(path, "/")
				path = strings.Replace(path, "/", "-", -1)
				nameParts = append(nameParts, path)
			}

			service.Name = strings.Join(nameParts, "-")
		}
	}
}
