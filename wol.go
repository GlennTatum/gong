package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func main() {

	// Get a list of all interfaces

	/*
	 * Command Options:
	 * -i : Interface to be used
	 * -h : Hardware address to be found
	 *
	 */

	iface_in := flag.String("i", "Interface", "Network Interface e.g. enp5s0")
	hardwareaddr_in := flag.String("d", "Mac Address", "Network Hardware Address e.g. 01:23:45:67:89:ab")

	flag.Parse()

	fmt.Println("Using Interface ", *iface_in)
	fmt.Println("Sending Magic Packet to Desktop ", *hardwareaddr_in)

	// Scan for the interface

	// Get the IPNet and HardwareAddr Values

	iface, err := net.InterfaceByName(*iface_in)
	if err != nil {

		// Interface is not found return failure

		panic(err)
	}

	ifaces, _ := iface.Addrs()

	_, net, err := net.ParseCIDR(ifaces[0].String())
	if err != nil {
		panic(err)
	}

	send(net, &iface.HardwareAddr)

	// Write Magic Packet
}

func send(iface *net.IPNet, hardwareaddr *net.HardwareAddr) {

	handle, err := pcap.OpenLive(iface.String(), 9, true, pcap.BlockForever)
	if err != nil {
		panic(err)
	}

	writeWOL(handle, iface, hardwareaddr)

	defer handle.Close()

}

func writeWOL(handle *pcap.Handle, iface *net.IPNet, hardwareaddr *net.HardwareAddr) {

	// https://en.wikipedia.org/wiki/Wake-on-LAN

	// https://pkg.go.dev/github.com/google/gopacket#hdr-Creating_Packet_Data

	eth := layers.Ethernet{
		SrcMAC:       *hardwareaddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}

	udp := layers.UDP{

		SrcPort:  9,
		DstPort:  9,
		Length:   102,
		Checksum: 0,
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{}

	payload := generateWOLPayload(hardwareaddr)

	gopacket.SerializeLayers(buf, opts,
		&eth,
		&udp,
		gopacket.Payload(payload),
	)

	handle.WritePacketData(buf.Bytes())

}

func generateWOLPayload(hardwareaddr *net.HardwareAddr) []byte {

	macaddr := strings.Split(hardwareaddr.String(), ":")

	fmt.Println(macaddr)

	payload := []byte{

		// Broadcast MAC Address

		uint8(0xff),
		uint8(0xff),
		uint8(0xff),
		uint8(0xff),
		uint8(0xff),
		uint8(0xff),
	}

	// 16 Repetitions of the target computers MAC Address

	for i := 0; i < 16; i++ {

		for i := range macaddr {

			toHex, _ := hex.DecodeString(macaddr[i])

			payload = append(payload, toHex...)
		}

	}

	return payload

}
