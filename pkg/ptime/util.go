package ptime

import "time"

const (
	DateLayout        = "2006-01-02"
	TimeLayout        = "2006-01-02 15:04:05"
	RFC3339Layout     = "2006-01-02T15:04:05"
	RFC3339NanoLayout = "2006-01-02T15:04:05.000"
)

// ParseTime parses a formatted string to a time, use the system's local time zone.
func ParseTime(s string) (time.Time, error) { return time.ParseInLocation(TimeLayout, s, time.Local) }

// ParseDate parses a formatted string to a time, use the system's local time zone.
func ParseDate(s string) (time.Time, error) { return time.ParseInLocation(DateLayout, s, time.Local) }
