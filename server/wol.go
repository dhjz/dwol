package server

import (
	"bytes"
	"encoding/hex"
	"net"
	"strings"
)

func SendWOL(mac string, broadcast string, port int) error {
	mac = strings.TrimSpace(mac)
	mac = strings.ReplaceAll(mac, "-", ":")
	mac = strings.ReplaceAll(mac, ".", ":")

	macBytes, err := hex.DecodeString(strings.ReplaceAll(mac, ":", ""))
	if err != nil {
		return err
	}

	if len(macBytes) != 6 {
		return nil
	}

	packet := make([]byte, 102)
	copy(packet, bytes.Repeat([]byte{0xFF}, 6))

	for i := 6; i < 102; i += 6 {
		copy(packet[i:i+6], macBytes)
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   net.ParseIP(broadcast),
		Port: port,
	})
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write(packet)
	return err
}
