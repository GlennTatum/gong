package main

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
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

	fmt.Println("Using Interface: ", *iface_in)
	fmt.Println("Sending Magic Packet to Desktop: ", *hardwareaddr_in)

	in, err := net.InterfaceByName(*iface_in)
	if err != nil {
		fmt.Println("e")
		panic(err)
	}

	hw_addr := in.HardwareAddr.String()

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr(hw_addr),
		DstMAC: net.HardwareAddr(*hardwareaddr_in),
	}

	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	p := generateWOLPayload(*hardwareaddr_in)

	gopacket.SerializeLayers(buf, opts, &eth, gopacket.Payload(p))

	conn, err := net.Dial("udp", "255.255.255.0:9")
	if err != nil {
		panic(err)
	}

	conn.Write(buf.Bytes())

}

func generateWOLPayload(hardwareaddr string) []byte {

	macaddr := strings.Split(hardwareaddr, ":")

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

	fmt.Println(payload)

	return payload

}
