package validator

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type RuleSet struct {
	Conflicts []ConflictRule `json:"conflicts" yaml:"conflicts"`
	Requires  []RequireRule  `json:"requires" yaml:"requires"`
}

type ConflictRule struct {
	Name    string   `json:"name" yaml:"name"`
	Scope   []string `json:"scope" yaml:"scope"`
	Plugins []string `json:"plugins" yaml:"plugins"`
}

type RequireRule struct {
	Name   string   `json:"name" yaml:"name"`
	Scope  []string `json:"scope" yaml:"scope"`
	Plugin string   `json:"plugin" yaml:"plugin"`
	Fields []string `json:"fields" yaml:"fields"`
}

func DefaultRules() RuleSet {
	return RuleSet{
		Conflicts: []ConflictRule{
			{
				Name:    "limit-req conflicts with limit-count",
				Scope:   []string{"route", "service", "plugin_config"},
				Plugins: []string{"limit-req", "limit-count"},
			},
		},
		Requires: []RequireRule{
			{
				Name:   "key-auth requires key",
				Scope:  []string{"consumer"},
				Plugin: "key-auth",
				Fields: []string{"key"},
			},
		},
	}
}

func LoadRules(path string) (RuleSet, error) {
	if path == "" {
		return DefaultRules(), nil
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return RuleSet{}, err
	}

	ext := filepath.Ext(path)
	var rs RuleSet
	switch ext {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(data, &rs)
	case ".json":
		err = json.Unmarshal(data, &rs)
	default:
		err = errors.New("unsupported rules file type; use .yaml, .yml, or .json")
	}
	if err != nil {
		return RuleSet{}, err
	}

	return rs, nil
}
