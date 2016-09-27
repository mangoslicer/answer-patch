package datastores

import (
	"database/sql"
	"testing"
)

var GlobalDB *sql.DB

func TestConnectToPostgres(t *testing.T) {
	if err := ConnectToPostgres().Ping(); err != nil {
		t.Error(err)
	}
}
