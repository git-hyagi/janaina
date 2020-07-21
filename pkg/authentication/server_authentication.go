package authentication

import (
	"context"
	"github.com/git-hyagi/janaina/pkg/databases"
	config "github.com/git-hyagi/janaina/pkg/loadconfigs"
	"github.com/git-hyagi/janaina/pkg/types"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServer struct {
	jwtManager *JWTManager
}

var database, db_user, db_pass, db_addr, table = config.LoadConfig()
var dsn = db_user + ":" + db_pass + "@tcp(" + db_addr + "/" + database
var connection, _ = databases.Connect(database, db_user, db_pass, db_addr)
var db = &databases.DbConnection{Conn: connection}

// check the user password
func checkPassword(user *types.Person, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// creates a new auth server
func NewAuthServer(jwtManager *JWTManager) *AuthServer {
	return &AuthServer{jwtManager}
}

// Login method from auth service validates user password and return a *LoginResponse (token) if logged in successfully
func (server *AuthServer) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {

	user, err := db.FindPersonByUsername(req.GetUsername())

	// if user not found or error during database query, return err
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Could no find user: %v", err)
	}

	// if password does not match with user, return err
	if !checkPassword(&user, req.GetPassword()) {
		return nil, status.Errorf(codes.Internal, "Incorrect username/password")
	}

	// create a new token for this user
	token, err := server.jwtManager.Generate(&user)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Failed to generate token for user: %v", err)
	}

	res := &LoginResponse{AccessToken: token}
	return res, nil
}
