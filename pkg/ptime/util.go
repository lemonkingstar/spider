package ptime

import (
	"time"

	"github.com/lemonkingstar/spider/pkg/plog"
)

const (
	DateLayout         = "2006-01-02"
	TimeLayout         = "2006-01-02 15:04:05"
	RFC3339Layout      = "2006-01-02T15:04:05Z07:00"
	RFC3339MilliLayout = "2006-01-02T15:04:05.999Z07:00"
)

// ParseTime parses a formatted string to a time, use the system's local time zone.
func ParseTime(s string) (time.Time, error) { return time.ParseInLocation(TimeLayout, s, time.Local) }

// ParseTimeRFC3339 parses a formatted string to a time, auto analysis the time zone.
func ParseTimeRFC3339(s string) (time.Time, error) { return time.Parse(RFC3339Layout, s) }

// ParseDate parses a formatted string to a time, use the system's local time zone.
func ParseDate(s string) (time.Time, error) { return time.ParseInLocation(DateLayout, s, time.Local) }

func FormatTime(t time.Time) string            { return t.Format(TimeLayout) }
func FormatDate(t time.Time) string            { return t.Format(DateLayout) }
func FormatTimeRFC3339(t time.Time) string     { return t.Format(RFC3339Layout) }
func FormatTimeRFC3339Nano(t time.Time) string { return t.Format(RFC3339MilliLayout) }
func FormatNow() string                        { return FormatTime(time.Now()) }
func FormatNowRFC3339() string                 { return FormatTimeRFC3339(time.Now()) }
func FormatNowRFC3339Nano() string             { return FormatTimeRFC3339Nano(time.Now()) }

func Time2Utc(t time.Time) time.Time   { return t.UTC() }
func Time2Local(t time.Time) time.Time { return t.Local() }

func Now() time.Time      { return time.Now() }
func NowUnix() int64      { return time.Now().Unix() }
func NowUnixMilli() int64 { return time.Now().UnixMilli() }

func ParseUnix(sec int64) time.Time       { return time.Unix(sec, 0) }
func ParseUnixMilli(msec int64) time.Time { return time.UnixMilli(msec) }

func Local2Utc(localStr string) (string, error) {
	t, err := time.ParseInLocation(TimeLayout, localStr, time.Local)
	if err != nil {
		return "", err
	}
	return FormatTimeRFC3339(t.UTC()), nil
}

func Utc2Local(utcStr string) (string, error) {
	t, err := time.ParseInLocation(RFC3339Layout, utcStr, time.UTC)
	if err != nil {
		return "", err
	}
	return FormatTime(t.Local()), nil
}

// TimeElapsed calculates time elapsed.
func TimeElapsed(name string) func() int64 {
	start := time.Now()
	return func() int64 {
		millisecond := time.Since(start) / time.Millisecond
		plog.Infof("[Time Elapsed / %s] took %dms.", name, millisecond)
		return int64(millisecond)
	}
}
