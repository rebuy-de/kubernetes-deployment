package settings

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"
	"github.com/rebuy-de/kubernetes-deployment/pkg/templates"

	yaml "gopkg.in/yaml.v2"
)

var (
	DefaultLocation = gh.Location{
		Ref: "master",
	}
)

type Defaults struct {
	Location       gh.Location      `yaml:",inline"`
	TemplateValues templates.Values `yaml:"template-values"`
}

type Settings struct {
	Default  Defaults `yaml:"default"`
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

func Read(location string, client gh.Client) (*Settings, error) {
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

func ReadFromGitHub(filename string, client gh.Client) (*Settings, error) {
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

func (s *Settings) Clean() {
	s.Default.Location.Path = filepath.Clean(strings.Trim(s.Default.Location.Path, "/")) + "/"

	for i := range s.Services {
		service := &s.Services[i]

		service.Location.Defaults(s.Default.Location)
		service.Location.Defaults(DefaultLocation)

		service.TemplateValues.Defaults(s.Default.TemplateValues)

		service.Location.Path = filepath.Clean(strings.Trim(service.Location.Path, "/")) + "/"

		if strings.TrimSpace(service.Name) == "" {
			nameParts := []string{}
			if service.Location.Owner != s.Default.Location.Owner {
				nameParts = append(nameParts, service.Location.Owner)
			}

			if service.Location.Repo != s.Default.Location.Repo {
				nameParts = append(nameParts, service.Location.Repo)
			}

			if service.Location.Path != s.Default.Location.Path {
				path := service.Location.Path
				path = strings.Trim(path, "/")
				path = strings.Replace(path, "/", "-", -1)
				nameParts = append(nameParts, path)
			}

			service.Name = strings.Join(nameParts, "-")
		}
	}
}
