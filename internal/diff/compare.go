package diff

import (
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func deepEqual(a, b any) bool {
	return cmp.Equal(a, b, cmpopts.EquateEmpty())
}
