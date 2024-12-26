package sockets

import (
	"crypto/tls"
	"fmt"
	"httpFromScratch/tlsConfig"
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

func (tcp *TCPConnection) CreateConnection(protos string) net.Listener {
	tlsConfig, err := tlsConfig.CreateConfig("/home/beau/.tls/cert.pem", "/home/beau/.tls/key.pem", protos)
	if err != nil {
		fmt.Println("Cannot set up tcp. Error with TLS Config:", err)
		panic(err)
	}
	conn, err := tls.Listen("tcp", tcp.Host+":"+strconv.Itoa(int(tcp.Port)), tlsConfig)
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
