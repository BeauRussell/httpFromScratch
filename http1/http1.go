package http1

import "httpFromScratch/sockets"

func Server() {
	connection := sockets.TCPConnection{
		Host: "localhost",
		Port: 8080,
	}

	listener := connection.CreateConnection()

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
}
