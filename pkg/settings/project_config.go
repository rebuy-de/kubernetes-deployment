package settings

import (
	"fmt"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v2"
)

type ProjectConfig struct {
	Services Services `json:"services"`
	Settings Settings `json:"settings"`
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

func (c *ProjectConfig) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			data, err := yaml.Marshal(c)
			if err == nil {
				io.WriteString(s, string(data))
				return
			}
		}
		fallthrough
	case 's', 'q':
		fmt.Fprintf(s, "foobar")
	}
}
