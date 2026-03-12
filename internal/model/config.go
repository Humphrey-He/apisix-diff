package model

import (
	"strings"
)

type Config struct {
	Routes        []Route        `json:"routes" yaml:"routes"`
	Upstreams     []Upstream     `json:"upstreams" yaml:"upstreams"`
	Services      []Service      `json:"services" yaml:"services"`
	Consumers     []Consumer     `json:"consumers" yaml:"consumers"`
	PluginConfigs []PluginConfig `json:"plugin_configs" yaml:"plugin_configs"`
}

func (c *Config) Normalize() {
	for i := range c.Routes {
		c.Routes[i].Normalize()
	}
	for i := range c.Upstreams {
		c.Upstreams[i].Normalize()
	}
	for i := range c.Services {
		c.Services[i].Normalize()
	}
	for i := range c.Consumers {
		c.Consumers[i].Normalize()
	}
	for i := range c.PluginConfigs {
		c.PluginConfigs[i].Normalize()
	}
}

type Route struct {
	ID         string                 `json:"id" yaml:"id"`
	Name       string                 `json:"name" yaml:"name"`
	URI        string                 `json:"uri" yaml:"uri"`
	URIs       []string               `json:"uris" yaml:"uris"`
	Methods    []string               `json:"methods" yaml:"methods"`
	UpstreamID string                 `json:"upstream_id" yaml:"upstream_id"`
	ServiceID  string                 `json:"service_id" yaml:"service_id"`
	Plugins    map[string]any         `json:"plugins" yaml:"plugins"`
	Status     int                    `json:"status" yaml:"status"`
	Priority   int                    `json:"priority" yaml:"priority"`
	Labels     map[string]string      `json:"labels" yaml:"labels"`
	Vars       []any                  `json:"vars" yaml:"vars"`
	Metadata   map[string]any         `json:"metadata" yaml:"metadata"`
}

func (r *Route) Normalize() {
	for i, m := range r.Methods {
		r.Methods[i] = strings.ToUpper(m)
	}
}

func (r Route) Key() string {
	if r.ID != "" {
		return r.ID
	}
	if r.Name != "" {
		return r.Name
	}
	if r.URI != "" {
		return r.URI
	}
	if len(r.URIs) > 0 {
		return r.URIs[0]
	}
	return "route_unknown"
}

type Upstream struct {
	ID      string            `json:"id" yaml:"id"`
	Name    string            `json:"name" yaml:"name"`
	Type    string            `json:"type" yaml:"type"`
	Nodes   Nodes             `json:"nodes" yaml:"nodes"`
	Timeout TimeoutConfig     `json:"timeout" yaml:"timeout"`
	Labels  map[string]string `json:"labels" yaml:"labels"`
}

func (u *Upstream) Normalize() {}

func (u Upstream) Key() string {
	if u.ID != "" {
		return u.ID
	}
	if u.Name != "" {
		return u.Name
	}
	return "upstream_unknown"
}

type TimeoutConfig struct {
	Connect string `json:"connect" yaml:"connect"`
	Send    string `json:"send" yaml:"send"`
	Read    string `json:"read" yaml:"read"`
}

type Service struct {
	ID      string                 `json:"id" yaml:"id"`
	Name    string                 `json:"name" yaml:"name"`
	Plugins map[string]any         `json:"plugins" yaml:"plugins"`
	Labels  map[string]string      `json:"labels" yaml:"labels"`
	Upstream *Upstream             `json:"upstream" yaml:"upstream"`
	UpstreamID string              `json:"upstream_id" yaml:"upstream_id"`
	Metadata map[string]any        `json:"metadata" yaml:"metadata"`
}

func (s *Service) Normalize() {}

func (s Service) Key() string {
	if s.ID != "" {
		return s.ID
	}
	if s.Name != "" {
		return s.Name
	}
	return "service_unknown"
}

type Consumer struct {
	ID       string                 `json:"id" yaml:"id"`
	Username string                 `json:"username" yaml:"username"`
	Plugins  map[string]any         `json:"plugins" yaml:"plugins"`
	Labels   map[string]string      `json:"labels" yaml:"labels"`
	Metadata map[string]any         `json:"metadata" yaml:"metadata"`
}

func (c *Consumer) Normalize() {}

func (c Consumer) Key() string {
	if c.ID != "" {
		return c.ID
	}
	if c.Username != "" {
		return c.Username
	}
	return "consumer_unknown"
}

type PluginConfig struct {
	ID      string                 `json:"id" yaml:"id"`
	Name    string                 `json:"name" yaml:"name"`
	Plugins map[string]any         `json:"plugins" yaml:"plugins"`
	Labels  map[string]string      `json:"labels" yaml:"labels"`
	Metadata map[string]any        `json:"metadata" yaml:"metadata"`
}

func (p *PluginConfig) Normalize() {}

func (p PluginConfig) Key() string {
	if p.ID != "" {
		return p.ID
	}
	if p.Name != "" {
		return p.Name
	}
	return "plugin_config_unknown"
}
