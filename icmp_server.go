package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

func main() {

	host := "10.0.0.9"

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

	for i := 1; i <= 10; i++ {
		//Make ICMP Message
		wm := icmp.Message{
			Type: ipv4.ICMPTypeEcho,
			Code: 0,
			Body: &icmp.Echo{
				ID: os.Getpid() & 0xffff, Seq: i,
				Data: []byte("WORLD"),
			},
		}
		wb, err := wm.Marshal(nil)
		if err != nil {
			log.Fatalf("Marshal: %v", err)
		}
		//Send ICMP Packet
		start := time.Now()
		n, err := c.WriteTo(wb, dst_ip)
		if err != nil {
			log.Fatalf("WriteTo: %v", err)
		}

		c.SetReadDeadline(time.Now().Add(10 * time.Second))
		reply := make([]byte, 1500)
		n, peer_ip, err := c.ReadFrom(reply)
		duration := time.Since(start)

		rm, err := icmp.ParseMessage(ipv4.ICMPTypeEcho.Protocol(), reply[:n])
		fmt.Println(rm.Body, peer_ip, duration)
	}
}
