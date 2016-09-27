package datastores

import (
	"reflect"
	"testing"

	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/settings"
)

var GlobalUserStore *UserStore

func init() {
	settings.SetPreproductionEnv()
	GlobalUserStore = &UserStore{DB: ConnectToPostgres()}

}

func TestFindUserByID(t *testing.T) {

	expectedUser := &models.User{ID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", HashedPassword: "$2a$10$lWqqb7MhwH7YryO4DyjdeOsFQ9hK7qxZ8PPcm6qjuNlM47KNInHMK"}

	retrievedUser, err, _ := GlobalUserStore.FindUser("id", "0c1b2b91-9164-4d52-87b0-9c4b444ee62d")
	if err != nil {
		t.Error(err)
	}

	if retrievedUser == nil {
		t.Errorf("Expected and did not recieve %#v", expectedUser)
	} else {
		compareUsers(t, expectedUser, retrievedUser)
	}
}

func TestFindUserByUsername(t *testing.T) {

	expectedUser := &models.User{ID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", HashedPassword: "$2a$10$lWqqb7MhwH7YryO4DyjdeOsFQ9hK7qxZ8PPcm6qjuNlM47KNInHMK"}

	retrievedUser, err, _ := GlobalUserStore.FindUser("username", "Tester1")
	if err != nil {
		t.Error(err)
	}

	if retrievedUser == nil {
		t.Errorf("Expected and did not recieve %#v", expectedUser)
	} else {
		compareUsers(t, expectedUser, retrievedUser)
	}
}

func TestStoreUserWithNewCredentials(t *testing.T) {

	err, _ := GlobalUserStore.StoreUser("TestUser", "$2a$10$iziTEDykz1SgOVWhLuBxeeBiZFJdD6GfTO0vA06IJTafiPfSu4QYq")
	if err != nil {
		t.Error(err)
	}
}

func TestStoreUserWithExistingUserCredentials(t *testing.T) {

	err, _ := GlobalUserStore.StoreUser("Tester1", "$2a$10$iziTEEyoz1SgOVWhLuBxeeBiZFJdD6GfTO0vA06IJTafiPfSu4QYq")
	if err.Error() != "The provided username is not unique" {
		t.Error(err)
	}
}

func compareUsers(t *testing.T, x *models.User, y *models.User) {

	// Avoids the complication of parsing postgres timestamp values to golang time.Time
	standardizeTime(&x.CreatedAt, &y.CreatedAt)

	if !reflect.DeepEqual(x, y) {
		t.Errorf("Expected %#v, but recieved %#v", x, y)
	}
}

/*
func TestIsUsernameRegistered(t *testing.T) {
	if !GlobalUserStore.IsUsernameRegistered("Tester2") {
		t.Errorf("IsUsernameUnique function failed to recognize \"Tester2\" as a non unique username")
	}
}
*/
