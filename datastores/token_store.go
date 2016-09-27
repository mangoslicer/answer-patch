package datastores

import (
	"log"

	"github.com/garyburd/redigo/redis"
)

type TokenStoreServices interface {
	StoreToken(string, string, int) error
	IsTokenStored(string) (bool, error)
}

type JWTStore struct {
	Conn redis.Conn
}

func (store *JWTStore) StoreToken(userID, signedToken string, exp int) error {
	_, err := store.Conn.Do("SET", userID, signedToken)
	if err != nil {
		log.Fatal(err)
		return InternalErr
	}

	_, err = store.Conn.Do("EXPIRE", userID, exp)
	if err != nil {
		log.Fatal(err)
		return InternalErr
	}

	return nil
}

func (store *JWTStore) IsTokenStored(userID string) (bool, error) {

	val, err := store.Conn.Do("GET", userID)
	if err != nil {
		log.Fatal(err)
		return false, InternalErr
	} else if val == nil {
		return false, nil
	}

	return true, nil
}
