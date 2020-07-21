# token.go
- basically creates a new token encrypted with a `secretKey` and with `tokenDuration` duration
- the token payload will be filled with:
  - a username
  - a method (thinking in changing this to another name, maybe route)
  - some standard jwt fields, overwriting only the duration


# client_authentication.go (service AuthServiceClient Login method implementation)
- Login request
- it expects an username and password
- it receives a token from `LoginResponse`(server)
- the struct has a service (protobuf service) a username(string) and a password(string)

# server_authentication.go (service AuthServiceServer Login method implementation)
- Login response to the `LoginRequest`(client) service with a token
- makes a query to the database with the username provided
- check if the user returned from the database has the same password (hashed) as the one provided during loging
- if the password matches, generates a new token using `jwtManager.Generate` method
- return the token

# client_auth_interceptor.go (attach a token and an "authorization" header to the handler in the case of an authenticated method )
- Returns a `unary interceptor` for the methods (routes) with a token in case of an authenticated route
- the authMethods field is only responsible to route the request to an authenticated or not path
- the `unary interceptor` will just passthrough the request if this is not an authenticated method and will append the **token** and an "**authorization**" header in case of an authenticated route/method
- when the `clientauthinterceptor` is called, the "constructor" asks for a new token

# server_auth_interceptor.go ()
- The `server interceptor` will receive a token and a role. The method will verify if:
  - the metadata has an "**authorization**" header
     - if so, verify if the token provided is valid (through the `jwtManager.Verify` method) 
  - the role present in the token matches the one came from the metadata requisition
- after doing the above checks (for the authenticated methods) just returns the handler
