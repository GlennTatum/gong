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
	Build() gopacket.SerializeBuffer
	Write(buf gopacket.SerializeBuffer)
}

type UDP struct {
	ip      string
	src     string
	dst     string
	options gopacket.SerializeOptions
}

type NetworkError struct {
	err error
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

func (u *UDP) Write(buf gopacket.SerializeBuffer) {

	// Check Network Connectivity

	conn, err := net.Dial("udp", fmt.Sprintf("%s:0", u.ip))
	if err != nil {
		panic(err)
	}

	conn.Write(buf.Bytes())
}

func (u *UDP) Build() gopacket.SerializeBuffer {

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

	return buf
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

// This should return some form of an implemented error interface

func checkConn() string {
	return "Network is not properly configured"
}

func main() {
	p := NewUDP("192.168.76.255", "b4:2e:99:3c:39:c9", "18:C0:4D:36:EE:91")

	p.Write(p.Build())
}
