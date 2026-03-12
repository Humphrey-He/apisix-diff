// Package render prints human-readable diff plans.
package render

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/awesomeProject/apidiff/internal/diff"
)

// Options configures output behavior.
type Options struct {
	// Color enables ANSI color codes in output.
	Color bool
}

// RenderPlan writes a plan-style diff to the provided writer.
// It groups field-level changes by top-level field name.
func RenderPlan(w io.Writer, plan diff.Plan, opts Options) {
	adds := 0
	changes := 0
	deletes := 0
	for _, c := range plan.Changes {
		switch c.Type {
		case diff.ChangeAdd:
			adds++
		case diff.ChangeModify:
			changes++
		case diff.ChangeDelete:
			deletes++
		}
	}

	fmt.Fprintf(w, "Plan: %d to add, %d to change, %d to delete.\n\n", adds, changes, deletes)

	sorted := make([]diff.Change, len(plan.Changes))
	copy(sorted, plan.Changes)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].ResourceType == sorted[j].ResourceType {
			return sorted[i].Key < sorted[j].Key
		}
		return sorted[i].ResourceType < sorted[j].ResourceType
	})

	for _, c := range sorted {
		switch c.Type {
		case diff.ChangeAdd:
			fmt.Fprintf(w, "%s %s.%s\n", colorize("+", colorGreen, opts.Color), c.ResourceType, c.Key)
		case diff.ChangeDelete:
			fmt.Fprintf(w, "%s %s.%s\n", colorize("-", colorRed, opts.Color), c.ResourceType, c.Key)
		case diff.ChangeModify:
			fmt.Fprintf(w, "%s %s.%s\n", colorize("~", colorYellow, opts.Color), c.ResourceType, c.Key)
			fields := diff.FieldDiff(c.Before, c.After)
			if len(fields) == 0 {
				fmt.Fprintln(w, "  (no field-level diff available)")
				continue
			}

			groups := groupFieldChanges(fields)
			groupKeys := make([]string, 0, len(groups))
			for k := range groups {
				groupKeys = append(groupKeys, k)
			}
			sort.Strings(groupKeys)

			for _, group := range groupKeys {
				fmt.Fprintf(w, "  %s:\n", group)
				for _, ch := range groups[group] {
					fmt.Fprintf(w, "    %s: %s %s %s\n", ch.Path, formatValue(ch.Before), colorize("->", colorDim, opts.Color), formatValue(ch.After))
				}
			}
		}
	}
}

type groupedChange struct {
	Path   string
	Before any
	After  any
}

func groupFieldChanges(fields []diff.FieldChange) map[string][]groupedChange {
	out := map[string][]groupedChange{}
	for _, ch := range fields {
		// Group by the first path segment for readability.
		group, leaf := splitGroupPath(ch.Path)
		out[group] = append(out[group], groupedChange{Path: leaf, Before: ch.Before, After: ch.After})
	}
	for group := range out {
		sort.Slice(out[group], func(i, j int) bool {
			return out[group][i].Path < out[group][j].Path
		})
	}
	return out
}

func splitGroupPath(path string) (string, string) {
	trimmed := strings.TrimPrefix(path, ".")
	if trimmed == "" {
		return "root", path
	}

	parts := strings.Split(trimmed, ".")
	if len(parts) == 1 {
		return parts[0], parts[0]
	}

	group := parts[0]
	leaf := strings.Join(parts[1:], ".")
	return group, leaf
}

func formatValue(v any) string {
	if v == nil {
		return "null"
	}
	switch t := v.(type) {
	case string:
		return fmt.Sprintf("%q", t)
	case fmt.Stringer:
		return t.String()
	}

	if b, err := json.Marshal(v); err == nil {
		text := string(b)
		if strings.HasPrefix(text, "{") || strings.HasPrefix(text, "[") {
			return text
		}
	}

	return fmt.Sprintf("%v", v)
}

const (
	colorReset  = "\x1b[0m"
	colorRed    = "\x1b[31m"
	colorGreen  = "\x1b[32m"
	colorYellow = "\x1b[33m"
	colorDim    = "\x1b[2m"
)

func colorize(text, color string, enabled bool) string {
	if !enabled {
		return text
	}
	return color + text + colorReset
}
