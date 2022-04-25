package utils

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsPrivateAddress(t *testing.T) {
	tests := []struct {
		name    string
		address string
		want    bool
	}{
		{"t1", "127.0.0.1", true},
		{"t2", "192.168.192.168", true},
		{"t3", "180.149.134.141", false},
		{"t4", "153.186.180.56", false},
		{"t5", "10.0.0.10", true},
		{"t6", "ABC", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsPrivateAddress(tt.address); got != tt.want {
				t.Errorf("IsPrivateAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetRealIp(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		ip         string
		remoteAddr string
		want       string
	}{
		{"t1", "X-Real-Ip", "127.0.0.1", "200.200.200.200", "200.200.200.200"},
		{"t2", "X-Forwarded-For", "192.168.0.1", "200.200.200.200", "200.200.200.200"},
		{"t3", "X-Real-Ip", "123.123.123.123", "200.200.200.200", "123.123.123.123"},
		{"t4", "X-Forwarded-For", "124.124.124.124", "200.200.200.200", "124.124.124.124"},
		{"t5", "X-Forwarded-For", "AAA", "127.0.0.1:9090", "127.0.0.1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/this", http.NoBody)
			req.Header.Set(tt.tag, tt.ip)
			req.RemoteAddr = tt.remoteAddr
			if got := GetRealIp(req); got != tt.want {
				t.Errorf("GetRealIp() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetIpYlVer(t *testing.T) {
	tests := []struct {
		name       string
		tag        string
		ip         string
		remoteAddr string
		want       string
	}{
		{"t1", "X-Real-Ip", "127.0.0.1", "200.200.200.200", "127.0.0.1"},
		{"t2", "X-Forwarded-For", "192.168.0.1", "200.200.200.200", "192.168.0.1"},
		{"t3", "X-Real-Ip", "123.123.123.123,9.9.9.9", "200.200.200.200", "123.123.123.123"},
		{"t4", "X-Forwarded-For", "ABC", "200.200.200.200", "ABC"},
		{"t5", "X-Forwarded-Fo", "AAA", "127.0.0.1", "127.0.0.1"},
		{"t6", "", "", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/this", http.NoBody)
			req.Header.Set(tt.tag, tt.ip)
			req.RemoteAddr = tt.remoteAddr
			if got := GetIpYlVer(req); got != tt.want {
				t.Errorf("GetRealIp() = %v, want %v", got, tt.want)
			}
		})
	}
}
