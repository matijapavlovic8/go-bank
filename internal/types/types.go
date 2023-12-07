package types

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type Role string

const (
	AdminRole Role = "admin"
	UserRole  Role = "user"
)

type User struct {
	ID                int       `json:"ID"`
	FirstName         string    `json:"name"`
	LastName          string    `json:"lastName"`
	MemberSince       time.Time `json:"memberSince"`
	EncryptedPassword string    `json:"encryptedPassword"`
	Role              Role      `json:"role"`
}

type UserDto struct {
	ID          int       `json:"ID"`
	FirstName   string    `json:"name"`
	LastName    string    `json:"lastName"`
	MemberSince time.Time `json:"memberSince"`
	Role        Role      `json:"role"`
}

type Account struct {
	AccountNumber int       `json:"accountNumber"`
	Balance       float64   `json:"balance"`
	OwnerID       int       `json:"ownerID"`
	Created       time.Time `json:"created"`
}

type LoginRequest struct {
	ID       int    `json:"id"`
	Password string `json:"password"`
}

type LoginResponse struct {
	ID    int    `json:"id"`
	Token string `json:"token"`
}

type CreateUserRequest struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Password  string `json:"password"`
	Role      string `json:"role"`
}

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}

func NewUser(FirstName string, LastName string, Password string, Role string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	role, err := StringToRole(Role)
	if err != nil {
		return nil, err
	}
	return &User{
		FirstName:         FirstName,
		LastName:          LastName,
		EncryptedPassword: string(encpw),
		MemberSince:       time.Now(),
		Role:              role,
	}, nil
}

func NewAccount(OwnerID int) *Account {
	return &Account{
		Balance: 0,
		OwnerID: OwnerID,
		Created: time.Now(),
	}
}

func StringToRole(s string) (Role, error) {
	switch s {
	case string(AdminRole):
		return AdminRole, nil
	case string(UserRole):
		return UserRole, nil
	default:
		return "", errors.New("Invalid role")
	}
}
