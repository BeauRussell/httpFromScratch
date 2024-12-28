package http2

import (
	"fmt"
	"httpFromScratch/framePackaging"
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
	fmt.Println(string(buf))
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

	for {
		frame := &framePackaging.Frame{}
		err := frame.ParseFrame(conn)
		if err != nil {
			fmt.Println("Error parsing frame:", err)
			return
		}

		switch frame.Type {
		case framePackaging.FrameData:
			fmt.Println("Got frame data, building response:", string(frame.Payload))
			sendDataPackets(conn, frame)
		case framePackaging.FrameHeaders:
			fmt.Println("HEADERS frame received")
		case framePackaging.FrameSettings:
			fmt.Println("SETTINGS frame received")
		default:
			fmt.Printf("Unknown frame type: %d\n", frame.Type)
		}
	}
}

func sendDataPackets(conn net.Conn, receivedFrame *framePackaging.Frame) {
	headerFrame := &framePackaging.Frame{}
	headers := map[string]string{
		":status":      "200",
		"content-type": "text/plain",
	}
	headerPacket, err := headerFrame.BuildHeadersFrame(headers, receivedFrame.StreamID)
	fmt.Println(headerPacket)
	if err != nil {
		fmt.Println("Failed to build header packet", err)
		return
	}
	_, err = conn.Write(headerPacket)
	if err != nil {
		fmt.Println("Failed to send header packet", err)
		return
	}

	dataFrame := &framePackaging.Frame{}
	dataPacket, err := dataFrame.BuildDataFrame(string(receivedFrame.Payload), receivedFrame.StreamID)
	if err != nil {
		fmt.Println("Failed to build data packet:", err)
		return
	}
	fmt.Println(dataPacket)
	_, err = conn.Write(dataPacket)
	if err != nil {
		fmt.Println("Failed to send data packet:", err)
		return
	}
}
