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
	Send()
}

type UDP struct {
	ip      string
	src     string
	dst     string
	options gopacket.SerializeOptions
}

func NewUDP(ip string, src string, dst string) Packet {

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	return &UDP{
		ip:      ip,
		src:     src,
		dst:     dst,
		options: opts,
	}

}

/*

The struct UDP implements the Packet Interface

Should Interface Functions recieve values of its own type? (Send(p Packet))

When there is an interface is it suitable to call a types method within a function that accepts an interface? (p.Send())

*/

func Write(p Packet) {
	// "Generic" that checks net
}

/*
How is u able to be referenced if the only parameter was the interface Packet?

Why can p not access the types attributes? (p.src)

How is u able to be referenced if the only parameter to Send() is an interface?
*/

func (u *UDP) Send() {

	fmt.Println("Sending Packet ", u, reflect.TypeOf(u.src))

	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr(u.src),
		DstMAC: net.HardwareAddr(u.dst),
	}

	buf := gopacket.NewSerializeBuffer()

	payload := []byte{uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff), uint8(0xff)}

	for i := 0; i < 16; i++ {
		payload = append(payload, deviceStringToHex(u.dst)...)
	}

	gopacket.SerializeLayers(buf, u.options, &eth, gopacket.Payload(payload))

	conn, err := net.Dial("udp", fmt.Sprintf("%s:0", u.ip))
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
	p := NewUDP("192.168.76.255", "28:d0:ea:80:38:9c", "18:C0:4D:36:EE:91")

	p.Send()
}
