package model

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type Nodes map[string]int

type nodeItem struct {
	Host   string `json:"host" yaml:"host"`
	Port   int    `json:"port" yaml:"port"`
	Weight int    `json:"weight" yaml:"weight"`
}

func (n *Nodes) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.MappingNode {
		var m map[string]int
		if err := value.Decode(&m); err != nil {
			return err
		}
		*n = Nodes(m)
		return nil
	}
	if value.Kind == yaml.SequenceNode {
		var items []nodeItem
		if err := value.Decode(&items); err != nil {
			return err
		}
		out := Nodes{}
		for _, it := range items {
			if it.Port == 0 && it.Host != "" {
				return fmt.Errorf("node port required for host %s", it.Host)
			}
			key := it.Host
			if it.Port > 0 {
				key = fmt.Sprintf("%s:%d", it.Host, it.Port)
			}
			weight := it.Weight
			if weight == 0 {
				weight = 1
			}
			out[key] = weight
		}
		*n = out
		return nil
	}
	return fmt.Errorf("unsupported nodes format")
}

func (n *Nodes) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return nil
	}
	if data[0] == '{' {
		var m map[string]int
		if err := json.Unmarshal(data, &m); err != nil {
			return err
		}
		*n = Nodes(m)
		return nil
	}
	if data[0] == '[' {
		var items []nodeItem
		if err := json.Unmarshal(data, &items); err != nil {
			return err
		}
		out := Nodes{}
		for _, it := range items {
			if it.Port == 0 && it.Host != "" {
				return fmt.Errorf("node port required for host %s", it.Host)
			}
			key := it.Host
			if it.Port > 0 {
				key = fmt.Sprintf("%s:%d", it.Host, it.Port)
			}
			weight := it.Weight
			if weight == 0 {
				weight = 1
			}
			out[key] = weight
		}
		*n = out
		return nil
	}
	return fmt.Errorf("unsupported nodes format")
}

func (n Nodes) ToAddressList() []string {
	out := make([]string, 0, len(n))
	for key := range n {
		out = append(out, key)
	}
	return out
}
