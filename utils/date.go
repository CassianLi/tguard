package utils

import "time"

// GetCurrentYearMonth returns the current year month
func GetCurrentYearMonth(t time.Time) (string, string) {
	return t.Format("2006"), t.Format("01")
}
