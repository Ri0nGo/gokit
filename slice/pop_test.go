package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPop(t *testing.T) {
	testCase := []struct {
		name      string
		slice     []int
		index     int
		wantVal   int
		wantSlice []int
		wantErr   error
	}{
		{
			name:      "Index middle",
			slice:     []int{10, 20, 30},
			index:     1,
			wantVal:   20,
			wantSlice: []int{10, 30},
		},
		{
			name:      "Index first",
			slice:     []int{10, 20, 30},
			index:     0,
			wantVal:   10,
			wantSlice: []int{20, 30},
		},
		{
			name:      "Index last",
			slice:     []int{10, 20, 30},
			index:     2,
			wantVal:   30,
			wantSlice: []int{10, 20},
		},
		{
			name:      "Index -1",
			slice:     []int{10, 20, 30},
			index:     -1,
			wantVal:   0,
			wantSlice: nil,
			wantErr:   errIndexOutOfSlice,
		},
		{
			name:      "Index out of",
			slice:     []int{10, 20, 30},
			index:     12,
			wantVal:   0,
			wantSlice: nil,
			wantErr:   errIndexOutOfSlice,
		},
	}
	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, val, err := Pop[int](tc.slice, tc.index)
			assert.Equal(t, tc.wantErr, err)
			if err != nil {
				return
			}
			assert.Equal(t, tc.wantSlice, result)
			assert.Equal(t, tc.wantVal, val)
		})
	}
}
