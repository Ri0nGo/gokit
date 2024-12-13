package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestNextWeekDay(t *testing.T) {
	testCase := []struct {
		name          string
		targetWeekDay time.Weekday
		input         time.Time
		want          time.Time
	}{
		{
			name:          "获取下周一",
			targetWeekDay: time.Monday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2024, time.December, 30, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周二",
			targetWeekDay: time.Tuesday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2024, time.December, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周三",
			targetWeekDay: time.Wednesday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2025, time.January, 1, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周四",
			targetWeekDay: time.Thursday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2025, time.January, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周五",
			targetWeekDay: time.Friday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2025, time.January, 3, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周六",
			targetWeekDay: time.Saturday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2025, time.January, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			name:          "获取下周日",
			targetWeekDay: time.Sunday,
			input:         time.Date(2024, time.December, 29, 0, 0, 0, 0, time.UTC),
			want:          time.Date(2025, time.January, 5, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			result := NextWeekDay(tc.input, tc.targetWeekDay)
			assert.Equal(t, result, tc.want)
		})
	}
}
