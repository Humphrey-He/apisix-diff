package render

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/awesomeProject/apidiff/internal/diff"
)

func RenderPlan(w io.Writer, plan diff.Plan) {
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
			fmt.Fprintf(w, "+ %s.%s\n", c.ResourceType, c.Key)
		case diff.ChangeDelete:
			fmt.Fprintf(w, "- %s.%s\n", c.ResourceType, c.Key)
		case diff.ChangeModify:
			fmt.Fprintf(w, "~ %s.%s\n", c.ResourceType, c.Key)
			fields := diff.FieldDiff(c.Before, c.After)
			if len(fields) == 0 {
				fmt.Fprintln(w, "  (no field-level diff available)")
				continue
			}
			for _, ch := range fields {
				fmt.Fprintf(w, "  %s: %s -> %s\n", ch.Path, formatValue(ch.Before), formatValue(ch.After))
			}
		}
	}
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
