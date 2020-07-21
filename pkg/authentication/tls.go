package authentication

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
)

// load server crt and key files
func LoadServerTLS(serverCrt, serverKey string) (credentials.TransportCredentials, error) {

	// Load private/public key pair
	certs, err := tls.LoadX509KeyPair(serverCrt, serverKey)
	if err != nil {
		return nil, err
	}

	// struct to config the TLS server
	config := &tls.Config{
		Certificates: []tls.Certificate{certs},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

// load server CA to be trusted by client
func LoadServerCA(caFile string) (credentials.TransportCredentials, error) {

	// Load CA from server's certificate
	serverCA, err := ioutil.ReadFile(caFile)
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(serverCA) {
		return nil, fmt.Errorf("Failed to append cert to the cert pool")
	}

	// struct to config the TLS
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil
}
