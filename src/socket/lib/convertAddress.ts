
// These conversions are very much unnecessary as node connections at the best level doesn't need this. But it helps the learning process


/**
	* Takes in IPv4 Address and converts it to binary for socket use
	* @param {string} ip - IPv4 Address
	* @returns {Uint8Array} - Binary representation of the IPv4 Address
*/
function ip4ToBinary(ip: string): Uint8Array{
	const octets = ip.split('.');
	let binary = 0;

	if (octets.length !== 4) {
		throw new Error('Invalid IPv4 Address');
	}

	const binaryIp  = new Uint8Array(4);
	for (let i = 0; i < octets.length; i++) {
		const octect = parseInt(octets[i], 10);
		if (isNaN(octect) ||octect < 0 || octect > 255) {
			throw new Error('Invalid IPv4 Address');
		}

		// Storing a number into a Uint8Array auto converts it to binary
		binaryIp[i] = binary;
	}

	return binaryIp;
}

/**
	*  Takes in IPv6 Address and converts it to binary for socket use
	*  @param {string} ip - IPv6 Address
	*  @returns {Uint8Array} - Binary representation of the IPv6 Address
*/
function ip6ToBinary(ip: string): Uint8Array {
	const hextets: string[] = ip.split(':');
	if (hextets.length !== 8) {
		throw new Error('Invalid IPv6 Address');
	}
	
	const bytes = [];

	for (const hextet of hextets) {
		const binary = parseInt(hextet, 16);
		if (isNaN(binary) || binary < 0 || binary > 65535) {
			throw new Error('Invalid IPv6 Address');
		}

		bytes.push(binary >> 8);
		bytes.push(binary & 0xFF);
	}

	return new Uint8Array(bytes);
}

/**
	* Takes in IPv4 Address and converts it to binary for socket use
	* @param {Uint8Array} binaries - IPv4 Address in Binary Array
	* @returns {string} - presentation representation of the IPv4 Address
*/
function binaryToIp4(binaries: Uint8Array): string {
	const octets = new Array(4);
	for (let i = 0; i < binaries.length; i++) {
		// uInt8Array stores numbers in binary, but will output them as base 10
		octets[i] = binaries[i].toString();
	}

	if (octets.length !== 4) {
		throw new Error('Invalid IPv4 Address');
	}

	return octets.join('.');
}

/**
	* Takes in IPv4 Address and converts it to binary for socket use
	* @param {Uint8Array} binaries - IPv6 address in Binary Array
	* @returns {string} - presentation representation of the IPv4 Address
*/
function binaryToIp6(binaries: Uint8Array): string {
	const hextets = new Array(8);
	for (let i = 0; i < binaries.length; i += 2) {
		// uInt8Array stores numbers in binary, but will output them as base 10
		hextets[i / 2] = (binaries[i] << 8 | binaries[i + 1]).toString(16);
	}

	if (hextets.length !== 8) {
		throw new Error('Invalid IPv6 Address');
	}

	return hextets.join(':');
}

function ipToBinary(version: string, ip: string) {
	if (version === 'ipv4') {
		return ip4ToBinary(ip);
	} else if (version === 'ipv6') {
		return ip6ToBinary(ip);
	}
}

function binaryToIp(version: string, binary: Uint8Array) {
	if (version === 'ipv4') {
		return binaryToIp4(binary);
	} else if (version === 'ipv6') {
		return binaryToIp6(binary);
	}
}

export { ipToBinary, binaryToIp };

