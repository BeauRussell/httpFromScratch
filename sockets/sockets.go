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
	Host     string
	Port     int
	Listener net.Listener
}

type UDPConnection struct {
	Host     string
	Port     int
	Listener *net.UDPConn
}

func (tcp *TCPConnection) CreateConnection() net.Listener {
	conn, err := net.Listen("tcp", tcp.Host+":"+strconv.Itoa(int(tcp.Port)))
	if err != nil {
		fmt.Println("Error setting up TCP listener:", err)
		panic(err)
	}

	tcp.Listener = conn

	return conn
}

func (udp *UDPConnection) CreateConnection() *net.UDPConn {
	addr := net.UDPAddr{
		Port: udp.Port,
		IP:   net.ParseIP(udp.Host),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		fmt.Println("Error setting up UDP listener:", err)
		panic(err)
	}

	udp.Listener = conn
	return conn
}
