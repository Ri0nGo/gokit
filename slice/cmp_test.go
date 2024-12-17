package slice

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMin(t *testing.T) {
	testCase := []struct {
		name  string
		slice []float64
		want  float64
		err   error
	}{
		{
			name:  "empty slice",
			slice: []float64{},
			err:   ErrSliceEmpty,
		},
		{
			name:  "min slice test 1",
			slice: []float64{1.1, 2.4, 1.28, 5.10, 11.1, 120, 0.8},
			want:  0.8,
		},
		{
			name:  "min slice test 2",
			slice: []float64{1.1, 2.4, 1.28, 5.10, -11.1, 120, 0.8},
			want:  -11.1,
		},
		{
			name:  "big float slice test 3",
			slice: []float64{11231.1, 214342.4, 12314231.28, 12313215.10, -112311.1, 123132120, 5552310.8},
			want:  -112311.1,
		},
		{
			name:  "fist value is min",
			slice: []float64{-1011231.1, 214342.4, 12314231.28, 12313215.10, -5311.1, 123132120, 5552310.8},
			want:  -1011231.1,
		},
		{
			name:  "last value is min",
			slice: []float64{1011231.1, 214342.4, 12314231.28, 12313215.10, -5311.1, 123132120, -5552310.8},
			want:  -5552310.8,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Min(tc.slice)
			assert.Equal(t, tc.want, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestMax(t *testing.T) {
	testCase := []struct {
		name  string
		slice []float64
		want  float64
		err   error
	}{
		{
			name:  "empty slice",
			slice: []float64{},
			err:   ErrSliceEmpty,
		},
		{
			name:  "max slice test 1",
			slice: []float64{1.1, 2.4, 1.28, 5.10, 11.1, 120, 0.8},
			want:  120,
		},
		{
			name:  "min slice test 2",
			slice: []float64{1.1, 2.4, 1.28, 5.10, -11.1, 0.1, 0.8},
			want:  5.10,
		},
		{
			name:  "big float slice test 3",
			slice: []float64{11231.1, 214342.4, 429687.28, 32112.10, -112311.1, 9923132120, 52310.8},
			want:  9923132120,
		},
		{
			name:  "fist value is max",
			slice: []float64{9923132120, 214342.4, 12314231.28, 12313215.10, -5311.1, 1232120, 5552310.8},
			want:  9923132120,
		},
		{
			name:  "last value is min",
			slice: []float64{1011231.1, 214342.4, 12314231.28, 1231321.10, -5311.1, 123132120, 9923132120},
			want:  9923132120,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Max(tc.slice)
			assert.Equal(t, tc.want, result)
			assert.Equal(t, tc.err, err)
		})
	}
}

func TestAvg(t *testing.T) {
	testCase := []struct {
		name  string
		slice []float64
		want  float64
		err   error
	}{
		{
			name:  "empty slice",
			slice: []float64{},
			err:   ErrSliceEmpty,
		},
		{
			name:  "avg slice test 1",
			slice: []float64{1, 2, 3, 4, 5},
			want:  3,
		},
		{
			name:  "min slice test 2",
			slice: []float64{-1, -2, -3, -4, -5},
			want:  -3,
		},
		{
			name:  "zero float slice test 3",
			slice: []float64{0, 0, 0, 0, 0},
			want:  0,
		},
		{
			name:  "fist value is max",
			slice: []float64{10, 100, 1000, 10000, 100000},
			want:  22222,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result, err := Avg(tc.slice)
			assert.Equal(t, tc.want, result)
			assert.Equal(t, tc.err, err)
		})
	}
}