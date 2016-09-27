package datastores

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/patelndipen/AP1/models"
)

type UserStoreServices interface {
	FindUser(string, string) (*models.User, error, int)
	StoreUser(string, string) (error, int)
	//	IsUsernameRegistered(string) (bool, error, int)
}

type UserStore struct {
	DB *sql.DB
}

func (store *UserStore) FindUser(filter, searchVal string) (*models.User, error, int) {

	queryStmt := `SELECT id, username, hashed_password, created_at FROM  ap_user WHERE ` + filter + ` =$1`

	row, err := store.DB.Query(queryStmt, searchVal)
	if err != nil {
		log.Fatal(err)
		return nil, InternalErr, http.StatusInternalServerError
	} else if !row.Next() {
		return nil, errors.New("No user exists with the provided credential"), http.StatusOK
	}

	user := new(models.User)

	err = row.Scan(&user.ID, &user.Username, &user.HashedPassword, &user.CreatedAt)
	if err != nil {
		log.Fatal(err)
		return nil, InternalErr, http.StatusInternalServerError
	}

	return user, nil, http.StatusOK

}

func (store *UserStore) StoreUser(username, hashedpassword string) (error, int) {
	/*
		row, err := store.DB.Query(`SELECT id FROM ap_user WHERE username = $1 AND hashed_password = $2`, username, hashedpassword)
		if err != nil {
			log.Fatal(err)
			return err, http.StatusInternalServerError
		} else if row.Next() {
			return nil, http.StatusOK
		}
	*/
	return transact(store.DB, func(tx *sql.Tx) (error, int) {
		_, err := tx.Exec(`INSERT INTO ap_user(username, hashed_password) values($1, $2)`, username, hashedpassword)
		if err != nil {
			return evaluateSQLError(err)
		}
		return nil, http.StatusOK
	})
}

/*
func (store *UserStore) IsUsernameRegistered(username string) bool {

	row, err := store.DB.Query(`SELECT username FROM ap_user WHERE username = $1`, username)
	if err != nil {
		log.Fatal(err)
	}

	return row.Next()
}
*/
