package datastores

import (
	"testing"

	"github.com/patelndipen/AP1/settings"
)

var GlobalRepStore *RepStore

func init() {

	settings.SetPreproductionEnv()
	GlobalRepStore = &RepStore{ConnectToMongoCol()}

	populateMongoCol(GlobalRepStore.Col)
}

func TestFindRep(t *testing.T) {

	expectedRep := 5 //the populateMongoCol set the key of {category:"testing", userID:"0"} to a rep of 5

	retrievedRep, err := GlobalRepStore.FindRep("testing", "0")
	if err != nil {
		t.Error(err)
	}

	if expectedRep != retrievedRep {
		t.Errorf("Expected the rep of {category:\"testing\", userID:\"0\"} to be 5, but the FindRep method returned %d", retrievedRep)
	}
}

func TestUpdateRep(t *testing.T) {

	expectedRep := 6

	err := GlobalRepStore.UpdateRep("testing", "1", 1)
	if err != nil {
		t.Error(err)
	}

	retrievedRep, err := GlobalRepStore.FindRep("testing", "1")
	if err != nil {
		t.Error(err)
	}

	if expectedRep != retrievedRep {
		t.Errorf("Expected the rep of {category:\"testing\", userID:\"1\"} to be 6, but the FindRep method returned %d", retrievedRep)
	}
}

func TestUpdateRepWithNewUserID(t *testing.T) {

	expectedRep := 10

	err := GlobalRepStore.UpdateRep("testing", "2", 5) // The userID of 2 does not yet exist in the collection, therefore UpdateRep should insert the userID into the collection
	if err != nil {
		t.Error(err)
	}

	retrievedRep, err := GlobalRepStore.FindRep("testing", "2")
	if err != nil {
		t.Error(err)
	}

	if expectedRep != retrievedRep {
		t.Errorf("Expected the rep of {category:\"testing\", userID:\"1\"} to be 10, but the FindRep method returned %d", retrievedRep)
	}

}
