package settings

import (
	"io/ioutil"
	"strings"

	"github.com/pkg/errors"
	"github.com/rebuy-de/kubernetes-deployment/pkg/gh"

	yaml "gopkg.in/yaml.v2"
)

type Settings struct {
	TemplateValues TemplateValues `yaml:"template-values"`
	Services       Services       `yaml:"services"`
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

func ReadFromGitHub(location string, client gh.Client) (*Settings, error) {
	data, err := client.GetContents(location)
	if err != nil {
		return nil, errors.Wrapf(err, "could not download file '%s'", location)
	}
	return FromBytes([]byte(data))
}
