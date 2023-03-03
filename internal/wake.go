package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"reflect"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

type Packet interface {
	Send(p Packet)
}

type UDP struct {
	ip      string
	port    string
	src     []byte
	dst     []byte
	options gopacket.SerializeOptions
}

func NewUDP(ip string, port string, src string, dst string) Packet {

	var src_conv = deviceStringToHex(src)
	var dst_conv = deviceStringToHex(dst)

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	return &UDP{
		ip:      ip,
		port:    port,
		src:     src_conv,
		dst:     dst_conv,
		options: opts,
	}

}

/*

The struct UDP implements the Packet Interface

Should Interface Functions recieve values of its own type? (Send(p Packet))

When there is an interface is it suitable to call a types method within a function that accepts an interface? (p.Send())

*/

func Send(p Packet) {

	/* Check network connectivity (interface availability etc.) */

	p.Send(p)
}

/*
How is u able to be referenced if the only parameter was the interface Packet?

Why can p not access the types attributes? (p.src)

How is u able to be referenced if the only parameter to Send() is an interface?
*/

func (u *UDP) Send(p Packet) {

	fmt.Println("Sending Packet ", u, reflect.TypeOf(u.src))

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr(u.src),
		DstMAC: net.HardwareAddr(u.dst),
	}

	buf := gopacket.NewSerializeBuffer()

	var payload = func(dst []byte) []byte {

		var b []byte

		for i := 0; i < 6; i++ {
			b = append(b, uint8(0xff))
		}

		for i := 0; i < 16; i++ {
			b = append(b, u.dst...)
		}

		return b
	}

	pld := payload(u.dst)

	gopacket.SerializeLayers(buf, u.options, &eth, gopacket.Payload(pld))

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%s", u.ip, u.port))
	if err != nil {
		panic(err)
	}

	conn.Write(buf.Bytes())

	fmt.Println("Packet Sent!", buf.Bytes())
}

func deviceStringToHex(s string) []byte {

	var hw []byte

	normal := strings.Split(s, ":")

	for i := range normal {
		h, err := hex.DecodeString(normal[i])
		if err != nil {
			panic(err)
		}

		hw = append(hw, h[0])
	}

	return hw

}

func main() {
	p := NewUDP("255.255.255.0", "0", "28:d0:ea:80:38:9c", "28:d0:ea:80:38:9c")

	Send(p)
}
