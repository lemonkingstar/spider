package pnet

import (
	"bytes"
	"errors"
	"net"
	"strconv"
	"strings"
)

func IsInternal(ip net.IP) bool {
	return (ip[0] == 10) ||
		(ip[0] == 172 && ip[1] >= 16 && ip[1] <= 31) ||
		(ip[0] == 192 && ip[1] == 168)
}
func IsExternal(ip net.IP) bool { return !IsInternal(ip) && ip.IsGlobalUnicast() }
func isUp(v net.Flags) bool { return v&net.FlagUp == net.FlagUp }

func GetInternalIP() (string, error) { return getInterfaceIP(true) }
func GetExternalIP() (string, error) { return getInterfaceIP(false) }

func getInterfaceIP(internal bool) (string, error) {
	inters, err := net.Interfaces()
	if err != nil { return "", err }

	for _, inter := range inters {
		if !isUp(inter.Flags) { continue }
		addrList, err := inter.Addrs()
		if err != nil { continue }
		for _, addr := range addrList {
			if ipNet, ok := addr.(*net.IPNet); ok {
				if ip4 := ipNet.IP.To4(); ip4 != nil {
					if internal && IsInternal(ip4) ||
						!internal && IsExternal(ip4) {
						return ip4.String(), nil
					}
				}
			}
		}
	}
	return "", errors.New("interface not found")
}

// IP2Int ip string to int value
func IP2Int(ip string) int {
	ips := strings.Split(ip, ".")
	result := 0
	pos := 24
	for _, v := range ips {
		tmp, _ := strconv.Atoi(v)
		tmp = tmp << pos
		result = result | tmp
		pos -= 8
	}
	return result
}

// Int2IP ip int value to ip string
func Int2IP(value int) string {
	buf := bytes.NewBufferString("")
	ips := make([]string, 4)
	for i := 3; i >= 0; i-- {
		tmp := value & 0xFF
		ips[i] = strconv.Itoa(tmp)
		value = value >> 8
	}
	for i := 0; i < 4; i++ {
		buf.WriteString(ips[i])
		if i < 3 {
			buf.WriteString(".")
		}
	}
	return buf.String()
}
