package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"github.com/awesomeProject/apidiff/internal/model"
	"gopkg.in/yaml.v3"
)

func LoadFile(path string) (model.Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return model.Config{}, err
	}

	ext := filepath.Ext(path)

	var cfg model.Config
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &cfg)
	case ".json":
		err = json.Unmarshal(data, &cfg)
	default:
		err = errors.New("unsupported file type; use .yaml, .yml, or .json")
	}

	if err != nil {
		return model.Config{}, err
	}

	cfg.Normalize()
	return cfg, nil
}
