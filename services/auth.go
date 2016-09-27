package services

import (
	"errors"
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/patelndipen/AP1/datastores"
	"github.com/patelndipen/AP1/settings"
	"golang.org/x/crypto/bcrypt"
)

const (
	JWTLife     = 72
	StoreOffset = 60
)

type Token struct {
	SignedToken string `json:"token"`
}

type AuthServices interface {
	Login(string, string) (*Token, error)
	Logout(string) error
	RefreshToken() (*Token, error)
}

type AuthContext struct {
	UserID     string
	Exp        time.Time
	TokenStore datastores.TokenStoreServices
}

var InternalErr = errors.New("Internal Error")

func NewAuthContext(ts datastores.TokenStoreServices) *AuthContext {
	return &AuthContext{UserID: "", Exp: time.Time{}, TokenStore: ts}
}

func (ac *AuthContext) Login(enteredPassword, hashedPassword string) (*Token, error) {
	if !authenticate(hashedPassword, enteredPassword) {
		return nil, errors.New("Credentials are incorrect")
	}

	return setTokenClaims(ac.UserID)
}

func (ac *AuthContext) Logout(signedToken string) error {

	storeTime := int(ac.Exp.Sub(time.Now()).Seconds() + StoreOffset)

	err := ac.TokenStore.StoreToken(ac.UserID, signedToken, storeTime)
	if err != nil {
		return InternalErr
	}

	return nil
}

func (ac *AuthContext) RefreshToken() (*Token, error) {
	return setTokenClaims(ac.UserID)
}

func authenticate(hashedPassword, enteredPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword)) == nil
	return true

}

func setTokenClaims(userID string) (*Token, error) {
	token := jwt.New(jwt.SigningMethodRS512)
	token.Claims["iat"] = time.Now().Unix()
	token.Claims["exp"] = time.Now().Add(time.Hour * time.Duration(JWTLife)).Unix()
	token.Claims["sub"] = userID

	signedToken, err := token.SignedString(settings.GetPrivateKey())
	if err != nil {
		log.Fatal(err)
		return nil, InternalErr
	}
	return &Token{SignedToken: signedToken}, nil
}
