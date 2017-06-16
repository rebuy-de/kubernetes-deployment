package fake

import yaml "gopkg.in/yaml.v2"

func YAML(obj interface{}) string {
	raw, err := yaml.Marshal(&obj)
	if err != nil {
		panic(err)
	}

	return string(raw)
}
