package authentication

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"time"
)

type ClientAuthInterceptor struct {
	ClientAuth  *ClientAuth
	authMethods map[string]bool
	accessToken string
}

// create a new ClientAuthInterceptor
func NewClientAuthInterceptor(ClientAuth *ClientAuth, authMethods map[string]bool, refreshDuration time.Duration) (*ClientAuthInterceptor, error) {
	interceptor := &ClientAuthInterceptor{ClientAuth: ClientAuth, authMethods: authMethods}
	err := interceptor.scheduleRefreshToken(refreshDuration)
	if err != nil {
		return nil, err
	}

	return interceptor, nil
}

// unary interceptor with the token attached
func (interceptor *ClientAuthInterceptor) Unary() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		log.Printf("----> unary client interceptor: %s", method)

		// if this route/method is an authenticated one attach a token to the invoker (handler)
		if interceptor.authMethods[method] {
			return invoker(interceptor.attachToken(ctx), method, req, reply, cc, opts...)
		}

		// if the route/method is not an authenticated method/route, just return the invoker (handler)
		return invoker(ctx, method, req, reply, cc, opts...)
	}

}

// attach token to the method "authorization" header
func (interceptor *ClientAuthInterceptor) attachToken(ctx context.Context) context.Context {
	return metadata.AppendToOutgoingContext(ctx, "authorization", interceptor.accessToken)
}

// refresh client token in a "refreshDuration" interval
func (interceptor *ClientAuthInterceptor) scheduleRefreshToken(refreshDuration time.Duration) error {

	// before entering in a infine loop through a go routine, ask for a new token
	err := interceptor.refreshToken()
	if err != nil {
		return err
	}

	go func() {
		wait := refreshDuration

		for {
			time.Sleep(wait)                  // wait for refreshDuration until asking for a new token
			err := interceptor.refreshToken() // ask for a new token
			if err != nil {
				wait = time.Second // in the case of an error, wait 1 second before asking for another token
			} else {
				wait = refreshDuration // wait for refreshDuration again before asking for a new token
			}
		}
	}()

	return nil
}

// refresh client token doing a re-login
func (interceptor *ClientAuthInterceptor) refreshToken() error {

	// the Login method returns a new access token
	accessToken, err := interceptor.ClientAuth.Login()

	if err != nil {
		return err
	}
	interceptor.accessToken = accessToken
	log.Printf("token refreshed: %v", accessToken)
	return nil
}
