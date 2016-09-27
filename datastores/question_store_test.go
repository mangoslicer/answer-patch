package datastores

import (
	"reflect"
	"testing"
	"time"

	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/settings"
)

var GlobalQuestionStore *QuestionStore

func init() {
	settings.SetPreproductionEnv()
	GlobalQuestionStore = &QuestionStore{ConnectToPostgres()}
}

func TestFindPostByID(t *testing.T) {

	expectedQuestion := &models.Question{ID: "0a24c4cd-4c73-42e4-bcca-3844d088de85", UserID: "85c3bdbc-5882-4571-aaee-e46a32713e91", Username: "Tester3", Category: "Balling", Title: "Can Jordans make me a sick baller?", Content: "I need to improve my game", Upvotes: 10, EditCount: 4, PendingCount: 3}

	expectedAnswer := &models.Answer{ID: "b50f0224-3fda-435b-a8a6-8257fcbf5aa7", QuestionID: "0a24c4cd-4c73-42e4-bcca-3844d088de85", UserID: "baeee18f-45db-4e68-81c4-25671beaab5f", Username: "Tester6", IsCurrentAnswer: true, Content: "Yeah, get the ones with the neon laces", Upvotes: 26, ReqUpvotes: 20}

	retreivedQuestion, retreivedAnswer, err, _ := GlobalQuestionStore.FindPostByID(expectedQuestion.ID)
	if err != nil {
		t.Error(err)
	}

	// Avoids reflect.DeepEqual from complaining about timestamp difference
	retreivedQuestion.SubmittedAt = time.Time{}
	retreivedAnswer.LastEditedAt = time.Time{}

	if !reflect.DeepEqual(retreivedQuestion, expectedQuestion) {
		t.Errorf("\n\n\nExpected:\n %+v\n Recieved:\n %+v\n\n", expectedQuestion, retreivedQuestion)

	}

	if !reflect.DeepEqual(retreivedAnswer, expectedAnswer) {
		t.Errorf("\n\n\nExpected:\n %+v\n Recieved:\n %+v\n\n", expectedAnswer, retreivedAnswer)
	}

}

func TestFindQuestionsByPostedBy(t *testing.T) {

	expectedQuestions := []*models.Question{&models.Question{ID: "526c4576-0e49-4e90-b760-e6976c698574", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "City Dining", Title: "Where is the best sushi place?", Content: "I have cravings", Upvotes: 15, EditCount: 8, PendingCount: 7}, &models.Question{ID: "38681976-4d2d-4581-8a68-1e4acfadcfa0", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "Gains", Title: "What should my squat to bench ratio be?", Content: "I need gains", Upvotes: 13, EditCount: 7, PendingCount: 6}}

	retreivedQuestions, err, _ := GlobalQuestionStore.FindQuestionsByFilter("posted-by", "Tester1")
	if err != nil {
		t.Error(err)
	}

	checkQuestionsForEquality(t, expectedQuestions, retreivedQuestions)
}

func TestFindQuestionsByAnsweredBy(t *testing.T) {

	expectedQuestions := []*models.Question{&models.Question{ID: "526c4576-0e49-4e90-b760-e6976c698574", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "City Dining", Title: "Where is the best sushi place?", Content: "I have cravings", Upvotes: 15, EditCount: 8, PendingCount: 7}}

	retreivedQuestions, err, _ := GlobalQuestionStore.FindQuestionsByFilter("answered-by", "Tester4")
	if err != nil {
		t.Error(err)
	}

	checkQuestionsForEquality(t, expectedQuestions, retreivedQuestions)

}

func TestSortQuestionsByUpvotes(t *testing.T) {

	//Test postComponent: "question", filter: "upvotes", order: "desc"
	expectedQuestions := []*models.Question{&models.Question{ID: "b19dc050-5ab2-417b-931c-d02445c27aca", UserID: "df38ea24-e67b-43c6-92bf-184cecee3003", Username: "Tester4", Category: "Gains", Title: "How can I convince people to skip leg day?", Content: "Please", Upvotes: 15, EditCount: 5, PendingCount: 4}, &models.Question{ID: "526c4576-0e49-4e90-b760-e6976c698574", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "City Dining", Title: "Where is the best sushi place?", Content: "I have cravings", Upvotes: 15, EditCount: 8, PendingCount: 7}, &models.Question{ID: "38681976-4d2d-4581-8a68-1e4acfadcfa0", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "Gains", Title: "What should my squat to bench ratio be?", Content: "I need gains", Upvotes: 13, EditCount: 7, PendingCount: 6}}

	retreivedQuestions, err, _ := GlobalQuestionStore.SortQuestions("question", "upvotes", "DESC", "0")
	if err != nil {
		t.Error(err)
	}

	checkQuestionsForEquality(t, expectedQuestions, retreivedQuestions)

}

func TestSortQuestionsByDate(t *testing.T) {

	//Test postComponent: "answer", filter: "date", order: "asc"
	expectedQuestions := []*models.Question{&models.Question{ID: "526c4576-0e49-4e90-b760-e6976c698574", UserID: "0c1b2b91-9164-4d52-87b0-9c4b444ee62d", Username: "Tester1", Category: "City Dining", Title: "Where is the best sushi place?", Content: "I have cravings", Upvotes: 15, EditCount: 8, PendingCount: 7}, &models.Question{ID: "0a24c4cd-4c73-42e4-bcca-3844d088de85", UserID: "85c3bdbc-5882-4571-aaee-e46a32713e91", Username: "Tester3", Category: "Balling", Title: "Can Jordans make me a sick baller?", Content: "I need to improve my game", Upvotes: 10, EditCount: 4, PendingCount: 3}, &models.Question{ID: "b19dc050-5ab2-417b-931c-d02445c27aca", UserID: "df38ea24-e67b-43c6-92bf-184cecee3003", Username: "Tester4", Category: "Gains", Title: "How can I convince people to skip leg day?", Content: "Please", Upvotes: 15, EditCount: 5, PendingCount: 4}}

	retreivedQuestions, err, _ := GlobalQuestionStore.SortQuestions("answer", "date", "ASC", "0")
	if err != nil {
		t.Error(err)
	}

	checkQuestionsForEquality(t, expectedQuestions, retreivedQuestions)
}

func TestStoreQuestion(t *testing.T) {

	err, _ := GlobalQuestionStore.StoreQuestion("{95954f28-a8c3-4e76-8c80-18de07931639}", "{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}", "Title", "Content and stuff")
	if err != nil {
		t.Error(err)
	}

	row, err := GlobalQuestionStore.DB.Query(`SELECT title FROM question WHERE title = 'Title'`)
	if err != nil {
		t.Error(err)
	} else if !row.Next() {
		t.Errorf("Failed to insert question into the ap1 database through QuestionStore's StoreQuestion method")
	}
}

func TestStoreQuestionWithForeignKeyViolation(t *testing.T) {

	//Nonexistent uuid provided for userID param
	err, _ := GlobalQuestionStore.StoreQuestion("{89f0b6aa-0399-4b31-8f24-cc4989f60391}", "{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}", "Different title", "Content and stuff")

	expectedErrMessage := "The provided user_id does not exist"

	if err.Error() != expectedErrMessage {
		t.Errorf("Expected an error message of %s, but recieved an error message of %s", expectedErrMessage, err.Error())
	}
}

func TestStoreQuestionWithUniqueConstraintViolation(t *testing.T) {

	// Title is not unique
	err, _ := GlobalQuestionStore.StoreQuestion("{89f0b6aa-0399-4b31-8f24-cc4989f60391}", "{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}", "Title", "Content and stuff")

	expectedErrMessage := "The provided title is not unique"

	if err.Error() != expectedErrMessage {
		t.Errorf("Expected an error message of %s, but recieved an error message of %s", expectedErrMessage, err.Error())
	}
}

func checkQuestionsForEquality(t *testing.T, expected []*models.Question, received []*models.Question) {

	if received == nil {
		t.Errorf("Recieved nil, but expected %#v", expected)
	}

	for i, _ := range expected {
		received[i].SubmittedAt = time.Time{}
		if !reflect.DeepEqual(received[i], expected[i]) {
			t.Errorf("\n\nExpected %#v,\n but received %#v\n\n", expected[i], received[i])
		}
	}
}
