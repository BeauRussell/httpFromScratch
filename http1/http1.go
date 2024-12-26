package http1

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"httpFromScratch/sockets"
	"io"
	"net"
	"os"
	"strconv"
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
			if !errors.Is(err, io.EOF) {
				fmt.Println("Error reading request line:", err)
			}
			return
		}
		requestLine = strings.TrimSpace(requestLine)

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

		contentLength := 0
		if val, ok := headers["Content-Length"]; ok {
			contentLength, err = strconv.Atoi(val)
			if err != nil {
				fmt.Println("Invalid content length:", val)
				return
			}
		}

		body := make([]byte, contentLength)
		_, err = io.ReadFull(reader, body)
		if err != nil {
			fmt.Println("Error reading body:", err)
			return
		}

		handleHTTPRequest(conn, requestLine, headers, body)
	}
}

func handleHTTPRequest(conn net.Conn, requestLine string, headers map[string]string, body []byte) {
	fmt.Println(headers)
	parts := strings.Split(requestLine, " ")
	if len(parts) < 3 {
		writeHTTPResponse(conn, 400, "Bad Request", "Invalid request line")
		return
	}

	method, path, version := parts[0], parts[1], parts[2]
	fmt.Printf("Method: %s, Path: %s, Version %s\n", method, path, version)

	if path == "/" {
		htmlString, err := loadHTML("templates/index.html")
		if err != nil {
			fmt.Println("Cannot load HTML file to send to connection:", err)
			return
		}
		writeHTTPResponse(conn, 200, "OK", htmlString)
	} else if path == "/post" && method == "POST" {
		htmlString, err := loadHTML("templates/post.html")
		if err != nil {
			fmt.Println("Cannot load HTML file to send to connection:", err)
			return
		}
		htmlString = handlePostData(htmlString, headers, body)
		writeHTTPResponse(conn, 200, "OK", htmlString)
	} else {
		htmlString, err := loadHTML("templates/404.html")
		if err != nil {
			fmt.Println("Cannot load HTML file to send to connection:", err)
			return
		}
		writeHTTPResponse(conn, 404, "Not Found", htmlString)
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

func loadHTML(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Printf("File to open file: %v\n", err)
		return "", err
	}
	defer file.Close()

	var htmlString string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		htmlString += scanner.Text() + "\n"
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return "", err
	}

	return htmlString, nil
}

func handlePostData(html string, headers map[string]string, body []byte) string {
	dataMap := map[string]string{
		"Title":   "Testing Post",
		"Content": string(body),
	}

	var buf bytes.Buffer

	switch headers["Content-Type"] {
	case "application/json":
		tmpl := template.Must(template.New("webpage").Parse(html))
		err := (tmpl.Execute(&buf, dataMap))
		if err != nil {
			fmt.Println("Failed to write to HTML template:", err)
			return ""
		}
	}

	return buf.String()
}
