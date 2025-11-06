package phash

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func Md5Sign(src string) string {
	has := md5.Sum([]byte(src))
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

func Md5Sign2Upper(src string) string {
	return strings.ToUpper(Md5Sign(src))
}
