package diff

import (
	"reflect"
	"sort"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

type FieldChange struct {
	Path   string
	Before any
	After  any
}

func FieldDiff(before, after any) []FieldChange {
	collector := &fieldCollector{changes: map[string]FieldChange{}}
	cmp.Equal(before, after, cmpopts.EquateEmpty(), cmp.Reporter(collector))

	out := make([]FieldChange, 0, len(collector.changes))
	for _, ch := range collector.changes {
		out = append(out, ch)
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].Path < out[j].Path
	})
	return out
}

type fieldCollector struct {
	path    cmp.Path
	changes map[string]FieldChange
}

func (f *fieldCollector) PushStep(ps cmp.PathStep) {
	f.path = append(f.path, ps)
}

func (f *fieldCollector) PopStep() {
	f.path = f.path[:len(f.path)-1]
}

func (f *fieldCollector) Report(result cmp.Result) {
	if result.Equal() {
		return
	}
	if len(f.path) == 0 {
		return
	}

	last := f.path.Last()
	if !isFieldStep(last) {
		return
	}

	before, after := last.Values()
	path := f.path.String()

	f.changes[path] = FieldChange{
		Path:   path,
		Before: safeInterface(before),
		After:  safeInterface(after),
	}
}

func isFieldStep(step cmp.PathStep) bool {
	switch step.(type) {
	case cmp.StructField, cmp.MapIndex, cmp.SliceIndex, cmp.ArrayIndex:
		return true
	default:
		return false
	}
}

func safeInterface(v reflect.Value) any {
	if !v.IsValid() {
		return nil
	}
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nil
		}
		return safeInterface(v.Elem())
	}
	return v.Interface()
}
