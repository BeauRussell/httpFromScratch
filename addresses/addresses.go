package addresses

import (
	"encoding/binary"
	"net"
)

type addrinfo struct {
	ai_flags     int
	ai_family    int
	ai_socktype  int
	ai_protocol  int
	ai_canonname string
	ai_addr      *sockaddr
	ai_next      *addrinfo
}

type sockaddr struct {
	sa_family uint16
	sa_data   [14]byte
}

type in_addr struct {
	s_addr uint32
}

type sockaddr_in struct {
	sin_family int16
	sin_port   uint16
	sin_addr   in_addr
	sin_zero   [8]byte
}

type sockaddr_in6 struct {
	sin6_family   uint16
	sin6_port     uint16
	sin6_flowinfo uint32
	sin6_scope_id uint32
}

type in6_addr struct {
	s6_addr [16]byte
}

func (a *addrinfo) htons(value uint16) uint16 {
	var b [2]byte
	binary.BigEndian.PutUint16(b[:], value)
	return binary.BigEndian.Uint16(b[:])
}

func (a *addrinfo) ntohs(value uint16) uint16 {
	return a.htons(value)
}

func (a *addrinfo) htonl(value uint32) uint32 {
	var b [4]byte
	binary.BigEndian.PutUint32(b[:], value)
	return binary.BigEndian.Uint32(b[:])
}

func (a *addrinfo) ntohl(value uint32) uint32 {
	return a.htonl(value)
}

func (a *addrinfo) inet_pton4(ip4 string) net.IP {
	return net.ParseIP(ip4).To4()
}

func (a *addrinfo) inet_pton6(ip6 string) net.IP {
	return net.ParseIP(ip6).To16()
}

func (a *addrinfo) inet_ntop4(ip4 [4]byte) string {
	return net.IP(ip4[:]).String()
}

func (a *addrinfo) inet_ntop6(ip6 [16]byte) string {
	return net.IP(ip6[:]).String()
}
