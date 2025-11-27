package putil

import "testing"

func TestGenerateUUID(t *testing.T) {
	t.Log(GetRandString(32))
	uuid, _ := GenerateUUID()
	t.Log(uuid)
	t.Log(UUID())
}
