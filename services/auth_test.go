package services

import (
	"testing"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/patelndipen/AP1/settings"
)

type MockTokenStore struct {
	StoreTime int
}

func init() {
	settings.SetPreproductionEnv()
}

func (store *MockTokenStore) StoreToken(key, val string, exp int) error {
	store.StoreTime = exp //Records the exp param with which Logout called StoreToken in order for the TestLogoutStoreTimeCalculation func to test whether Logout calculated the correct StoreTime
	return nil
}

func (store *MockTokenStore) IsTokenStored(key string) (bool, error) {
	return false, nil
}

func TestLoginWithIncorrectPassword(t *testing.T) {

	ac := new(AuthContext)

	_, err := ac.Login("incorrect password", "Hash")
	if err == nil {
		t.Errorf("Login did not detect that the provided password does not match the hashedpassword of the existing user.")
	}
}

func TestLoginWithCorrectPassword(t *testing.T) {

	ac := &AuthContext{UserID: "0"}
	retrievedToken, err := ac.Login("password", "$2a$10$F4B95tmhW6VfQ.l.mhUj6Ow3Eg5dViTiJDsFzh8VaQr9Urd70LP9W")
	if err != nil {
		t.Error(err)
	}

	parsedToken, err := jwt.Parse(retrievedToken.SignedToken, func(token *jwt.Token) (interface{}, error) {
		return settings.GetPublicKey(), nil
	})
	if err != nil {
		t.Error(err)
	}

	parsedTokenUserID, ok := parsedToken.Claims["sub"].(string)
	if !ok {
		t.Errorf("The parsedToken \"sub\" claim does not have an underlying type of string")
	}

	if parsedTokenUserID != ac.UserID {
		t.Errorf("Expected to retrieve a token with a \"sub\" claim of %s, but the retrieved token has a \"sub\" claim of %s", ac.UserID, parsedTokenUserID)
	}
}

func TestLogoutStoreTimeCalculation(t *testing.T) {

	mockTokenStore := &MockTokenStore{StoreTime: 0}
	ac := &AuthContext{Exp: time.Now(), TokenStore: mockTokenStore}

	err := ac.Logout("Signed Token") // The Logout method's calculated storetime is recorded in mockTokenStore's StoreTime field by MockTokenStore's StoreToken method
	if err != nil {
		t.Error(err)
	}

	if mockTokenStore.StoreTime < (StoreOffset-5) && mockTokenStore.StoreTime != StoreOffset {
		t.Errorf("Expected the calculated store time of the token to be either equal to StoreOffset or greater than StoreOffset - 5, because the AuthContext struct, with which the Logout method was called, contained an Exp field of value time.Now(), but the calculated store time passed to StoreToken by the Logout method was %d seconds", mockTokenStore.StoreTime)

	}
}
