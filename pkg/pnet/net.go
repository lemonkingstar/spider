package pnet

import (
	"errors"
	"fmt"
	"math"
	"net"
	"net/url"
)

// URLTranscode
// url转码，url转码只是为了符合url的规范，因为在标准的url规范中中文和很多的字符是不允许出现在url中的，比如空格、@字符等
func URLTranscode(src string) string {
	return url.QueryEscape(src)
}

// URLDeTranscode
// url解码
func URLDeTranscode(src string) (string, error) {
	return url.QueryUnescape(src)
}

// GetSubnetCIDR 根据CIDR均分子网
// 子网最小个数为 len(switches), 剩余子网使用默认名称 default-x命名返回
func GetSubnetCIDR(cidr string, switches []string) (map[string]string, error) {
	ip, cidrNet, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil, err
	}

	maskBit := int(math.Ceil(math.Log2(float64(len(switches)))))
	ones, bits := cidrNet.Mask.Size()
	if maskBit >= (bits - ones) {
		return nil, errors.New("sub cidr not enough")
	}

	result := map[string]string{}
	netMask := ones + maskBit
	netCount := 1 << maskBit
	ip4 := ip.To4()
	if ip4 == nil {
		return nil, errors.New("to4 converts error")
	}
	ipUint := uint32(ip4[0])<<24 | uint32(ip4[1])<<16 | uint32(ip4[2])<<8 | uint32(ip4[3])
	for i := 0; i < netCount; i++ {
		offset := uint32(i << (bits - netMask))
		ipUintMod := ipUint | offset
		ipNext := make([]byte, net.IPv4len)
		ipNext[0] = byte(ipUintMod >> 24)
		ipNext[1] = byte((ipUintMod >> 16) & 0xFF)
		ipNext[2] = byte((ipUintMod >> 8) & 0xFF)
		ipNext[3] = byte(ipUintMod & 0xFF)
		subIPNet := net.IPNet{
			IP:   ipNext,
			Mask: net.CIDRMask(netMask, bits),
		}
		if i < len(switches) {
			result[switches[i]] = subIPNet.String()
		} else {
			result[fmt.Sprintf("default-%d", i)] = subIPNet.String()
		}
	}
	return result, nil
}
