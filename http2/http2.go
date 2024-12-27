package http2

import (
	"fmt"
	"httpFromScratch/sockets"
	"net"
)

const (
	http2Preface = "PRI * HTTP/2.0\r\n\r\nSM\r\n\r\n"
)

func Server() {
	connection := sockets.TCPConnection{
		Host: "localhost",
		Port: 8081,
	}

	listener := connection.CreateConnection("h2")
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Failed to accept new connection:", err)
			return
		}

		go handleHTTPConnection(conn)
	}
}

func handleHTTPConnection(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, len(http2Preface))
	_, err := conn.Read(buf)
	if err != nil || string(buf) != http2Preface {
		fmt.Println("Invalid HTTP/2 preface:", err)
		return
	}

	fmt.Println("Received HTTP/2 preface")

	// Create HTTP/2 settings frame
	settingsFrame := []byte{
		0x00, 0x00, 0x00, //Length (0 currently)
		0x04,                   // Type: SETTINGS
		0x00,                   // Flags
		0x00, 0x00, 0x00, 0x00, // Stream Identifier
	}
	_, err = conn.Write(settingsFrame)
	if err != nil {
		fmt.Println("Failed to send settings frame:", err)
		return
	}

	// Simple HTTP setup
	for {
		frameHeader := make([]byte, 9)
		_, err := conn.Read(frameHeader)
		if err != nil {
			fmt.Println("Error reading frame header:", err)
			return
		}

		// Parse frame header
		length := int(frameHeader[0])<<16 | int(frameHeader[1])<<8 | int(frameHeader[2])
		frameType := frameHeader[3]
		flags := frameHeader[4]
		streamId := int(frameHeader[5])<<24 | int(frameHeader[6])<<16 | int(frameHeader[7])<<8 | int(frameHeader[8])&0x7FFFFFFF

		fmt.Printf("Received frame: type=%d, flags=%d, streamID=%d, length=%d\n", frameType, flags, streamId, length)

		//Read the payload
		payload := make([]byte, length)
		_, err = conn.Read(payload)
		if err != nil {
			fmt.Println("error reading frame payload:", err)
			return
		}

		// Process frames (simplified)
		switch frameType {
		case 0x01: // HEADERS
			fmt.Println("Received HEADERS frame")
		case 0x00: // DATA
			fmt.Println("Received DATA frame")
			fmt.Println("Payload:", string(payload))
		default:
			fmt.Printf("Unknown frame type: %d\n", frameType)
		}
	}
}
