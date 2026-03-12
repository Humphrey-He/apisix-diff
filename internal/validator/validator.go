package validator

import (
	"context"
	"fmt"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/awesomeProject/apidiff/internal/model"
)

// Options configures validation behavior.
type Options struct {
	// SkipReachability disables TCP checks against upstream nodes.
	SkipReachability bool
	// RulesPath points to an optional YAML/JSON rule set file.
	RulesPath string
}

// ValidateConfig validates a config snapshot against semantic rules.
// It performs network reachability checks unless disabled.
func ValidateConfig(ctx context.Context, cfg model.Config, opts Options) error {
	if !opts.SkipReachability {
		if err := validateUpstreamReachability(ctx, cfg.Upstreams); err != nil {
			return err
		}
	}

	rules, err := LoadRules(opts.RulesPath)
	if err != nil {
		return err
	}

	if err := validatePluginsWithRules(cfg, rules); err != nil {
		return err
	}

	return nil
}

func validateUpstreamReachability(ctx context.Context, upstreams []model.Upstream) error {
	for _, u := range upstreams {
		for addr := range u.Nodes {
			if err := dialAddress(ctx, addr); err != nil {
				return fmt.Errorf("upstream %s node %s unreachable: %w", u.Key(), addr, err)
			}
		}
	}
	return nil
}

func dialAddress(ctx context.Context, addr string) error {
	if !strings.Contains(addr, ":") {
		return fmt.Errorf("invalid address %s", addr)
	}

	d := net.Dialer{Timeout: 2 * time.Second}
	conn, err := d.DialContext(ctx, "tcp", addr)
	if err != nil {
		return err
	}
	return conn.Close()
}

func validatePluginsWithRules(cfg model.Config, rules RuleSet) error {
	for _, r := range cfg.Routes {
		if err := validatePluginRules("route", r.Key(), r.Plugins, rules); err != nil {
			return err
		}
	}
	for _, s := range cfg.Services {
		if err := validatePluginRules("service", s.Key(), s.Plugins, rules); err != nil {
			return err
		}
	}
	for _, pc := range cfg.PluginConfigs {
		if err := validatePluginRules("plugin_config", pc.Key(), pc.Plugins, rules); err != nil {
			return err
		}
	}
	for _, c := range cfg.Consumers {
		if err := validatePluginRules("consumer", c.Key(), c.Plugins, rules); err != nil {
			return err
		}
	}
	return nil
}

func validatePluginRules(scope, key string, plugins map[string]any, rules RuleSet) error {
	if plugins == nil {
		return nil
	}
	for _, rule := range rules.Conflicts {
		if !scopeMatches(rule.Scope, scope) {
			continue
		}
		if hasAllPlugins(plugins, rule.Plugins) {
			name := rule.Name
			if name == "" {
				name = strings.Join(rule.Plugins, " + ")
			}
			return fmt.Errorf("%s %s violates conflict rule: %s", scope, key, name)
		}
	}

	for _, rule := range rules.Requires {
		if !scopeMatches(rule.Scope, scope) {
			continue
		}
		if _, ok := plugins[rule.Plugin]; !ok {
			continue
		}
		for _, field := range rule.Fields {
			if !pluginFieldExists(plugins, rule.Plugin, field) {
				name := rule.Name
				if name == "" {
					name = fmt.Sprintf("%s requires %s", rule.Plugin, field)
				}
				return fmt.Errorf("%s %s violates require rule: %s", scope, key, name)
			}
		}
	}

	for _, rule := range rules.RequireAny {
		if !scopeMatches(rule.Scope, scope) {
			continue
		}
		if _, ok := plugins[rule.Plugin]; !ok {
			continue
		}
		if !pluginHasAnyField(plugins, rule.Plugin, rule.Fields) {
			name := rule.Name
			if name == "" {
				name = fmt.Sprintf("%s requires one of %s", rule.Plugin, strings.Join(rule.Fields, ", "))
			}
			return fmt.Errorf("%s %s violates require_one_of rule: %s", scope, key, name)
		}
	}

	for _, rule := range rules.DenyFields {
		if !scopeMatches(rule.Scope, scope) {
			continue
		}
		if _, ok := plugins[rule.Plugin]; !ok {
			continue
		}
		for _, field := range rule.Fields {
			if pluginFieldExists(plugins, rule.Plugin, field) {
				name := rule.Name
				if name == "" {
					name = fmt.Sprintf("%s forbids %s", rule.Plugin, field)
				}
				return fmt.Errorf("%s %s violates deny_fields rule: %s", scope, key, name)
			}
		}
	}

	for _, rule := range rules.RegexRules {
		if !scopeMatches(rule.Scope, scope) {
			continue
		}
		if rule.Plugin != "" {
			if _, ok := plugins[rule.Plugin]; !ok {
				continue
			}
		}
		if rule.Field == "" || rule.Pattern == "" {
			continue
		}
		value, ok := pluginFieldValue(plugins, rule.Plugin, rule.Field)
		if !ok {
			continue
		}
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("%s %s violates regex rule: %s (field not string)", scope, key, ruleName(rule, "regex"))
		}
		re, err := regexp.Compile(rule.Pattern)
		if err != nil {
			return fmt.Errorf("invalid regex pattern for rule %s: %w", ruleName(rule, "regex"), err)
		}
		if !re.MatchString(str) {
			return fmt.Errorf("%s %s violates regex rule: %s", scope, key, ruleName(rule, "regex"))
		}
	}

	return nil
}

func ruleName(rule RegexRule, fallback string) string {
	if rule.Name != "" {
		return rule.Name
	}
	if rule.Plugin != "" {
		return fmt.Sprintf("%s.%s matches %s", rule.Plugin, rule.Field, rule.Pattern)
	}
	return fallback
}

func scopeMatches(scopes []string, scope string) bool {
	if len(scopes) == 0 {
		return true
	}
	for _, s := range scopes {
		if s == scope {
			return true
		}
	}
	return false
}

func hasAllPlugins(plugins map[string]any, names []string) bool {
	if len(names) == 0 {
		return false
	}
	for _, name := range names {
		if _, ok := plugins[name]; !ok {
			return false
		}
	}
	return true
}

func pluginHasAnyField(plugins map[string]any, pluginName string, fields []string) bool {
	for _, field := range fields {
		if pluginFieldExists(plugins, pluginName, field) {
			return true
		}
	}
	return false
}

func pluginFieldValue(plugins map[string]any, pluginName, field string) (any, bool) {
	if pluginName != "" {
		raw, ok := plugins[pluginName]
		if !ok {
			return nil, false
		}
		current, ok := raw.(map[string]any)
		if !ok {
			return nil, false
		}
		return nestedFieldValue(current, field)
	}
	if field == "" {
		return nil, false
	}

	parts := strings.Split(field, ".")
	if len(parts) < 2 {
		return nil, false
	}
	plugin := parts[0]
	return pluginFieldValue(plugins, plugin, strings.Join(parts[1:], "."))
}

func nestedFieldValue(current map[string]any, field string) (any, bool) {
	parts := strings.Split(field, ".")
	for i, p := range parts {
		val, ok := current[p]
		if !ok {
			return nil, false
		}
		if i == len(parts)-1 {
			return val, true
		}
		child, ok := val.(map[string]any)
		if !ok {
			return nil, false
		}
		current = child
	}
	return nil, false
}

func pluginFieldExists(plugins map[string]any, pluginName, field string) bool {
	raw, ok := plugins[pluginName]
	if !ok {
		return false
	}
	current, ok := raw.(map[string]any)
	if !ok {
		return false
	}
	parts := strings.Split(field, ".")
	for i, p := range parts {
		val, ok := current[p]
		if !ok {
			return false
		}
		if i == len(parts)-1 {
			return !isEmptyValue(val)
		}
		child, ok := val.(map[string]any)
		if !ok {
			return false
		}
		current = child
	}
	return false
}

func isEmptyValue(v any) bool {
	if v == nil {
		return true
	}
	switch t := v.(type) {
	case string:
		return t == ""
	case []any:
		return len(t) == 0
	case map[string]any:
		return len(t) == 0
	default:
		return false
	}
}
