package ptype

import (
	"database/sql/driver"
	"time"

	"github.com/lemonkingstar/spider/pkg/pconv"
)

// NewTimestamp returns *Timestamp at time t.
func NewTimestamp(t time.Time) *Timestamp {
	return &Timestamp{
		Seconds: t.Unix(),
		Nanos:   int32(t.Nanosecond()),
	}
}

// AsTime converts *Timestamp to a time.Time.
func (t *Timestamp) AsTime() time.Time {
	return time.Unix(t.GetSeconds(), int64(t.GetNanos()))
}

// MarshalJSON implements the json.Marshaler interface.
func (t *Timestamp) MarshalJSON() ([]byte, error) {
	ta := t.AsTime()
	return ta.MarshalJSON()
}

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var ta time.Time
	err := ta.UnmarshalJSON(data)
	if err != nil {
		return err
	}
	t.Seconds = ta.Unix()
	t.Nanos = int32(ta.Nanosecond())
	return nil
}

// Value implements the driver.Valuer interface.
func (t *Timestamp) Value() (driver.Value, error) {
	return t.AsTime(), nil
}

// Scan implements the sql.Scanner interface.
func (t *Timestamp) Scan(src any) error {
	vt, ok := src.(time.Time)
	if ok {
		t.Seconds = vt.Unix()
		t.Nanos = int32(vt.Nanosecond())
		return nil
	}
	vi, err := pconv.Type2Int64(src)
	if err != nil {
		return err
	}
	t.Seconds = vi
	return nil
}
