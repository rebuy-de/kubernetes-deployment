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

func (c *ProjectConfig) MergeConfig(localConfig *ProjectConfig) {

	if localConfig.Settings.Kubeconfig != nil {
		c.Settings.Kubeconfig = localConfig.Settings.Kubeconfig
	}

	if localConfig.Settings.Output != nil {
		c.Settings.Output = localConfig.Settings.Output
	}

	if localConfig.Settings.Sleep != nil {
		c.Settings.Sleep = localConfig.Settings.Sleep
	}

	if localConfig.Settings.SkipShuffle != nil {
		c.Settings.SkipShuffle = localConfig.Settings.SkipShuffle
	}

	if localConfig.Settings.SkipFetch != nil {
		c.Settings.SkipFetch = localConfig.Settings.SkipFetch
	}

	if localConfig.Settings.SkipDeploy != nil {
		c.Settings.SkipDeploy = localConfig.Settings.SkipDeploy
	}

	if localConfig.Settings.RetrySleep != nil {
		c.Settings.RetrySleep = localConfig.Settings.RetrySleep
	}

	if localConfig.Settings.RetryCount != nil {
		c.Settings.RetryCount = localConfig.Settings.RetryCount
	}

	if localConfig.Settings.IgnoreDeployFailures != nil {
		c.Settings.IgnoreDeployFailures = localConfig.Settings.IgnoreDeployFailures
	}

	tempMap := make(map[string]string)

	for _, templateValue := range *c.Settings.TemplateValues {
		tempMap[templateValue.Key] = templateValue.Value
	}

	for _, templateValue := range *localConfig.Settings.TemplateValues {
		tempMap[templateValue.Key] = templateValue.Value
	}

	c.Settings.TemplateValuesMap = tempMap
}
