// Package apisix fetches live configuration from the APISIX Admin API.
// It converts API responses into internal model types for diffing.
package apisix

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/awesomeProject/apidiff/internal/model"
	"github.com/go-resty/resty/v2"
)

// Client wraps Admin API access with base URL, token, and timeout settings.
// It is safe for reuse across requests.
type Client struct {
	adminURL string
	token    string
	timeout  time.Duration
}

// NewClient creates a new Admin API client with the provided settings.
func NewClient(adminURL, token string, timeout time.Duration) *Client {
	return &Client{adminURL: adminURL, token: token, timeout: timeout}
}

// FetchAll pulls all supported APISIX resources via the Admin API.
// It performs read-only calls and returns a normalized config snapshot.
func (c *Client) FetchAll(ctx context.Context) (model.Config, error) {
	restyClient := resty.New().SetBaseURL(c.adminURL).SetTimeout(c.timeout)
	if c.token != "" {
		restyClient.SetHeader("X-API-KEY", c.token)
	}

	routes, err := fetchList[model.Route](ctx, restyClient, "/apisix/admin/routes")
	if err != nil {
		return model.Config{}, err
	}
	upstreams, err := fetchList[model.Upstream](ctx, restyClient, "/apisix/admin/upstreams")
	if err != nil {
		return model.Config{}, err
	}
	services, err := fetchList[model.Service](ctx, restyClient, "/apisix/admin/services")
	if err != nil {
		return model.Config{}, err
	}
	consumers, err := fetchList[model.Consumer](ctx, restyClient, "/apisix/admin/consumers")
	if err != nil {
		return model.Config{}, err
	}
	pluginConfigs, err := fetchList[model.PluginConfig](ctx, restyClient, "/apisix/admin/plugin_configs")
	if err != nil {
		return model.Config{}, err
	}

	cfg := model.Config{
		Routes:        routes,
		Upstreams:     upstreams,
		Services:      services,
		Consumers:     consumers,
		PluginConfigs: pluginConfigs,
	}
	cfg.Normalize()
	return cfg, nil
}

type listResp struct {
	List []json.RawMessage `json:"list"`
}

type listItem[T any] struct {
	Value T `json:"value"`
}

// fetchList reads a list endpoint and decodes the APISIX list schema.
// APISIX returns either list items with "value" or raw resources, so we support both.
func fetchList[T any](ctx context.Context, client *resty.Client, path string) ([]T, error) {
	resp, err := client.R().SetContext(ctx).Get(path)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("apisix admin error: %s", resp.Status())
	}

	var lr listResp
	if err := json.Unmarshal(resp.Body(), &lr); err != nil {
		return nil, err
	}

	out := make([]T, 0, len(lr.List))
	for _, raw := range lr.List {
		var wrapper listItem[T]
		if err := json.Unmarshal(raw, &wrapper); err == nil {
			out = append(out, wrapper.Value)
			continue
		}
		var item T
		if err := json.Unmarshal(raw, &item); err != nil {
			return nil, err
		}
		out = append(out, item)
	}
	return out, nil
}
