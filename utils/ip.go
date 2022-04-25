package utils

import (
	"net"
	"net/http"
	"strings"
)

var maxCidrBlocks = [...]string{
	"127.0.0.1/8",    // localhost
	"10.0.0.0/8",     // 24-bit block
	"172.16.0.0/12",  // 20-bit block
	"192.168.0.0/16", // 16-bit block
	"169.254.0.0/16", // link local address
	"::1/128",        // localhost IPv6
	"fc00::/7",       // unique local address IPv6
	"fe80::/10",      // link local address IPv6
}

var ipnets []*net.IPNet

func init() {
	ipnets = make([]*net.IPNet, len(maxCidrBlocks))
	for i, v := range maxCidrBlocks {
		_, ipnet, _ := net.ParseCIDR(v)
		ipnets[i] = ipnet
	}
}

/* 不為私人IP且為有效IP */
func IsValidIp(address string) bool {
	ip := net.ParseIP(address)
	if ip == nil {
		return false
	}
	for _, value := range ipnets {
		if value.Contains(ip) {
			return false
		}
	}
	return true
}

/* 確定是私人IP */
func IsPrivateAddress(address string) bool {
	ipAddress := net.ParseIP(address)
	if ipAddress == nil {
		return false
	}
	for _, value := range ipnets {
		if value.Contains(ipAddress) {
			return true
		}
	}
	return false
}

/* 取得真實IP 標籤若是私人IP則往下判斷 */
func GetRealIp(req *http.Request) string {
	xRealIP := req.Header.Get("X-Real-Ip")
	xForwardedFor := req.Header.Get("X-Forwarded-For")
	if xRealIP != "" {
		for _, address := range strings.Split(xRealIP, ",") {
			address = strings.TrimSpace(address)
			if IsValidIp(address) { // 有效IP
				return address
			}
		}
	}
	if xForwardedFor != "" {
		for _, address := range strings.Split(xForwardedFor, ",") {
			address = strings.TrimSpace(address)
			if IsValidIp(address) { // 有效IP
				return address
			}
		}
	}
	index := strings.Index(req.RemoteAddr, ":")
	if index != -1 {
		return req.RemoteAddr[:index]
	}
	return req.RemoteAddr
}

/* YL判斷IP版本 */
func GetIpYlVer(req *http.Request) string {
	xreal := req.Header.Get("x-real-ip")
	xforwarded := req.Header.Get("x-forwarded-for")
	var ip string
	if xreal != "" {
		ip = xreal
	} else if xforwarded != "" {
		ip = xforwarded
	} else {
		index := strings.Index(req.RemoteAddr, ":")
		if index != -1 {
			ip = req.RemoteAddr[:index]
		} else {
			ip = req.RemoteAddr
		}
	}
	ips := strings.Split(ip, ",")
	return strings.ReplaceAll(ips[0], `::ffff:`, "")
}
