package archive

import (
	"strings"
	"time"
)

func parseDay(day string) (time.Time, bool) {
	now := time.Now().Local()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
	switch strings.ToLower(day) {
	case "сегодня":
		return today, true
	case "вчера":
		return today.AddDate(0, 0, -1), true
	case "позавчера":
		return today.AddDate(0, 0, -2), true
	default:
		tm, err := time.ParseInLocation("2006-01-02", day, time.Local)
		if err != nil {
			return time.Time{}, false
		}
		return tm, true
	}
}

func parseTime(tm string) (time.Duration, bool) {
	t, err := time.ParseInLocation("15:04:05", tm, time.Local)
	if err != nil {
		return 0, false
	}
	return time.Duration(t.Hour())*time.Hour + time.Duration(t.Minute())*time.Minute + time.Duration(t.Second())*time.Second, true
}
