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


			message:=string(applicationLayer.Payload())
			fmt.Println(message)

		}

	}
}

