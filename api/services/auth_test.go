package services

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/mangoslicer/answer-patch/settings"
)

/**
 * Mock of the datastore package's TokenStore struct
 * The StoreTime field represents the amount of time for which a token of a logged-out user should be stored
*/
type MockTokenStore struct {
	StoreTime int
}

/**
 * Mock AuthContext available for testing functions, specifically the unit tests for the Login method
*/
var (
	userID = "0" // Arbitrary userID for testing purposes
	ac = &AuthContext{UserID: userID}
)

func init() {
	settings.SetPreproductionEnv()
}

/**
 * Mock implementation of StoreToken method from the TokenStoreServices interface
 * Records the exp parameter for later use by TestLogoutStoreTimeCalculation
*/
func (store *MockTokenStore) StoreToken(key, val string, exp int) error {
	store.StoreTime = exp
	return nil
}

/**
 * Mock implementation of StoreToken method from the TokenStoreServices interface
 * IsTokenStored always returns false because
*/
func (store *MockTokenStore) IsTokenStored(key string) (bool, error) {
	return false, nil
}

/**
 * Tests that Login rejects passwords which do not meet the password requirements defined in auth.go
*/
func TestLoginWithInvalidPassword(t *testing.T) {
	invalidPasswordMessage := "Login did not detect that %s is an invalid password was provided"
  var invalidPasswords = []string {
				"", // Does not pass the password length requirement
				"no_capital_letters", // Lacks special character and capital letter
				"Capital", //Lacks special character
				"lowercase!", // Lacks capital letter
	}
	for _, invalidPassword := range invalidPasswords {
	 		_, err := ac.Login(invalidPassword, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") //Arbitrary hash for testing purposes
 			if err.Error() != "Invalid password recieved"  {
 				t.Errorf(invalidPasswordMessage, invalidPassword)
 		}
	}
}

/**
 * Tests that Login accepts passwords which do meet the password requirements defined in auth.go
*/
func TestLoginWithValidPassword(t *testing.T) {
	validPasswordMessage := "Login interpreted %s as an invalid password"
	var validPasswords = []string {
			"Abcde$",
			"O'reillyAutoParts",
	}
 for _, validPassword := range validPasswords {
		_, err := ac.Login(validPassword, "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") //Arbitrary hash for testing purposes
		if err.Error() == "Invalid password recieved"  {
			t.Errorf(validPasswordMessage, validPassword)
	}
 }
}

/**
 * Tests that the Login method correctly detects incorrect password and does not return a JSON Web Token
*/
func TestLoginWithIncorrectPassword(t *testing.T) {
	_, err := ac.Login("incorrect password", "2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824") // The password string does not match the hash string
	if err == nil {
		t.Errorf("Login did not detect that the provided password does not match the hashedpassword of the existing user.")
	}
}

/**
 * Tests that the Login method correctly detects password and does return a JSON Web Token
 * Parses JSON Web Token to check whether the subscriber of the token matches the user ID of the user who requested the token
*/
func TestLoginWithCorrectPassword(t *testing.T) {
 hashedPassword := "$2a$10$XgxVoZidTuxugFcAkBipBeSYxFJSLv/w0t1Lt7ihTOu0cThCBHgA." // Hash of "Passw!rd" string literal

/*
  userID :="0"
	ac := &AuthContext{UserID: userID}
*/
	retrievedToken, err := ac.Login("Passw!rd", hashedPassword)
	if err != nil {
		t.Error(err)
	} else if retrievedToken == nil {
		t.Errorf("Correct password was provided, but the Login method did not return a JSON Web Token")
	}

	parsedToken, err := jwt.Parse(retrievedToken.SignedToken, func(token *jwt.Token) (interface{}, error) {
		return settings.GetPublicKey(), nil
	})
	if err != nil {
		t.Error(err)
	}

	parsedTokenClaims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.Errorf("Error reading parsed JSON Web Token")
	}

	if parsedTokenClaims["sub"] != ac.UserID {
		t.Errorf("Expected to retrieve a token with a \"sub\" claim of %s, but the retrieved token has a \"sub\" claim of %s", ac.UserID, parsedTokenClaims["sub"])
	}
}

/**
 * Tests whether the Logout method correctly calculates the time for which the token of a logged-out user should be stored
*/
func TestLogoutStoreTimeCalculation(t *testing.T) {

	mockTokenStore := &MockTokenStore{StoreTime: 0}
	exp := time.Now().Add(time.Second * time.Duration(5)) // Time value 5 seconds ahead of current time
	logoutAccountContext := &AuthContext{Exp: exp, TokenStore: mockTokenStore}

	signedToken := "WIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ" // Arbitrary signed JSON Web Token
	err := logoutAccountContext.Logout(signedToken) // The Logout method's calculated storetime is recorded in mockTokenStore's StoreTime field by MockTokenStore's StoreToken method
	if err != nil {
		t.Error(err)
	}

	maxBoundOnStoreTimeCalc := exp.Sub(time.Now()).Seconds() // Represents the maximum amount of error (~5 seconds) which is allowed by the Logout storetime calculation
	if mockTokenStore.StoreTime > (int(maxBoundOnStoreTimeCalc) + StoreOffset)  {
		t.Errorf("Expected the calculated store time of the logged-out user's token to be less five seconds, but the store time was %d", mockTokenStore.StoreTime)
	}
}
