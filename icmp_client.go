package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"net"
	"os"
	"log"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
)

var (
	device       string = "eth0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 1
	handle       *pcap.Handle
)

func send_icmp(dst_ip *net.IPAddr,IcmpPcaket *icmp.PacketConn, host_ip string,send_data string){

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
	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
	log.Fatal(err)
	}
	defer handle.Close()

	// Set filter
	var filter string = "icmp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("CAPTURE ICMP PACKETS")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for packet := range packetSource.Packets() {

		ipLayer:=packet.Layer(layers.LayerTypeIPv4)
		ip, _ :=ipLayer.(*layers.IPv4)
		fmt.Printf("%v From %s to %s\n", time.Now(),ip.SrcIP, ip.DstIP)

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {


			host:=string(applicationLayer.Payload())

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

			send_icmp(dst_ip,c,host,"192.168.3.20")


		}

	}
}

