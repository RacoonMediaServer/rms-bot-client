package archive

import "time"

func parseDay(day string) (time.Time, bool) {
	now := time.Now()
	tm := time.Time{}
	tm = tm.AddDate(now.Year(), int(now.Month()), now.Day())
	switch day {
	case "Вчера":
		tm = tm.Sub
	}
	return time.Time{}, false
}

func parseTime(tm string) (time.Duration, bool) {
	return time.Duration(0), false
}
