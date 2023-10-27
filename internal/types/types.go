package types

import (
	"golang.org/x/crypto/bcrypt"
	"time"
)

type User struct {
	ID                int       `json:"ID"`
	FirstName         string    `json:"name"`
	LastName          string    `json:"lastName"`
	MemberSince       time.Time `json:"memberSince"`
	EncryptedPassword string    `json:"encryptedPassword"`
}

type UserDto struct {
	ID          int       `json:"ID"`
	FirstName   string    `json:"name"`
	LastName    string    `json:"lastName"`
	MemberSince time.Time `json:"memberSince"`
}

type Account struct {
	AccountNumber int       `json:"accountNumber"`
	Balance       float32   `json:"balance"`
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
}

func (u *User) ValidPassword(pw string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(pw)) == nil
}

func NewUser(FirstName string, LastName string, Password string) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         FirstName,
		LastName:          LastName,
		EncryptedPassword: string(encpw),
		MemberSince:       time.Now(),
	}, nil
}

func NewAccount(OwnerID int) *Account {
	return &Account{
		Balance: 0,
		OwnerID: OwnerID,
		Created: time.Now(),
	}
}
