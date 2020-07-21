package main

import (
	"github.com/git-hyagi/janaina/pkg/authentication"
	"github.com/git-hyagi/janaina/pkg/person"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

const (
	secretKey     = "mySuperSecretPassword"
	tokenDuration = 60 * time.Minute
	serverCrt     = "deploy/server.crt"
	serverKey     = "deploy/server.key"
)

func methods() map[string][]string {
	return map[string][]string{
		"/person.ManagePerson/CreateMedic":    {"admin"},
		"/person.ManagePerson/CreatePatient":  {"admin"},
		"/person.ManagePerson/CreateAdmin":    {"admin"},
		"/person.ManagePerson/RemovePerson":   {"admin"},
		"/person.ManagePerson/FindPerson":     {"admin"},
		"/person.ManagePerson/UpdateEmail":    {"admin"},
		"/person.ManagePerson/UpdatePhone":    {"admin"},
		"/person.ManagePerson/UpdateFacebook": {"admin"},
		"/person.ManagePerson/UpdateName":     {"admin"},
		"/person.ManagePerson/UpdateSurname":  {"admin"},
		"/person.ManagePerson/UpdatePassword": {"admin"},
	}
}

func main() {

	// auxiliar struct to make grpc methods requests
	managePerson := person.Server{}

	// token
	jwtManager := authentication.NewJWTManager(secretKey, tokenDuration)
	authServer := authentication.NewAuthServer(jwtManager)

	// tls
	tlsFiles, err := authentication.LoadServerTLS(serverCrt, serverKey)
	if err != nil {
		log.Fatal("Error loading TLS files: %v", err)
	}

	// create a listener on port 9000
	listener, err := net.Listen("tcp", ":9000")

	if err != nil {
		log.Fatal("Failed to bind port 9000: %v", err)
	}

	// authentication interceptor
	interceptor := authentication.NewServerAuthInterceptor(jwtManager, methods())

	// start grpc server
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(interceptor.Unary()), grpc.Creds(tlsFiles))
	authentication.RegisterAuthServiceServer(grpcServer, authServer)
	person.RegisterManagePersonServer(grpcServer, &managePerson)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal("Fatal to serve gRPC server over port 9000: %v", err)
	}
}
