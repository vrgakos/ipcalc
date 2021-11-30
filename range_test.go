package ipcalc

import (
	"testing"
)

func TestRangeParse(t *testing.T) {
	var tests = []struct {
		start, end string
		wantedSize string
	}{
		{"192.168.0.1", "192.168.0.2", "2"},
		{"192.168.0.5", "192.168.0.9", "5"},
		{"192.0.0.100", "192.0.0.199", "100"},
		{"192.0.0.199", "192.0.0.100", ""},
		{"192.0.0.199", "2001:db8::ff", ""},
		{"2001:db8::1", "2001:db8::2", "2"},
		{"2001:db8::1", "2001:db8::ffff", "65535"},
		{"2001:eb8::1", "2001:a::1", ""},
	}

	for _, tt := range tests {
		r, _ := ParseRange(tt.start, tt.end)

		if r == nil {
			if tt.wantedSize == "" {
				continue
			}
			t.Errorf("got nil, wanted success")
			break
		}

		if r.Size.String() != tt.wantedSize {
			t.Errorf("got %s, want %s", r.Size.String(), tt.wantedSize)
		}
	}
}
