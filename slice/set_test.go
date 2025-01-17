package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetSlice(t *testing.T) {
	testCase := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "normal slice",
			input: []string{"a", "b", "c", "a", "d", "c"},
			want:  []string{"a", "b", "c", "d"},
		},
		{
			name:  "empty slice",
			input: []string{},
			want:  []string(nil),
		},
		{
			name:  "not remove replicate",
			input: []string{"a", "b", "c", "d", "e", "f"},
			want:  []string{"a", "b", "c", "d", "e", "f"},
		},
		{
			name:  "all replicate",
			input: []string{"a", "a", "a", "a", "a"},
			want:  []string{"a"},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result := SetSlice[string](tc.input)
			assert.Equal(t, tc.want, result)
		})
	}
}
