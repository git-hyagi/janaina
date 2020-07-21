package authentication

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/git-hyagi/janaina/pkg/types"
	"time"
)

type JWTManager struct {
	secretKey     string
	tokenDuration time.Duration
}

// struct with the jwt standard claims
type UserClaim struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Method   string `json:"method"`
}

// create a jwt manager
func NewJWTManager(secretKey string, tokenDuration time.Duration) *JWTManager {
	return &JWTManager{secretKey, tokenDuration}
}

// Create and sign a new jwt token for user
func (manager *JWTManager) Generate(user *types.Person) (string, error) {

	claims := UserClaim{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(manager.tokenDuration).Unix(),
		},
		Username: user.Name,
		Method:   user.Role,
	}

	// here is where the jwt token is created
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)

	// sign the token with the secretKey
	return token.SignedString([]byte(manager.secretKey))
}

// Verify the token
func (manager *JWTManager) Verify(accessToken string) (*UserClaim, error) {

	// parse the token claim with the token, an empty userclaim and the signature
	// the valid() does not check the content/payload from the token (UserClaim)
	// it only checks if it is not expired and signed by a trusted manager
	token, err := jwt.ParseWithClaims(
		accessToken,
		&UserClaim{},
		func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)
			if !ok {
				return nil, fmt.Errorf("Error with token sign method")
			}
			return []byte(manager.secretKey), nil
		},
	)

	if err != nil {
		return nil, fmt.Errorf("Invalid token: %w", err)
	}

	// convert the claims in a *UserClaim object
	claims, ok := token.Claims.(*UserClaim)
	if !ok {
		return nil, fmt.Errorf("Invalid token claim")
	}

	return claims, nil
}
