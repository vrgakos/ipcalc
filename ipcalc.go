package ipcalc

import (
	"encoding/binary"
	"net"

	"github.com/vrgakos/uint128"
)

func ipToInt(ip net.IP) (uint128.Uint128, int) {
	ipv4 := ip.To4()
	if ipv4 != nil {
		return uint128.New(
			uint64(binary.BigEndian.Uint32(ipv4)),
			0,
		), 32
	}

	return uint128.New(
		binary.BigEndian.Uint64(ip[8:]),
		binary.BigEndian.Uint64(ip[:8]),
	), 128
}

func intToIPv4(i uint128.Uint128) net.IP {
	b := make([]byte, 8)
	binary.BigEndian.PutUint64(b, i.Lo)
	return net.IPv4(b[4], b[5], b[6], b[7])
}

func intToIPv6(i uint128.Uint128) net.IP {
	b := make([]byte, 16)
	binary.BigEndian.PutUint64(b[8:], i.Lo)
	binary.BigEndian.PutUint64(b[:8], i.Hi)
	return net.IP(b)
}

func maskToInt(ones, bits int) uint128.Uint128 {
	mask := net.CIDRMask(ones, bits)
	if mask == nil {
		// invalid parameters...
	}

	if len(mask) == 4 {
		return uint128.New(
			uint64(binary.BigEndian.Uint32(mask)),
			0,
		)
	}

	if len(mask) == 16 {
		return uint128.New(
			binary.BigEndian.Uint64(mask[8:]),
			binary.BigEndian.Uint64(mask[:8]),
		)
	}

	return uint128.Zero
}
