package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandom(t *testing.T) {
	testCase := []struct {
		name   string
		length int
		want   int
	}{
		{
			name:   "special length string",
			length: 10,
			want:   10,
		}, {
			name:   "one length string",
			length: 1,
			want:   1,
		}, {
			name:   "length is zero",
			length: 0,
			want:   0,
		}, {
			name:   "length is negative number",
			length: -1,
			want:   0,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result := RandString(tc.length)
			assert.Equal(t, tc.want, len(result), tc.name)
		})
	}
}
