package diff

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// deepEqual compares two values with empty slices/maps treated as equal.
func deepEqual(a, b any) bool {
	return cmp.Equal(a, b, cmpopts.EquateEmpty())
}
