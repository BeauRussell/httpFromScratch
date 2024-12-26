package tlsConfig

import (
	"crypto/tls"
	"fmt"
)

func CreateConfig(certPath string, keyPath string, protos string) (*tls.Config, error) {
	config := &tls.Config{
		Certificates: []tls.Certificate{},
		NextProtos:   []string{},
	}
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, fmt.Errorf("Failed to load TLS certificate: %v", err)
	}

	config.Certificates = append(config.Certificates, cert)
	config.NextProtos = append(config.NextProtos, protos)

	return config, nil
}
