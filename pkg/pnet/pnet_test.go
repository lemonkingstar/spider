package pnet

import (
	"testing"
)

func TestIPString2Int(t *testing.T) {
	ip, _ := GetInternalIP()
	val := IP2Int(ip)
	t.Log(val)
	t.Log(Int2IP(val))
}

func TestGetSubnetCIDR(t *testing.T) {
	switches := []string{"cn-shanghai-a", "cn-shanghai-b"}
	t.Log(GetSubnetCIDR("10.35.224.0/22", switches))
}
