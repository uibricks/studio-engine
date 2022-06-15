package utils

import "time"

var DbDateLayout = "2006-01-02 15:04:05"

func DbFormatDate(date time.Time) string {
	return date.Format(DbDateLayout)
}
