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

	ip, _ := in.Addrs()

	ipv4_local := strings.Split(ip[0].String(), ".")

	ipv4_local[3] = "255"

	ipv4_broadcast := strings.Join(ipv4_local, ".")

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

	// https://forums.ivanti.com/s/article/Understanding-Wake-On-LAN?language=en_US

	conn, err := net.Dial("udp", fmt.Sprintf("%s:0", ipv4_broadcast))
	if err != nil {
		panic(err)
	}

	conn.Write(buf.Bytes())

}

func generateWOLPayload(hardwareaddr string) []byte {

	macaddr := strings.Split(hardwareaddr, ":")

	payload := []byte{uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff)}

	for i := 0; i < 16; i++ {

		for i := range macaddr {

			toHex, _ := hex.DecodeString(macaddr[i])

			payload = append(payload, toHex...)
		}

	}

	return payload

}
