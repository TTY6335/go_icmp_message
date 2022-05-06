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

func send_icmp(dst_ip *net.IPAddr,IcmpPcaket *icmp.PacketConn, host_ip string,send_data string){
	fmt.Printf("%s %s\n",host_ip,send_data)

	//Make ICMP Message
	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho,
		Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff,
			Data: []byte(send_data),
		},
	}
	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatalf("Marshal: %v", err)
	}


	//Send ICMP Packet
	start := time.Now()
	_,_ = IcmpPcaket.WriteTo(wb, dst_ip)
	if err != nil {
		log.Fatalf("WriteTo: %v", err)
	}

	IcmpPcaket.SetReadDeadline(time.Now().Add(10 * time.Second))
	reply := make([]byte, 1500)
	_, peer_ip, err := IcmpPcaket.ReadFrom(reply)
	duration := time.Since(start)

	fmt.Println(peer_ip,start,duration,start.Add(duration))
}

func main() {

	host := "10.0.0.1"

	//Start Listen ICPM
	c, err := icmp.ListenPacket("ip4:icmp", "0.0.0.0")
	if err != nil {
		log.Fatalf("ListenPacket: %v", err)
	}
	defer c.Close()

	//Resolve Destination IP DNS
	dst_ip, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		log.Fatalf("ResolveIPAddr: %v", err)
	}

	for i := 1; i <= 20; i++ {

		send_icmp(dst_ip,c,host,"192.168.3.20")

	}
}
