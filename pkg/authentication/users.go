package authentication

import (
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// user account
type User struct {
	Name           string
	HashedPassword string
	Role           string
}

// create new user
func Newuser(username, password, role string) (*User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("Error trying to hash provided password: %w", err)
	}

	user := &User{Name: username, HashedPassword: string(hashedPassword), Role: role}
	return user, nil
}

// check the user password
func (user *User) checkPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	return err == nil
}

// clone user
func (user *User) cloneUser() *User {
	return &User{
		Name:           user.Name,
		HashedPassword: user.HashedPassword,
		Role:           user.Role,
	}
}
