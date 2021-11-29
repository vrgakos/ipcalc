package main

import (
	"fmt"

	"github.com/vrgakos/ipcalc"
)

func main() {
	baseSub1 := ipcalc.NewSubnet("192.168.100.0/24")
	baseSub1.Meta = "base"
	fmt.Println(baseSub1.Print())
	fmt.Println()

	/*rand.Seed(time.Now().UnixNano())
	var err error
	var i uint8
	i = 0
	for {
		fmt.Scanln()

		hostSub := fmt.Sprintf("192.168.100.%d/32", rand.Int()%256)
		i++

		fmt.Print("Push ", hostSub, " ")
		_, err = baseSub1.Insert(ipcalc.NewSubnet(hostSub))
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("SIKER")
		//fmt.Println(baseSub1.Print())
		//fmt.Println()

	}*/

	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.141/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.74/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.10/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.234/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.238/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.129/32"))
	baseSub1.Insert(ipcalc.NewSubnet("192.168.100.226/32"))
	fmt.Println(baseSub1.Insert(ipcalc.NewSubnet("192.168.100.128/25")))
	fmt.Println(baseSub1.Insert(ipcalc.NewSubnet("192.168.100.224/27")))

	fmt.Println(baseSub1.Print())

	var i uint8 = 0
	for {
		var str string
		fmt.Scanln(&str)

		fmt.Print("Find ")
		fmt.Println(baseSub1.Find(ipcalc.NewSubnet(str)))

		fmt.Print("Lookup ")
		fmt.Println(baseSub1.Lookup(ipcalc.NewSubnet(str)))

		i++
	}

	/*
		baseSub2 := ipcalc.NewSubnet("192.168.100.0/24")
		err = baseSub2.Insert(baseSub1)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(baseSub2.Print())
	*/

	/*ranger := cidranger.NewPCTrieRanger()
	_, network1, _ := net.ParseCIDR("192.168.1.0/24")
	_, network2, _ := net.ParseCIDR("128.168.1.0/32")

	ranger.Insert(cidranger.NewBasicRangerEntry(*network1))
	ranger.Insert(cidranger.NewBasicRangerEntry(*network2))*/
}
