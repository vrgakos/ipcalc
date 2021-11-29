package ipcalc

import (
	"net"
	"testing"
)

func TestIpToInt(t *testing.T) {
	var tests = []struct {
		ipStr, valueStr string
		bits            uint8
	}{
		// https://www.ipaddressguide.com/ip
		{"0.0.0.0", "0", 32},
		{"0.0.0.1", "1", 32},
		{"1.1.1.1", "16843009", 32},
		{"10.0.0.0", "167772160", 32},
		{"144.51.10.96", "2419264096", 32},

		// https://www.ipaddressguide.com/ipv6-to-decimal
		{"::", "0", 128},
		{"::1", "1", 128},
		{"fe80:abba:edda:acdc::", "338292008029353851204481704463429533696", 128},
		{"fe80:abba:edda:acdc::1000", "338292008029353851204481704463429537792", 128},
		{"2001:db8::", "42540766411282592856903984951653826560", 128},
		{"2001:4860:4860::8888", "42541956123769884636017138956568135816", 128},
	}

	for _, tt := range tests {
		t.Run(tt.ipStr, func(t *testing.T) {
			ip := net.ParseIP(tt.ipStr)

			val, gotBits := ipToInt(ip)
			gotValue := val.String()
			if gotValue != tt.valueStr {
				t.Errorf("gotValue %s, want %s", gotValue, tt.valueStr)
			}

			if gotBits != int(tt.bits) {
				t.Errorf("gotBits %d, want %d", gotBits, tt.bits)
			}
		})
	}
}
