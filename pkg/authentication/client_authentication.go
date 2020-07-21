package authentication

import (
	"context"
	"google.golang.org/grpc"
	"time"
)

type ClientAuth struct {
	service  AuthServiceClient
	username string
	password string
}

// creates a new ClientAuth
func NewClientAuth(cc *grpc.ClientConn, username string, password string) *ClientAuth {
	service := NewAuthServiceClient(cc)
	return &ClientAuth{service, username, password}
}

// login user and return the access token
func (client *ClientAuth) Login() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &LoginRequest{Username: client.username, Password: client.password}

	// make a request to the login service
	res, err := client.service.Login(ctx, req)
	if err != nil {
		return "", err
	}

	// accessToken (a field from struct LoginResponse) is the message returned from Login service
	return res.GetAccessToken(), nil
}
