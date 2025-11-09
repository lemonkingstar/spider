package phash

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVerifyPassword(t *testing.T) {
	password := "password"
	b, _ := GenerateFromPassword(password)
	t.Log(b)
	assert.Nil(t, VerifyPassword(password, b))
}
