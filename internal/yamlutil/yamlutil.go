package yamlutil

import (
	"github.com/goccy/go-yaml"
)

func ObjectToYaml(object any) ([]byte, error) {
	bytes, err := yaml.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
