package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func NewPacket(addr string, src net.HardwareAddr, dst net.HardwareAddr) gopacket.SerializeBuffer {
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr(src.String()),
		DstMAC: net.HardwareAddr(dst.String()),
	}

	buf := gopacket.NewSerializeBuffer()

	gopacket.SerializeLayers(buf, options, &eth, gopacket.Payload(payload(dst)))

	return buf
}

func payload(dst net.HardwareAddr) []byte {
	var buf []byte

	for i := 0; i < 6; i++ {
		buf = append(buf, 255)
	}
	for i := 0; i < 16; i++ {
		buf = append(buf, DeviceStringToHex(dst.String())...)
	}

	return buf
}

func dial(addr string) (net.Conn, error) {
	conn, err := net.Dial("udp", fmt.Sprintf("%s:0", addr))
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func Wake(addr string, src net.HardwareAddr, dst net.HardwareAddr) error {
	conn, err := dial(addr)
	if err != nil {
		return err
	}

	p := NewPacket(addr, src, dst)
	conn.Write(p.Bytes())

	return nil
}

func DeviceStringToHex(s string) []byte {

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
	// TODO Add flags
	src := DeviceStringToHex("7e:33:a1:c8:ad:2a")
	dst := DeviceStringToHex("16:03:10:7f:76:76")
	Wake("192.168.76.255", src, dst)
}
