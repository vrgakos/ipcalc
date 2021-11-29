package ipcalc

import (
	"fmt"
	"math/rand"
	"net"
	"testing"
	"time"
)

func TestIntersect(t *testing.T) {
	var tests = []struct {
		s1, s2 string
		want   bool
	}{
		{"10.0.0.0/16", "10.0.0.0/16", true},
		{"10.0.0.0/16", "10.0.0.0/24", true},
		{"10.0.0.0/16", "10.0.1.0/24", true},
		{"10.0.0.0/16", "10.0.1.10/24", true},
		{"10.0.0.0/16", "10.0.255.0/24", true},

		{"10.0.0.0/16", "10.0.0.0/32", true},
		{"10.0.0.0/16", "10.0.1.0/32", true},
		{"10.0.0.0/16", "10.0.1.10/32", true},
		{"10.0.0.0/16", "10.0.255.0/32", true},

		{"10.0.0.0/16", "10.1.0.0/24", false},
		{"10.0.0.0/16", "10.1.3.4/32", false},

		{"11.0.0.0/16", "10.0.0.0/16", false},
		{"11.0.0.0/16", "10.0.0.0/24", false},
		{"11.0.0.0/16", "10.0.1.0/24", false},
		{"11.0.0.0/16", "10.0.1.10/24", false},
		{"11.0.0.0/16", "10.0.255.0/24", false},

		{"10.2.0.0/16", "10.0.0.0/16", false},
		{"10.2.0.0/16", "10.0.0.0/24", false},
		{"10.2.0.0/16", "10.0.1.0/24", false},
		{"10.2.0.0/16", "10.0.1.10/24", false},
		{"10.2.0.0/16", "10.0.255.0/24", false},
	}

	for _, tt := range tests {
		s1 := NewSubnet(tt.s1)
		s2 := NewSubnet(tt.s2)

		got := s1.Intersect(s2)
		if got != tt.want {
			t.Errorf("got %t, want %t", got, tt.want)
		}

		_, net1, _ := net.ParseCIDR(tt.s1)
		_, net2, _ := net.ParseCIDR(tt.s2)

		want2 := net2.Contains(net1.IP) || net1.Contains(net2.IP)
		if got != want2 {
			t.Errorf("got %t, want2 %t", got, want2)
		}
	}
}

func TestContains(t *testing.T) {
	var tests = []struct {
		s1, s2 string
		want   bool
	}{
		{"10.0.0.0/16", "10.0.0.0/16", true},
		{"10.0.0.0/16", "10.0.0.0/24", true},
		{"10.0.0.0/16", "10.0.1.0/24", true},
		{"10.0.0.0/16", "10.0.1.10/24", true},
		{"10.0.0.0/16", "10.0.255.0/24", true},

		{"10.0.0.0/16", "10.0.0.0/32", true},
		{"10.0.0.0/16", "10.0.1.0/32", true},
		{"10.0.0.0/16", "10.0.1.10/32", true},
		{"10.0.0.0/16", "10.0.255.0/32", true},

		{"10.0.0.0/16", "10.1.0.0/24", false},
		{"10.0.0.0/16", "10.1.3.4/32", false},

		{"11.0.0.0/16", "10.0.0.0/16", false},
		{"11.0.0.0/16", "10.0.0.0/24", false},
		{"11.0.0.0/16", "10.0.1.0/24", false},
		{"11.0.0.0/16", "10.0.1.10/24", false},
		{"11.0.0.0/16", "10.0.255.0/24", false},

		{"10.2.0.0/16", "10.0.0.0/16", false},
		{"10.2.0.0/16", "10.0.0.0/24", false},
		{"10.2.0.0/16", "10.0.1.0/24", false},
		{"10.2.0.0/16", "10.0.1.10/24", false},
		{"10.2.0.0/16", "10.0.255.0/24", false},
	}

	for _, tt := range tests {
		s1 := NewSubnet(tt.s1)
		s2 := NewSubnet(tt.s2)

		got := s1.Contains(s2)
		if got != tt.want {
			t.Errorf("got %t, want %t", got, tt.want)
		}

		_, net1, _ := net.ParseCIDR(tt.s1)
		_, net2, _ := net.ParseCIDR(tt.s2)

		want2 := net1.Contains(net2.IP)
		if got != want2 {
			t.Errorf("got %t, want2 %t", got, want2)
		}
	}
}

func init() {
	rand.Seed(time.Now().Unix())
}

func randIPv4Addr() string {
	return fmt.Sprintf("%d.%d.%d.%d", rand.Uint32()%256, rand.Uint32()%256, rand.Uint32()%256, rand.Uint32()%256)
}

func randIPv4Subnet() string {
	return fmt.Sprintf("%s/%d", randIPv4Addr(), rand.Uint32()%33)
}

func randIPv6Addr() string {
	return fmt.Sprintf("%x:%x:%x:%x:%x:%x:%x:%x", rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536, rand.Uint32()%65536)
}

func randIPv6Subnet() string {
	return fmt.Sprintf("%s/%d", randIPv6Addr(), rand.Uint32()%129)
}

func TestContainsRandomIpv4(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s1Str := randIPv4Subnet()
		s2Str := randIPv4Subnet()

		s1 := NewSubnet(s1Str)
		s2 := NewSubnet(s2Str)

		got := s1.Contains(s2)

		_, net1, _ := net.ParseCIDR(s1Str)
		_, net2, _ := net.ParseCIDR(s2Str)

		want := net1.Contains(net2.IP)
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
}

func TestContainsRandomIpv6(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s1Str := randIPv6Subnet()
		s2Str := randIPv6Subnet()

		s1 := NewSubnet(s1Str)
		s2 := NewSubnet(s2Str)

		got := s1.Contains(s2)

		_, net1, _ := net.ParseCIDR(s1Str)
		_, net2, _ := net.ParseCIDR(s2Str)

		want := net1.Contains(net2.IP)
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
}

func TestIntersectRandomIpv4(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s1Str := randIPv4Subnet()
		s2Str := randIPv4Subnet()

		s1 := NewSubnet(s1Str)
		s2 := NewSubnet(s2Str)

		got := s1.Intersect(s2)

		_, net1, _ := net.ParseCIDR(s1Str)
		_, net2, _ := net.ParseCIDR(s2Str)

		want := net1.Contains(net2.IP) || net2.Contains(net1.IP)
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
}

func TestIntersectRandomIpv6(t *testing.T) {
	for i := 0; i < 10000; i++ {
		s1Str := randIPv6Subnet()
		s2Str := randIPv6Subnet()

		s1 := NewSubnet(s1Str)
		s2 := NewSubnet(s2Str)

		got := s1.Intersect(s2)

		_, net1, _ := net.ParseCIDR(s1Str)
		_, net2, _ := net.ParseCIDR(s2Str)

		want := net1.Contains(net2.IP) || net2.Contains(net1.IP)
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	}
}

func TestGetCommonOnes(t *testing.T) {
	var tests = []struct {
		s1, s2 string
		want   uint8
	}{
		{"10.0.0.0/16", "10.0.0.0/16", 16},
		{"10.0.0.0/16", "10.0.0.0/24", 16},
		{"10.0.0.0/32", "10.0.0.0/32", 32},
		{"10.0.0.0/32", "10.0.0.1/32", 31},
		{"10.0.0.0/32", "10.0.0.2/32", 30},
		{"10.0.0.0/32", "10.0.0.3/32", 30},
		{"10.0.0.0/32", "10.0.0.255/32", 24},
		{"10.0.0.0/32", "10.0.0.255/32", 24},
		{"10.0.0.0/32", "255.0.0.0/8", 0},
		{"10.0.0.0/32", "10.0.0.0/8", 8},
	}

	for _, tt := range tests {
		s1 := NewSubnet(tt.s1)
		s2 := NewSubnet(tt.s2)

		got := s1.CommonOnes(s2, true)
		if got != tt.want {
			t.Errorf("got %d, want %d", got, tt.want)
		}
	}
}

func BenchmarkIntersect(b *testing.B) {
	s1 := NewSubnet(randIPv4Subnet())
	s2 := NewSubnet(randIPv4Subnet())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s1.Intersect(s2)
	}
}

func BenchmarkContains(b *testing.B) {
	s1 := NewSubnet(randIPv4Subnet())
	s2 := NewSubnet(randIPv4Subnet())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		s1.Contains(s2)
	}
}

func BenchmarkGoNativeIntersect(b *testing.B) {
	_, net1, _ := net.ParseCIDR(randIPv4Subnet())
	_, net2, _ := net.ParseCIDR(randIPv4Subnet())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		net2.Contains(net1.IP)
		net1.Contains(net2.IP)
	}
}

func BenchmarkGoNativeContains(b *testing.B) {
	_, net1, _ := net.ParseCIDR(randIPv4Subnet())
	_, net2, _ := net.ParseCIDR(randIPv4Subnet())

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		net1.Contains(net2.IP)
	}
}
