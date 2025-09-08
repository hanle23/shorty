package yamlutil

import (
	"github.com/goccy/go-yaml"
	"github.com/hanle23/shorty/internal/config"
)

func ObjectToYaml(object config.ShortcutFile) ([]byte, error) {
	bytes, err := yaml.Marshal(object)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}
