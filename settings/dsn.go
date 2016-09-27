package settings

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type PostgresDSN struct {
	Username string
	Password string
	DBName   string
	SSLMode  string
}

type RedisDSN struct {
	Addr     string
	Password string
}

type MongoDSN struct {
	Username string
	Password string
	Addr     string
	DBName   string
	ColName  string
}

func GetPostgresDSN() string {

	dsn := new(PostgresDSN)

	getDSN("postgres", dsn)

	return fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s", dsn.Username, dsn.Password, dsn.DBName, dsn.SSLMode)

}

func GetRedisDSN() *RedisDSN {

	dsn := new(RedisDSN)

	getDSN("redis", dsn)

	return dsn
}

func GetMongoDSN() *MongoDSN {

	dsn := new(MongoDSN)

	getDSN("mongo", dsn)

	return dsn

}

func getDSN(dbName string, dsnStruct interface{}) {

	content, err := ioutil.ReadFile("/home/dipen/go/src/github.com/patelndipen/AP1/settings/" + os.Getenv("GO_ENV") + "/" + dbName + ".json")
	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(content, dsnStruct)
	if err != nil {
		log.Fatal(err)
	}

}
