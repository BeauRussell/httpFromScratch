package sockets

import (
	"fmt"
	"net"
	"strconv"
)

type socketConnection interface {
	CreateConnection()
}

type TCPConnection struct {
	host     string
	port     int
	listener net.Listener
}

type UDPConnection struct {
	host     string
	port     int
	listener *net.UDPConn
}

func (tcp *TCPConnection) CreateConnection() net.Listener {
	conn, err := net.Listen("tcp", tcp.host+":"+strconv.Itoa(int(tcp.port)))
	if err != nil {
		fmt.Println("Error setting up TCP listener:", err)
		panic(err)
	}

	tcp.listener = conn

	return conn
}

func (udp *UDPConnection) CreateConnection() *net.UDPConn {
	addr := net.UDPAddr{
		Port: udp.port,
		IP:   net.ParseIP(udp.host),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error setting up UDP listener:", err)
		panic(err)
	}

	udp.listener = conn
	return conn
}
