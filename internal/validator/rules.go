// Package validator performs semantic checks on configs.
// It validates upstream reachability and plugin rule sets.
package validator

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// RuleSet defines the plugin validation rules loaded from disk.
type RuleSet struct {
	Conflicts  []ConflictRule   `json:"conflicts" yaml:"conflicts"`
	Requires   []RequireRule    `json:"requires" yaml:"requires"`
	RequireAny []RequireAnyRule `json:"require_one_of" yaml:"require_one_of"`
	DenyFields []DenyFieldRule  `json:"deny_fields" yaml:"deny_fields"`
	RegexRules []RegexRule      `json:"regex" yaml:"regex"`
}

// ConflictRule declares plugins that cannot be enabled together.
type ConflictRule struct {
	Name    string   `json:"name" yaml:"name"`
	Scope   []string `json:"scope" yaml:"scope"`
	Plugins []string `json:"plugins" yaml:"plugins"`
}

// RequireRule declares fields that must be present when a plugin is enabled.
type RequireRule struct {
	Name   string   `json:"name" yaml:"name"`
	Scope  []string `json:"scope" yaml:"scope"`
	Plugin string   `json:"plugin" yaml:"plugin"`
	Fields []string `json:"fields" yaml:"fields"`
}

// RequireAnyRule declares that at least one field is required for a plugin.
type RequireAnyRule struct {
	Name   string   `json:"name" yaml:"name"`
	Scope  []string `json:"scope" yaml:"scope"`
	Plugin string   `json:"plugin" yaml:"plugin"`
	Fields []string `json:"fields" yaml:"fields"`
}

// DenyFieldRule declares fields that must not appear for a plugin.
type DenyFieldRule struct {
	Name   string   `json:"name" yaml:"name"`
	Scope  []string `json:"scope" yaml:"scope"`
	Plugin string   `json:"plugin" yaml:"plugin"`
	Fields []string `json:"fields" yaml:"fields"`
}

// RegexRule declares a regex constraint for a plugin field.
type RegexRule struct {
	Name    string   `json:"name" yaml:"name"`
	Scope   []string `json:"scope" yaml:"scope"`
	Plugin  string   `json:"plugin" yaml:"plugin"`
	Field   string   `json:"field" yaml:"field"`
	Pattern string   `json:"pattern" yaml:"pattern"`
}

// DefaultRules returns a minimal built-in rule set.
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

// LoadRules loads a rules file from YAML or JSON.
// An empty path returns the default rule set.
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
