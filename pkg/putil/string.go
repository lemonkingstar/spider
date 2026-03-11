package putil

import (
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/matoous/go-nanoid/v2"
	"github.com/mozillazg/go-pinyin"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

func GenRandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	data := make([]byte, length)
	for i := 0; i < length; i++ {
		if x := r.Intn(3); x == 1 {
			data[i] = byte(r.Intn(10) + 48)
		} else {
			data[i] = byte(r.Intn(26) + 97)
		}
	}
	return string(data)
}

func Uuid(prefix string) string {
	var id string
	u, err := uuid.NewRandom()
	if err != nil {
		id = GenRandString(32)
	} else {
		id = u.String()
	}
	if prefix != "" {
		id = fmt.Sprintf("%s-%s", prefix, id)
	}
	return id
}

func ShortUuid(prefix string) string {
	id, err := gonanoid.New(17)
	if err != nil {
		id = GenRandString(17)
	}
	if prefix != "" {
		id = fmt.Sprintf("%s-%s", prefix, id)
	}
	return id
}

func Utf8ToGbk(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewEncoder())
	data, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return data, nil
}

func GbkToUtf8(b []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(b), simplifiedchinese.GBK.NewDecoder())
	data, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return data, nil
}

// GetInitials 获取文本首字母
func GetInitials(input string) string {
	if len(input) == 0 {
		return ""
	}

	for _, initials := range input {
		pinyinResult := pinyin.SinglePinyin(initials, pinyin.NewArgs())
		for _, p := range pinyinResult {
			if len(p) > 0 {
				return string(p[0])
			}
		}

		return string(initials)
	}

	return ""
}

// TextToPinyin 文本转拼音
func TextToPinyin(input string) string {
	pinyinSlice := pinyin.LazyPinyin(input, pinyin.Args{
		Style: pinyin.Normal,
		Fallback: func(r rune, a pinyin.Args) []string {
			return []string{string(r)}
		},
	})
	
	return strings.Join(pinyinSlice, "")
}
