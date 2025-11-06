package phash

import (
	"strconv"
	"testing"
	"time"
)

func TestSnowflake(t *testing.T) {
	sf := &SnowFlake{}
	_ = sf.Init(1, 5)
	for i := 0; i < 10; i++ {
		t.Log(sf.NextUUID())
	}
}

func parseSnowflake(uuid string) time.Time {
	n, _ := strconv.ParseInt(uuid, 16, 64)
	// 获取毫秒
	milliSec := n >> 22
	// 加上内置的毫秒偏移量(根据实际情况处理)
	milliSec += 1672502400000

	return time.UnixMilli(milliSec)
}

func TestParseSnowflake(t *testing.T) {
	sf := &SnowFlake{}
	_ = sf.Init(1, 5)
	// uuid := "60ba14e257982bff"
	uuid, _ := sf.NextUUID()
	t.Log(uuid)

	tm := parseSnowflake(uuid)
	// str := fmt.Sprintf("%d-%d-%d %d:%d:%d %d", tm.Year(), int(tm.Month()), tm.Day(), tm.Hour(), tm.Minute(), tm.Second(), tm.Nanosecond()/1e6)
	t.Log(tm)
}

func BenchmarkSnowflake(b *testing.B) {
	b.ResetTimer()
	sf := &SnowFlake{}
	_ = sf.Init(1, 5)
	for i := 0; i < b.N; i++ {
		b.Log(sf.NextUUID())
	}
}
