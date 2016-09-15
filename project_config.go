package main

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"github.com/rebuy-de/kubernetes-deployment/settings"
)

type ProjectConfig struct {
	Services *Services
	Settings *settings.Settings
}

func (c ProjectConfig) String() string {
	data, err := yaml.Marshal(c)
	if err != nil {
		panic(err)
	}

	return string(data)
}

func ReadProjectConfigFrom(filename string) (*ProjectConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config := new(ProjectConfig)
	err = yaml.Unmarshal([]byte(data), config)
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *ProjectConfig) WriteTo(filename string) error {
	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, data, 0644)
}
