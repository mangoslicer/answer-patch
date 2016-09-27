package datastores

import (
	"fmt"
	"log"

	"github.com/patelndipen/AP1/settings"
	"gopkg.in/mgo.v2"
)

func ConnectToMongoCol() *mgo.Collection {

	dsn := settings.GetMongoDSN()

	s, err := mgo.Dial(fmt.Sprintf("mongodb://%s:%s@localhost:%s/", dsn.Username, dsn.Password, dsn.Addr))
	if err != nil {
		log.Fatal(err)
	}

	s.SetMode(mgo.Monotonic, true)

	return s.DB(dsn.DBName).C(dsn.ColName)
}
