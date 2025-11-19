package yamlutil

import (
	"github.com/goccy/go-yaml"
	"github.com/hanle23/shorty/internal/types"
)

func ObjectToYaml(object types.Shortcut) ([]byte, error) {
	bytes, err := yaml.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
