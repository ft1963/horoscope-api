package main

import "time"

// 日付文字列をパースする共通関数
func parseDate(dateStr string) time.Time {
	if dateStr == "" {
		return time.Now()
	}
	// HTML5 date input format (YYYY-MM-DD)
	if t, err := time.Parse("2006-01-02", dateStr); err == nil {
		return t
	}
	// RFC3339 format
	if t, err := time.Parse(time.RFC3339, dateStr); err == nil {
		return t
	}
	return time.Now()
}
