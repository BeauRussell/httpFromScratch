package http1

import (
	"bufio"
	"fmt"
	"httpFromScratch/sockets"
	"net"
	"strings"
)

func Server() {
	connection := sockets.TCPConnection{
		Host: "localhost",
		Port: 8080,
	}

	listener := connection.CreateConnection()
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		requestLine, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading request line:", err)
			return
		}
		requestLine = strings.TrimSpace(requestLine)
		fmt.Println("Request Line:", requestLine)

		headers := make(map[string]string)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				fmt.Println("Error reading headers:", err)
			}
			line = strings.TrimSpace(line)
			if line == "" {
				break
			}
			parts := strings.SplitN(line, ":", 2)
			if len(parts) == 2 {
				headers[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
			}
		}

		handleHTTPRequest(conn, requestLine, headers)
	}
}

func handleHTTPRequest(conn net.Conn, requestLine string, headers map[string]string) {
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		writeHTTPResponse(conn, 400, "Bad Request", "Invalid request line")
		return
	}

	method, path, version := parts[0], parts[1], parts[2]
	fmt.Printf("Method: %s, Path: %s, Version %s\n", method, path, version)

	if path == "/" {
		writeHTTPResponse(conn, 200, "OK", "<h1>Welcome to the GO HTTP server!</h1>")
	} else {
		writeHTTPResponse(conn, 404, "Not Found", "<h1>Page not found</h1>")
	}
}

func writeHTTPResponse(conn net.Conn, statusCode int, statusText string, bodyText string) {
	response := fmt.Sprintf(
		"HTTP/1.1 %d %s\r\n"+
			"Content-Length: %d\r\n"+
			"Content-Type: text/html\r\n"+
			"Connection: close\r\n"+
			"\r\n%s",
		statusCode, statusText, len(bodyText), bodyText)
	conn.Write([]byte(response))
}
