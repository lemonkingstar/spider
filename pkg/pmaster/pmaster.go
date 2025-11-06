package pmaster

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/lemonkingstar/spider/pkg/plog"
)

var (
	logger = plog.WithField("[PACKET]", "pmaster")
)

var (
	machineID     int64 // 机器 id 占10位, 十进制范围是 [ 0, 1023 ]
	sn            int64 // 序列号占 12 位,十进制范围是 [ 0, 4095 ]
	lastTimeStamp int64 // 上次的时间戳(毫秒级), 1秒=1000毫秒, 1毫秒=1000微秒,1微秒=1000纳秒
)

func buildSnowflakeUUID() string {
	rand.Seed(time.Now().UnixNano())
	machineID = rand.Int63n(1024)
	machineID <<= 12
	sn = rand.Int63n(4096)

	curTimeStamp := time.Now().UnixNano() / 1000000
	rightBinValue := curTimeStamp & 0x1FFFFFFFFFF
	rightBinValue <<= 22
	id := rightBinValue | machineID | sn
	return fmt.Sprintf("%x", id)
}
