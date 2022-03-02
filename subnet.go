package ipcalc

import (
	"fmt"
	"net"
	"strings"

	"github.com/vrgakos/uint128"
)

type Subnet struct {
	isIPv6  bool            // true: this is an IPv6 address
	NetInt  uint128.Uint128 // network stored as a number
	MaskInt uint128.Uint128 // mask stored as a number (1 where network is "fixed")
	NetOnes uint8           // mask size

	Meta string
	// SubnetCidr string
	// Used       bool

	// TREE links
	isDummy bool
	parent  *Subnet
	// skippedBits uint8
	children []*Subnet
}

func NewSubnet(cidr string) *Subnet {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		return nil
	}

	ipIntBase, _ := ipToInt(net.IP)
	ones, bits := net.Mask.Size()
	subnet := &Subnet{
		isIPv6:   bits == 128,
		NetInt:   ipIntBase,
		NetOnes:  uint8(ones),
		children: make([]*Subnet, 2),
	}
	subnet.MaskInt = subnet.calcMaskInt()

	return subnet
}

func (s *Subnet) CloneBase() *Subnet {
	subnet := &Subnet{
		isIPv6:   s.isIPv6,
		NetInt:   s.NetInt,
		MaskInt:  s.MaskInt,
		NetOnes:  s.NetOnes,
		children: make([]*Subnet, 2),
	}
	return subnet
}

func (s *Subnet) CloneWithOnes(ones uint8) *Subnet {
	res := s.CloneBase()
	res.NetOnes = ones
	res.MaskInt = res.calcMaskInt()
	res.NetInt = res.NetInt.And(res.MaskInt)

	return res
}

func (s *Subnet) SameSubnet(s2 *Subnet) bool {
	return s.isIPv6 == s2.isIPv6 && s.NetOnes == s2.NetOnes && s.NetInt.Cmp(s2.NetInt) == 0
}

func (s *Subnet) Find(f *Subnet) (*Subnet, error) {
	if !s.Intersect(f) {
		return nil, fmt.Errorf("the find subnet have to be intersected")
	}

	child := s

	for child != nil {
		if f.SameSubnet(child) && !child.isDummy {
			return child, nil
		}

		bitPos := child.targetBitPosition()
		bitVal := f.bitValue(uint8(bitPos))

		child = child.children[bitVal]
	}

	return nil, fmt.Errorf("not found")
}

func (s *Subnet) Lookup(f *Subnet) (*Subnet, error) {
	if !s.Intersect(f) {
		return nil, fmt.Errorf("the lookup subnet have to be intersected")
	}

	best := s
	child := s

	for child != nil {
		bitPos := child.targetBitPosition()
		bitVal := f.bitValue(uint8(bitPos))

		child = child.children[bitVal]
		if child.Contains(f) && !child.isDummy {
			best = child
		}
	}

	return best, nil
}

func (s *Subnet) Insert(newChild *Subnet) (bool, error) {
	// fmt.Printf("Inserting %s into %s\n", newChild, s)
	// fmt.Println("----------------------------------------------------------")
	// fmt.Println("s.Print")
	// fmt.Println(s.Print())
	// fmt.Println("----------------------------------------------------------")

	if !s.Intersect(newChild) {
		return false, fmt.Errorf("the inserted subnet have to be intersected")
	}

	if newChild.NetOnes == s.NetOnes {
		if s.isDummy {
			s.isDummy = false
			return true, nil
		}
		return false, fmt.Errorf("already there")
	}

	if newChild.NetOnes < s.NetOnes {
		return false, fmt.Errorf("the inserted subnet is not smaller than the base")
	}

	bitPos := s.targetBitPosition()
	bitVal := newChild.bitValue(uint8(bitPos))
	// fmt.Printf(" %d. bitPos of newChild network is %d\n", bitPos, bitVal)

	existingChild := s.children[bitVal]
	// fmt.Printf(" existingChild is %v\n", existingChild)

	// No existing child, in this direction, just add a new one.
	if existingChild == nil {
		s.children[bitVal] = newChild
		newChild.parent = s
		return true, nil
	}

	commonOnes := existingChild.CommonOnes(newChild, false)
	// fmt.Printf(" existingChild and newChild has %d common ones\n", commonOnes)

	//divergingBitPos := commonOnes //- 1
	// fmt.Printf(" divergingBitPos=%d existingChild.targetBitPosition()=%d\n", divergingBitPos, existingChild.targetBitPosition())

	if existingChild.NetOnes > commonOnes {
		/*var asd uint8 = 32
		if newChild.isIPv6 {
			asd = 128
		}*/
		dummySubnet := newChild.CloneWithOnes(existingChild.CommonOnes(newChild, true))
		dummySubnet.isDummy = true
		fmt.Printf(" insert dummySubnet=%s\n", dummySubnet)

		// Place dummySubnet
		s.children[bitVal] = dummySubnet
		dummySubnet.parent = s

		// Relocate existingChild
		asd := dummySubnet.targetBitPosition()
		// fmt.Printf(" dummySubnet targetBitPosition=%d\n", asd)
		// fmt.Printf(" dummySubnet (uint8)targetBitPosition=%d\n", uint8(asd))
		dummySubnetBitValue := existingChild.bitValue(uint8(asd)) // NOT SURE, but looks good
		// fmt.Printf(" dummySubnet dummySubnetBitValue=%d\n", dummySubnetBitValue)
		dummySubnet.children[dummySubnetBitValue] = existingChild
		existingChild.parent = dummySubnet

		existingChild = dummySubnet
		// fmt.Println("----------------------------------------------------------")
		// fmt.Println(dummySubnet.Print())
		// fmt.Println("----------------------------------------------------------")
	}

	return existingChild.Insert(newChild)
}

func (s *Subnet) Intersect(s2 *Subnet) bool {
	if s == nil {
		return false
	}

	if s2 == nil {
		return false
	}

	return s.NetInt.And(s.MaskInt).Cmp(s2.NetInt.And(s2.MaskInt)) == 0 ||
		s.NetInt.And(s2.MaskInt).Cmp(s2.NetInt.And(s.MaskInt)) == 0
}

func (s *Subnet) Contains(s2 *Subnet) bool {
	if s == nil {
		return false
	}

	if s2 == nil {
		return false
	}

	return s.NetInt.And(s.MaskInt).Cmp(s2.NetInt.And(s.MaskInt)) == 0
}

func (s *Subnet) GetNetwork() net.IP {
	if s.isIPv6 {
		return intToIPv6(s.NetInt)
	}

	return intToIPv4(s.NetInt)
}

func (s *Subnet) GetNetworkStr() string {
	return s.GetNetwork().String()
}

//
func (s *Subnet) CommonOnes(s2 *Subnet, checkOnes bool) uint8 {
	// fmt.Printf(" CommonOnes %s %s\n", s, s2)
	// fmt.Print(s.DebugString())
	// fmt.Print(s2.DebugString())

	res := uint8(s.NetInt.Xor(s2.NetInt).LeadingZeros())
	if !s.isIPv6 {
		res -= 128 - 32
	}
	if !checkOnes {
		return res
	}

	// TODO: maybe replace with bit magic operations
	if s.NetOnes < res {
		res = s.NetOnes
	}
	if s2.NetOnes < res {
		res = s2.NetOnes
	}
	return res
}

func (s *Subnet) GetCidr() string {
	return fmt.Sprintf("%s/%d", s.GetNetwork(), s.NetOnes)
}

func (s *Subnet) DebugString() string {
	return fmt.Sprintf("%s\n net= %b\n mask=%b\n ones=%d\n\n", s.GetCidr(), s.NetInt.Big(), s.MaskInt.Big(), s.NetOnes)
}

func (s *Subnet) totalNumberOfBits() uint8 {
	if s.isIPv6 {
		return 128
	}

	return 32
}

func (s *Subnet) calcMaskInt() uint128.Uint128 {
	return maskToInt(int(s.NetOnes), int(s.totalNumberOfBits()))
}

func (s *Subnet) targetBitPosition() uint8 {
	return s.totalNumberOfBits() - s.NetOnes
}

func (s *Subnet) targetBitValue() uint8 {
	if s.GetBit(uint8(s.targetBitPosition())) {
		return 1
	}

	return 0
}

func (s *Subnet) bitValue(bitPos uint8) uint8 {
	if s.GetBit(bitPos) {
		return 1
	}

	return 0
}

func (s *Subnet) String() string {
	isDummyStr := ""
	if s.isDummy {
		isDummyStr = ",isDummy"
	}
	return fmt.Sprintf("%s (ones:%d,targetPos:%d%s)", s.GetCidr(), s.NetOnes, s.targetBitPosition(), isDummyStr)
}

func (s *Subnet) Print() string {
	children := []string{}
	padding := strings.Repeat("| ", s.level()+1)

	for bitVal, child := range s.children {
		if child == nil {
			continue
		}
		childStr := fmt.Sprintf("\n%s%d --> %s", padding, bitVal, child.Print())
		children = append(children, childStr)
	}

	return fmt.Sprintf("%s%s", s.String(), strings.Join(children, ""))
}

func (s *Subnet) GetBit(bit uint8) bool {
	return s.NetInt.GetBit(bit - 1)
}

func (s *Subnet) SetBit(bit uint8) {
	s.NetInt = s.NetInt.SetBit(bit - 1)
}

func (s *Subnet) ClearBit(bit uint8) {
	s.NetInt = s.NetInt.ClearBit(bit - 1)
}

func (s *Subnet) level() int {
	if s.parent == nil {
		return 0
	}
	return s.parent.level() + 1
}

func (s *Subnet) IsIPv6() bool {
	return s.isIPv6
}

func (s *Subnet) GetVersion() int8 {
	if s.isIPv6 {
		return 6
	} else {
		return 4
	}
}

func (s *Subnet) IsHostAddress() bool {
	return (s.isIPv6 && s.NetOnes == 128) || (!s.isIPv6 && s.NetOnes == 32)
}