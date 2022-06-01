package main

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	_ "net"
	_ "os"
	"log"
	"time"

	_ "golang.org/x/net/icmp"
	_ "golang.org/x/net/ipv4"
)

var (
	device       string = "enp4s0"
	snapshot_len int32  = 1024
	promiscuous  bool   = false
	err          error
	timeout      time.Duration = 1
	handle       *pcap.Handle
)

func main() {

	// Open device
	handle, err = pcap.OpenLive(device, snapshot_len, promiscuous, timeout)
	if err != nil {
	log.Fatal(err)
	}
	defer handle.Close()

	// Set filter capture only icmp
	var filter string = "icmp"
	err = handle.SetBPFFilter(filter)
	if err != nil {
	log.Fatal(err)
	}
	fmt.Println("CAPTURE ICMP PACKETS")

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	fmt.Printf("TIME, FROM, TO\n")
	for packet := range packetSource.Packets() {

		ipLayer:=packet.Layer(layers.LayerTypeIPv4)
		ip, _ :=ipLayer.(*layers.IPv4)
		fmt.Printf("%f,%s,%s\n", float64(time.Now().UnixNano())/1000000000,ip.SrcIP, ip.DstIP)

		applicationLayer := packet.ApplicationLayer()
		if applicationLayer != nil {

			message:=string(applicationLayer.Payload())
			fmt.Println(message)

		}

	}
}

