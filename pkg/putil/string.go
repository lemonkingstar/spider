package putil

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/google/uuid"
	"github.com/matoous/go-nanoid/v2"
)

func GenRandString(length int) string {
	r := rand.New(rand.NewSource(time.Now().Unix()))
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		if x := r.Intn(3); x == 1 {
			bytes[i] = byte(r.Intn(10) + 48)
		} else {
			bytes[i] = byte(r.Intn(26) + 97)
		}
	}
	return string(bytes)
}

func Uuid(prefix string) string {
	var id string
	u, err := uuid.NewRandom()
	if err != nil {
		id = GenRandString(36)
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
