package main

import(
	"fmt"
	"net"
)

func main() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, inter := range interfaces {
		addrs, err := inter.Addrs()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(addrs)
//		for _, a := range addrs {
//			if ipnet, ok := a.(*net.IPNet); ok {
//				if ipnet.IP.To4() != nil {
//					fmt.Println(ipnet.IP.String())
//				}
//			}
//		}
	}
}
