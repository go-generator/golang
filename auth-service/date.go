package auth_service

import "time"

func addSeconds(date time.Time, seconds int) time.Time {
	return date.Add(time.Second * time.Duration(seconds))
}

func addDays(date time.Time, days int) time.Time {
	return date.Add(time.Hour * time.Duration(days) * 24)
}

func compareDate(date1 time.Time, date2 time.Time) int {
	//return date1.Second() - date2.Second()
	return int(date1.Sub(date2).Seconds())
}
