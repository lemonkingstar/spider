package pginutil

import (
	"bytes"
	"io/ioutil"
	"net"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetRemoteIP(ctx *gin.Context) string {
	ip := ctx.Request.Header.Get("X-Real-IP")
	if ip != "" {
		return ip
	}

	ips := ctx.Request.Header.Get("X-Forwarded-For")
	if ips != "" {
		list := strings.Split(ips, ",")
		if len(list) > 0 {
			return list[0]
		}
	}

	return ctx.ClientIP()
}

func GetRemoteIP2(req *http.Request) string {
	ip := req.Header.Get("X-Real-IP")
	netIP := net.ParseIP(ip)
	if netIP != nil { return ip }

	ips := req.Header.Get("X-Forwarded-For")
	ipList := strings.Split(ips, ",")
	for _, ip := range ipList {
		netIP := net.ParseIP(ip)
		if netIP != nil { return ip }
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil { return "" }
	netIP = net.ParseIP(ip)
	if netIP != nil { return ip }
	return ""
}

func GetRemoteIP3(req *http.Request) string {
	ips := req.Header.Get("X-Forwarded-For")
	ipList := strings.Split(ips, ",")
	for _, ip := range ipList {
		netIP := net.ParseIP(ip)
		if netIP != nil { return ip }
	}

	ip, _, err := net.SplitHostPort(req.RemoteAddr)
	if err != nil { return "" }
	netIP := net.ParseIP(ip)
	if netIP != nil { return ip }

	ip = req.Header.Get("X-Real-IP")
	netIP = net.ParseIP(ip)
	if netIP != nil { return ip }
	return ""
}

func PeekRequest(req *http.Request) ([]byte, error) {
	if req.Body != nil {
		byt, err := ioutil.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		req.Body = ioutil.NopCloser(bytes.NewBuffer(byt))
		return byt, nil
	}
	return []byte{}, nil
}
