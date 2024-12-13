package utils

import "time"

// NextWeekDay 获取下周的指定日期，比如获取下周一
func NextWeekDay(date time.Time, targetDay time.Weekday) time.Time {
	weekday := date.Weekday()
	daysUntilTarget := (7 + int(targetDay) - int(weekday)) % 7
	if daysUntilTarget == 0 {
		daysUntilTarget = 7
	}
	nextWeekdayDate := date.AddDate(0, 0, daysUntilTarget)
	return nextWeekdayDate
}
