package ipcalc

import (
	"fmt"
	"net"

	"github.com/vrgakos/uint128"
)

type Range struct {
	// Subnet *Subnet
	Bits  int
	Start uint128.Uint128
	End   uint128.Uint128
	Size  uint128.Uint128
}

func NewRange(start, end net.IP) *Range {
	res := &Range{}
	var startBits, endBits int

	res.Start, startBits = ipToInt(start)
	res.End, endBits = ipToInt(end)

	if startBits != endBits {
		return nil
	}
	res.Bits = startBits

	if res.Start.Cmp(res.End) >= 0 {
		return nil
	}

	res.Size = res.End.Sub(res.Start).Add64(1)

	return res
}

func ParseRange(start, end string) (*Range, error) {
	startIp := net.ParseIP(start)
	if startIp == nil {
		return nil, fmt.Errorf("could not parse start ip")
	}

	endIp := net.ParseIP(end)
	if endIp == nil {
		return nil, fmt.Errorf("could not parse end ip")
	}

	return NewRange(startIp, endIp), nil
}

func (r *Range) GetIpByOffset64(offset uint64) net.IP {
	if r.Size.Cmp64(offset) <= 0 {
		return nil
	}

	if r.Bits == 32 {
		return intToIPv4(r.Start.Add64(offset))
	}

	if r.Bits == 128 {
		return intToIPv6(r.Start.Add64(offset))
	}

	return nil
}

func (r *Range) GetIpByOffset(offset uint128.Uint128) net.IP {
	if r.Size.Cmp(offset) <= 0 {
		return nil
	}

	if r.Bits == 32 {
		return intToIPv4(r.Start.Add(offset))
	}

	if r.Bits == 128 {
		return intToIPv6(r.Start.Add(offset))
	}

	return nil
}
