package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilter(t *testing.T) {
	testCase := []struct {
		name      string
		slice     []int
		match     filterFunc[int]
		wantSlice []int
	}{
		{
			name:  "index 0",
			slice: []int{1, 2, 3, 4},
			match: func(elem int) bool {
				return elem == 2
			},
			wantSlice: []int{2},
		},
		{
			name:  "index -1",
			slice: []int{1, 2, 3, 4},
			match: func(elem int) bool {
				return elem == 5
			},
			wantSlice: []int{},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result := Filter[int](tc.slice, tc.match)
			assert.Equal(t, result, tc.wantSlice)
		})
	}
}
