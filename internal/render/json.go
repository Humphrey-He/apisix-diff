package render

import (
	"encoding/json"
	"io"

	"github.com/awesomeProject/apidiff/internal/diff"
)

// PlanSummary is the top-level metadata for JSON output.
type PlanSummary struct {
	Add    int `json:"add"`
	Change int `json:"change"`
	Delete int `json:"delete"`
}

// PlanOutput is the JSON representation of a diff plan.
type PlanOutput struct {
	Summary PlanSummary   `json:"summary"`
	Changes []diff.Change `json:"changes"`
}

// RenderPlanJSON writes a machine-readable diff plan in JSON format.
func RenderPlanJSON(w io.Writer, plan diff.Plan) error {
	summary := PlanSummary{}
	for _, c := range plan.Changes {
		switch c.Type {
		case diff.ChangeAdd:
			summary.Add++
		case diff.ChangeModify:
			summary.Change++
		case diff.ChangeDelete:
			summary.Delete++
		}
	}

	payload := PlanOutput{Summary: summary, Changes: plan.Changes}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(payload)
}
