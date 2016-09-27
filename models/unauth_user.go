package models

import (
	"log"

	"code.google.com/p/go.crypto/bcrypt"
)

type UnauthUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

const (
	hashCost = 10
)

func NewUnauthUser() *UnauthUser {
	return &UnauthUser{}
}

func (unauth *UnauthUser) GetMissingFields() string {

	var missing string

	switch {
	case unauth.Username == "":
		missing = "username\n"
	case unauth.Password == "":
		missing += "password\n"
	}

	return missing
}

func (unauth *UnauthUser) HashPassword() string {
	hash, err := bcrypt.GenerateFromPassword([]byte(unauth.Password), hashCost)
	if err != nil {
		log.Fatal(err)
	}

	return string(hash)
}
