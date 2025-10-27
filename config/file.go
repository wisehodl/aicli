package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

func loadConfigFile(path string) (fileValues, error) {
	if path == "" {
		return fileValues{}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return fileValues{}, err
	}

	var raw map[string]interface{}
	if err := yaml.Unmarshal(data, &raw); err != nil {
		return fileValues{}, err
	}

	fv := fileValues{}
	if v, ok := raw["protocol"].(string); ok {
		fv.protocol = v
	}
	if v, ok := raw["url"].(string); ok {
		fv.url = v
	}
	if v, ok := raw["key_file"].(string); ok {
		fv.keyFile = v
	}
	if v, ok := raw["model"].(string); ok {
		fv.model = v
	}
	if v, ok := raw["fallback"].(string); ok {
		fv.fallback = v
	}
	if v, ok := raw["system_file"].(string); ok {
		fv.systemFile = v
	}

	return fv, nil
}
