package ptime

import (
	"testing"
	"time"
)

func TestTimeUtil(t *testing.T) {
	te := TimeElapsed("TestTimeUtil")
	defer func() { te() }()

	t.Log(FormatNow())
	t.Log(FormatNowRFC3339())
	t.Log(FormatNowRFC3339Nano())
	t.Log(FormatTimeRFC3339(time.Now().UTC()))
	t.Log(FormatTimeRFC3339Nano(time.Now().UTC()))

	utcStr, _ := Local2Utc("2025-11-24 21:43:50")
	t.Log(utcStr)
	localStr, _ := Utc2Local(utcStr)
	t.Log(localStr)
}
