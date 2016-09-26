package settings

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type ProjectConfig struct {
	Services *Services
	Settings *Settings
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
		return nil, fmt.Errorf("Could not open '%s':'%v'", filename, err)
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

	if localConfig.Settings.RetrySleep != nil {
		c.Settings.RetrySleep = localConfig.Settings.RetrySleep
	}

	if localConfig.Settings.RetryCount != nil {
		c.Settings.RetryCount = localConfig.Settings.RetryCount
	}

	tempMap := make(map[string]string)

	for _, templateValue := range *c.Settings.TemplateValues {
		tempMap[templateValue.Key] = templateValue.Value
	}

	if localConfig.Settings.TemplateValues != nil {
		for _, templateValue := range *localConfig.Settings.TemplateValues {
			tempMap[templateValue.Key] = templateValue.Value
		}
	}

	c.Settings.TemplateValuesMap = tempMap
}
