// Package diff computes semantic differences between config snapshots.
// It produces change sets used by the plan renderer.
package diff

import (
	"github.com/awesomeProject/apidiff/internal/model"
)

// ChangeType describes the type of change in a diff plan.
type ChangeType string

const (
	// ChangeAdd indicates a resource exists locally but not remotely.
	ChangeAdd ChangeType = "add"
	// ChangeModify indicates a resource exists in both but with different fields.
	ChangeModify ChangeType = "modify"
	// ChangeDelete indicates a resource exists remotely but not locally.
	ChangeDelete ChangeType = "delete"
)

// Change represents a resource-level change and its before/after state.
type Change struct {
	Type         ChangeType
	ResourceType string
	Key          string
	Before       any
	After        any
}

// Plan groups all changes computed between local and remote configs.
type Plan struct {
	Changes []Change
}

// HasChanges reports whether any change exists.
func (p Plan) HasChanges() bool {
	return len(p.Changes) > 0
}

// Compute compares two config snapshots and returns a change plan.
// It does not mutate either config.
func Compute(local model.Config, remote model.Config) Plan {
	plan := Plan{}

	plan.Changes = append(plan.Changes, diffCollection("route", toRouteMap(local.Routes), toRouteMap(remote.Routes))...)
	plan.Changes = append(plan.Changes, diffCollection("upstream", toUpstreamMap(local.Upstreams), toUpstreamMap(remote.Upstreams))...)
	plan.Changes = append(plan.Changes, diffCollection("service", toServiceMap(local.Services), toServiceMap(remote.Services))...)
	plan.Changes = append(plan.Changes, diffCollection("consumer", toConsumerMap(local.Consumers), toConsumerMap(remote.Consumers))...)
	plan.Changes = append(plan.Changes, diffCollection("plugin_config", toPluginConfigMap(local.PluginConfigs), toPluginConfigMap(remote.PluginConfigs))...)

	return plan
}

type keyer interface {
	Key() string
}

// diffCollection computes changes for a single resource type by key.
func diffCollection[T keyer](resourceType string, local map[string]T, remote map[string]T) []Change {
	changes := []Change{}
	for k, localItem := range local {
		if remoteItem, ok := remote[k]; ok {
			if !deepEqual(localItem, remoteItem) {
				changes = append(changes, Change{Type: ChangeModify, ResourceType: resourceType, Key: k, Before: remoteItem, After: localItem})
			}
			continue
		}
		changes = append(changes, Change{Type: ChangeAdd, ResourceType: resourceType, Key: k, After: localItem})
	}

	for k, remoteItem := range remote {
		if _, ok := local[k]; !ok {
			changes = append(changes, Change{Type: ChangeDelete, ResourceType: resourceType, Key: k, Before: remoteItem})
		}
	}
	return changes
}

func toRouteMap(items []model.Route) map[string]model.Route {
	out := make(map[string]model.Route, len(items))
	for _, it := range items {
		out[it.Key()] = it
	}
	return out
}

func toUpstreamMap(items []model.Upstream) map[string]model.Upstream {
	out := make(map[string]model.Upstream, len(items))
	for _, it := range items {
		out[it.Key()] = it
	}
	return out
}

func toServiceMap(items []model.Service) map[string]model.Service {
	out := make(map[string]model.Service, len(items))
	for _, it := range items {
		out[it.Key()] = it
	}
	return out
}

func toConsumerMap(items []model.Consumer) map[string]model.Consumer {
	out := make(map[string]model.Consumer, len(items))
	for _, it := range items {
		out[it.Key()] = it
	}
	return out
}

func toPluginConfigMap(items []model.PluginConfig) map[string]model.PluginConfig {
	out := make(map[string]model.PluginConfig, len(items))
	for _, it := range items {
		out[it.Key()] = it
	}
	return out
}
