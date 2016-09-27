package datastores

import (
	"log"
	"testing"

	"github.com/patelndipen/AP1/settings"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func init() {

	settings.SetPreproductionEnv()
}

func TestConnectToMongoCol(t *testing.T) {

	col := ConnectToMongoCol()

	if col == nil {
		t.Errorf("Expected a *mgo.Collection, but recieved nil")
	}

}

func populateMongoCol(col *mgo.Collection) {

	_, err := col.RemoveAll(bson.M{})
	if err != nil {
		log.Fatal(err)
	}

	err = col.Insert(bson.M{"_id": bson.M{"category": "testing", "userID": "0"}, "rep": 5})
	if err != nil {
		log.Fatal(err)
	}

	err = col.Insert(bson.M{"_id": bson.M{"category": "testing", "userID": "1"}, "rep": 5})
	if err != nil {
		log.Fatal(err)
	}

}
