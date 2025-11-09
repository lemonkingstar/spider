package phash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckPassword(t *testing.T) {
	password := "password"
	b, _ := GenerateFromPassword(password)
	t.Log(b)
	assert.Nil(t, CheckHashAndPassword(password, b))
}
