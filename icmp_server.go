package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)


func main() {

//	host := "192.168.3.21"
	host := "10.0.0.3"

	//Resolve Destination IP DNS
	dst_ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		log.Fatalf("ResolveIPAddr: %v", err)
	}

	//Start Listen ICPM
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()


	//Make ICMP Message
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
		Data: []byte(""),
		},
	}
	wb, err := wm.Marshal(nil)

	for i := 1; i <= 40; i++ {
		//Send ICMP Packet
		start := time.Now()
		n, err := c.WriteTo(wb, dst_ip)
		if err != nil {
			log.Fatalf("WriteTo: %v", err)
		}

		c.SetReadDeadline(time.Now().Add(1 * time.Second))
		duration := time.Since(start)
		reply := make([]byte, 1500)
		n, peer_ip, err := c.ReadFrom(reply)

		icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
		fmt.Println(peer_ip, duration)

		time.Sleep(time.Millisecond * 100)

	}
}
