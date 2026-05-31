package domain

import "fmt"

const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)

type User struct {
	ID           int
	Username     string
	Email        string
	Role         string
	PasswordHash string
}

func NewUser(username, email, role, passwordHash string) *User {
	return &User{
		Username:     username,
		Email:        email,
		Role:         role,
		PasswordHash: passwordHash,
	}
}

func (u *User) String() string {
	return fmt.Sprintf("User(id=%d, username=%s, email=%s, role=%s)", u.ID, u.Username, u.Email, u.Role)
}
