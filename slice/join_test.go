package slice

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestJoinByFunc(t *testing.T) {
	testCase := []struct {
		name   string
		input  []int
		sep    string
		output string
	}{
		{
			name:   "multi elem join",
			input:  []int{10, 20, 30},
			sep:    ",",
			output: "10,20,30",
		}, {
			name:   "single elem join",
			input:  []int{10},
			sep:    ",",
			output: "10",
		}, {
			name:   "not elem join",
			input:  []int{},
			sep:    ",",
			output: "",
		},
	}

	defaultFunc := func(i int) string {
		return strconv.Itoa(i)
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result := JoinByFunc(tc.input, tc.sep, defaultFunc)
			assert.Equal(t, tc.output, result, tc.name)
		})
	}
}
