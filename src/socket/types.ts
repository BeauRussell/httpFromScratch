interface SocketAddress {
	sa_family: number;
	sa_data: string;
}

interface AddressInfo {
	ai_flags: number;
	ai_family: number;
	ai_socktype: number;
	ai_protocol: number;
	ai_adrlen: number;
	ai_cannonname: string;

	ai_addr: SocketAddress;
}

interface SocketAddressIn {
	s_addr: number;
}

interface SocketAddressIn4 {
	sin4_family: number;
	sin4_port: number;
	sin4_addr: SocketAddressIn;
	sin4_zero: string;
}


interface SocketAddressIn6 {
	sin6_family: number;
	sin6_port: number;
	sin6_flowinfo: number;
	sin6_scope_id: string;
	sin6_addr:  SocketAddressIn;
}

export { AddressInfo, SocketAddress, SocketAddressIn, SocketAddressIn4, SocketAddressIn6 };
