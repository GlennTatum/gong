package main

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
)

func NewPacket(addr string, src net.HardwareAddr, dst net.HardwareAddr) (gopacket.SerializeBuffer, error) {
	options := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	eth := layers.Ethernet{
		SrcMAC: net.HardwareAddr(src.String()),
		DstMAC: net.HardwareAddr(dst.String()),
	}

	buf := gopacket.NewSerializeBuffer()
	payload, err := payload(dst)
	if err != nil {
		return nil, err
	}

	gopacket.SerializeLayers(buf, options, &eth, gopacket.Payload(payload))

	return buf, nil
}

func payload(dst net.HardwareAddr) ([]byte, error) {
	var buf []byte
	mac, err := DeviceStringToHex(dst.String())
	if err != nil {
		return nil, err
	}

	for i := 0; i < 6; i++ {
		buf = append(buf, 255)
	}
	for i := 0; i < 16; i++ {
		buf = append(buf, mac...)
	}

	return buf, nil
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

	p, err := NewPacket(addr, src, dst)
	if err != nil {
		return err
	}
	conn.Write(p.Bytes())

	return nil
}

func DeviceStringToHex(s string) ([]byte, error) {

	var hw []byte

	normal := strings.Split(s, ":")

	for i := range normal {
		h, err := hex.DecodeString(normal[i])
		if err != nil {
			return nil, err
		}
		hw = append(hw, h[0])
	}

	return hw, nil
}

func main() {
	// TODO Add flags
	src, err := DeviceStringToHex("7e:33:a1:c8:ad:2a")
	if err != nil {
		os.Exit(1)
	}

	dst, err := DeviceStringToHex("16:03:10:7f:76:76")
	if err != nil {
		os.Exit(1)
	}

	Wake("192.168.76.255", src, dst)
}
