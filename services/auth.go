package services

import (
	"errors"
	"log"
	"time"
	"regexp"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mangoslicer/answer-patch/datastores"
	"github.com/mangoslicer/answer-patch/settings"
	"golang.org/x/crypto/bcrypt"
)

/**
 * Constants for JSON Web Token configuration
*/
const (
	JWTLife     = 72 // Represents the amount of hours until a JSON Web Token will expire
	StoreOffset = 60 // Represents the amount of seconds beyond the token's expiration time that the token will be stored as invalid
)

/**
 * Constants for password validity verification
*/
const (
	PasswordLength = 6
)

type Token struct {
	SignedToken string `json:"token"`
}

type AuthServices interface {
	Login(string, string) (*Token, error)
	Logout(string) error
	RefreshToken() (*Token, error)
}

/**
 * AuthContext is a structure that is used by the request handlers to access data about the current user
 * Upon login, the user is ...
*/
type AuthContext struct {
	UserID     string
	Exp        time.Time
	TokenStore datastores.TokenStoreServices
}

/**
 * Creates a new AuthContext struct with the provided TokenStore
*/
func NewAuthContext(ts datastores.TokenStoreServices) *AuthContext {
	return &AuthContext{UserID: "", Exp: time.Time{}, TokenStore: ts}
}

/**
 * Provides a JSON Web Token for authenticated users
*/
func (ac *AuthContext) Login(enteredPassword, hashedPassword string) (*Token, error) {
	if !isValidPassword(enteredPassword) {
		return nil, errors.New("Invalid password recieved")
	} else if !authenticate(hashedPassword, enteredPassword) {
		return nil, errors.New("Credentials are incorrect")
	}

	return setTokenClaims(ac.UserID)
}

/**
 * The token parameter is stored in the TokenStore for an amount of hours determined by the sum of the remaining token life and the Offset
 * Tokens of logged-out users are stored until at least expiration such that the server can detect and block any HTTP requests with these tokens
*/
func (ac *AuthContext) Logout(signedToken string) error {

	storeTime := int(ac.Exp.Sub(time.Now()).Seconds() + StoreOffset)

	// Checking if storeTime is greater than 0 prevents the storage of expired JSON Web Tokens
	if storeTime > 0 {
		err := ac.TokenStore.StoreToken(ac.UserID, signedToken, storeTime)
		if err != nil {
			return errors.New("Internal Error")
		}
	}

	return nil
}

/**
 * Rather than parsing, delaying the expiration time, and signing the current token, RefreshToken simply returns a new token
*/
func (ac *AuthContext) RefreshToken() (*Token, error) {
	return setTokenClaims(ac.UserID)
}

/**
 * Authentication of user is performed by comparing the hash of the entered password with the stored hash of the user's password
*/
func authenticate(hashedPassword, enteredPassword string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(enteredPassword)) == nil
}

/**
 * Initializes JSON Web Token with initialization time, expiration time, and the userID
*/
func setTokenClaims(userID string) (*Token, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.MapClaims {
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * time.Duration(JWTLife)).Unix(),
		"sub": userID,
  })

	signedToken, err := token.SignedString(settings.GetPrivateKey())
	if err != nil {
		log.Fatal(err)
		return nil, errors.New("Internal Error")
	}
	return &Token{SignedToken: signedToken}, nil
}

/**
 * Checks whether password strings pass the character requirements of valid passwords
*/
func isValidPassword(password string) bool {
	capitalLetterRegex, _ := regexp.Compile("[A-Z]+")  // Checks whether the password contains at least one capital letter
	specialCharacterRegex, _ := regexp.Compile("\\W+") // Checks whether the password contains at least one special character
	return len(password) >= PasswordLength && capitalLetterRegex.MatchString(password) && specialCharacterRegex.MatchString(password)
}
