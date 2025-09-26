package utils

import (
	"log"
	"time"
)

// time format: ISO 8601, timezone format: "Asia/Jakarta"
// Usage example: FormatDateTime("2025-01-25T12:30:00Z", "Asia/Jakarta")
func parseISOToLocation(dateString, timezone string) time.Time {
	datetime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		log.Printf("Invalid time string: %v", err)
		return time.Time{}
	}

	localization, err := time.LoadLocation(timezone)
	if err != nil {
		log.Printf("Invalid timezone: %v", err)
		localization = time.UTC
	}

	return datetime.In(localization)
}

func FormatDateTime(dateString, timezone string) string {
	datetime := parseISOToLocation(dateString, timezone)
	if datetime.IsZero() {
		return ""
	}
	return datetime.Format("02 January 2006, 03:04 PM")
}

func FormatDate(dateString, timezone string) string {
	datetime := parseISOToLocation(dateString, timezone)
	if datetime.IsZero() {
		return ""
	}
	return datetime.Format("02 January 2006")
}
