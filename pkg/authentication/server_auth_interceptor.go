package authentication

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"log"
)

type ServerAuthInterceptor struct {
	jwtManager *JWTManager
	methods    map[string][]string
}

// returns a new serverAuthInterceptor
func NewServerAuthInterceptor(jwtManager *JWTManager, methods map[string][]string) *ServerAuthInterceptor {
	return &ServerAuthInterceptor{jwtManager, methods}
}

// returns a unary interceptor to authenticate unary rpc
func (interceptor *ServerAuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		log.Printf("---> unary interceptor: ", info.FullMethod)

		// call the authorized method to verify the token passed
		err = interceptor.authorize(ctx, info.FullMethod)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// verify the user's token
func (interceptor *ServerAuthInterceptor) authorize(ctx context.Context, method string) error {

	// if the method has no "route" associated, it means that there is no role for it
	// so, just return to the caller
	methods, ok := interceptor.methods[method]
	if !ok {
		return nil
	}

	md, ok := metadata.FromIncomingContext(ctx)
	log.Printf("metdata incoming from client: %v", md)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "metadata not provided")
	}

	// the token value comes from metadata["authorization"] header
	values := md["authorization"]

	// if the metadata has no authorization header it means that no token was provided
	if len(values) == 0 {
		return status.Errorf(codes.Unauthenticated, "authorization token not provided!")
	}

	// verify if the token is valid
	accessToken := values[0]
	claims, err := interceptor.jwtManager.Verify(accessToken)
	if err != nil {
		return status.Errorf(codes.Unauthenticated, "access token is invalid: %v", err)
	}

	// for each role defined for this method
	for _, role := range methods {

		// if a role matches with one comes from the token's payload
		// means that this token has access to this route/method
		// so, it is ok to proceed (just return nil)
		if role == claims.Method {
			return nil
		}
	}

	return status.Error(codes.PermissionDenied, "no permission to access this RPC")
}
