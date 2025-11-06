package putil

import (
	"crypto/rand"
	"fmt"
	randx "math/rand"
	"time"

	"github.com/satori/go.uuid"
)

// GetRandString 随机生成指定长度的字符串
func GetRandString(length int) string {
	r := randx.New(randx.NewSource(time.Now().Unix()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		b := r.Intn(26) + 97
		bytes[i] = byte(b)
	}
	return string(bytes)
}

// GenerateUUID 生成32位十六进制数字/以-分隔的全局唯一标识
// UUID Version 4
func GenerateUUID() (string, error) {
	buf := make([]byte, 16)
	_, err := rand.Read(buf)
	if err != nil { return "", err }
	return fmt.Sprintf("%x-%x-%x-%x-%x", buf[0:4], buf[4:6], buf[6:8], buf[8:10], buf[10:]), nil
}

// UUID Version 4
func UUID() string {
	return uuid.NewV4().String()
}
