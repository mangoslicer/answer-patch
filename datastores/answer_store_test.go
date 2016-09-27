package datastores

import (
	"log"
	"testing"

	"github.com/patelndipen/AP1/models"
	"github.com/patelndipen/AP1/settings"
)

var GlobalAnswerStore *AnswerStore

func init() {

	settings.SetPreproductionEnv()
	GlobalAnswerStore = &AnswerStore{ConnectToPostgres()}

	dropPostgresTables(GlobalAnswerStore.DB)
	initializePostgres(GlobalAnswerStore.DB)
	populatePostgres(GlobalAnswerStore.DB)

}

func TestIsAnswerSlotAvailable(t *testing.T) {

	slotTests := []struct {
		questionID string
		expected   bool
	}{
		{"38681976-4d2d-4581-8a68-1e4acfadcfa0", true},
		{"526c4576-0e49-4e90-b760-e6976c698574", true},
		{"0a24c4cd-4c73-42e4-bcca-3844d088de85", false},
	}

	for _, st := range slotTests {
		result, err := GlobalAnswerStore.IsAnswerSlotAvailable(st.questionID)
		if err != nil {
			t.Error(err)
		} else if result != st.expected {
			t.Errorf("Expected a result of %t for the question with an id of %s, but recieved %t", st.expected, st.questionID, result)
		}
	}
}

func TestStoreAnswer(t *testing.T) {

	newAnswer := &models.Answer{QuestionID: "{0a24c4cd-4c73-42e4-bcca-3844d088de85}", UserID: "{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}", Content: "very new", ReqUpvotes: 10}

	GlobalAnswerStore.StoreAnswer(newAnswer.QuestionID, newAnswer.UserID, newAnswer.Content, newAnswer.ReqUpvotes)

	row, err := GlobalAnswerStore.DB.Query(`SELECT content FROM answer WHERE user_id = $1::uuid AND content = $2`, newAnswer.UserID, newAnswer.Content)
	if err != nil {
		t.Error(err)
	} else if !row.Next() {
		t.Errorf("Failed to insert answer into the ap1 database through AnswerStore's StoreAnswer method")
	}
}

func TestStoreAnswerWithExistingAnswer(t *testing.T) {

	var pendingCount int

	existingAnswer := &models.Answer{QuestionID: "b19dc050-5ab2-417b-931c-d02445c27aca", UserID: "95954f28-a8c3-4e76-8c80-18de07931639", Content: "Convince them that small calfs are genetic", ReqUpvotes: 4}

	GlobalAnswerStore.StoreAnswer(existingAnswer.QuestionID, existingAnswer.UserID, existingAnswer.Content, existingAnswer.ReqUpvotes)

	row := GlobalAnswerStore.DB.QueryRow(`SELECT pending_count FROM question WHERE id = $1::uuid`, existingAnswer.QuestionID)
	err := row.Scan(&pendingCount)
	if err != nil {
		t.Error(err)
	}

	if pendingCount != 2 { //Two was the pending count before attempting to store the existing answer
		t.Errorf("Expected the pending count for the question with an ID of %s to remain two, because the StoreAnswer method does not write an answer, if the answer aready exists. The pending count retrieved is %d", existingAnswer.QuestionID, pendingCount)
	}
}

func TestCastVote(t *testing.T) {

	var expectedUserID, retrievedUserID string
	var originalUpvotes, retrievedUpvotes int

	answerID := "b50f0224-3fda-435b-a8a6-8257fcbf5aa7"

	row := GlobalAnswerStore.DB.QueryRow(`SELECT user_id, upvotes FROM answer WHERE id = $1::uuid`, answerID)
	err := row.Scan(&expectedUserID, &originalUpvotes)
	if err != nil {
		t.Error(err)
	}

	retrievedUserID, err, _ = GlobalAnswerStore.CastVote(answerID, 1)
	if err != nil {
		t.Error(err)
	}

	row = GlobalAnswerStore.DB.QueryRow(`SELECT upvotes FROM answer WHERE id = $1::uuid`, answerID)
	err = row.Scan(&retrievedUpvotes)
	if err != nil {
		t.Error(err)
	}

	if (originalUpvotes + 1) != retrievedUpvotes {
		t.Errorf("Expected the answer with an ID of %s to have %d upvotes, but the answer has %d upvotes", answerID, (originalUpvotes + 1), retrievedUpvotes)
	} else if expectedUserID != retrievedUserID {
		t.Errorf("Expected the answer with an ID of %s to have a user ID of %s, but the query returned a user ID of %s", answerID, expectedUserID, retrievedUserID)
	}
}

func TestCastVoteWithNonexistantAnswerID(t *testing.T) {

	_, err, _ := GlobalAnswerStore.CastVote("1da8f5f3-271e-4f35-a0dc-d2935effc524", -1) //the provided UUID does not exist

	if err.Error() != "No answer exists with the provided answer id" {
		t.Errorf("Expected the CastVote to return \"No answer exists with the provided answer id\", but CastVote returned %s", err.Error())
	}
}

func TestAssessAnswersWithNoQualifiedCurrentAnswers(t *testing.T) {

	questionID := "38681976-4d2d-4581-8a68-1e4acfadcfa0"

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)

	row, err := GlobalAnswerStore.DB.Query(`SELECT id FROM answer WHERE question_id = $1 AND is_current_answer = 'true'`, questionID)
	if err != nil {
		t.Error(err)
	} else if row.Next() {
		var currentAnswerID string
		err = row.Scan(&currentAnswerID)
		if err != nil {
			t.Error(err)
		}
		t.Errorf("Expected there to be no qualified current answer, but %s was detected to be the current answer", currentAnswerID)
	}

}

func TestAssessAnswersWithNewlyQualifiedCurrentAnswer(t *testing.T) {
	var retrievedUserID string
	questionID := "28a12532-bc7a-427c-8f55-b72b18df7c02"
	expectedUserID := "df38ea24-e67b-43c6-92bf-184cecee3003"

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)
	if err != nil {
		t.Error(err)
	}

	row, err := GlobalAnswerStore.DB.Query(`SELECT user_id FROM answer WHERE question_id = $1 AND is_current_answer = 'true'`, questionID)
	if err != nil {
		t.Error(err)
	}

	row.Next()
	err = row.Scan(&retrievedUserID)
	if err != nil {
		t.Error(err)
	}

	if expectedUserID != retrievedUserID {
		t.Errorf("Expected Tester4's (User ID of %s) answer to be the current answer, but the user ID of %s was detected to have the best possible current answer", questionID, expectedUserID, retrievedUserID)
	}
}

func findCurrentAnswerUserID(questionID string) string {
	var currentAnswerUserID string

	row, err := GlobalAnswerStore.DB.Query(`SELECT user_id FROM answer WHERE question_id = $1 AND is_current_answer = 'true'`, questionID)
	if err != nil {
		log.Fatal(err)
	}

	row.Next()
	err = row.Scan(&currentAnswerUserID)
	if err != nil {
		log.Fatal(err)
	}

	return currentAnswerUserID
}

func TestAssessAnswersWithNoNewCurrentAnswer(t *testing.T) {
	questionID := "526c4576-0e49-4e90-b760-e6976c698574"

	expectedCurrentAnswerUserID := findCurrentAnswerUserID(questionID)

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)
	if err != nil {
		t.Error(err)
	}

	retrievedCurrentAnswerUserID := findCurrentAnswerUserID(questionID)

	if expectedCurrentAnswerUserID != retrievedCurrentAnswerUserID {
		t.Errorf("Expected the user with an ID of %s to still have the best current answer for the question with an ID of %s, but the current answer was detected to be %s ", expectedCurrentAnswerUserID, questionID, retrievedCurrentAnswerUserID)
	}
}

func TestAssessAnswersWithNewCurrentAnswer(t *testing.T) {

	questionID := "b19dc050-5ab2-417b-931c-d02445c27aca"
	expectedCurrentAnswerUserID := "85c3bdbc-5882-4571-aaee-e46a32713e91"

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)

	if err != nil {
		t.Error(err)
	}

	retrievedCurrentAnswerUserID := findCurrentAnswerUserID(questionID)

	if expectedCurrentAnswerUserID != retrievedCurrentAnswerUserID {
		t.Errorf("Expected the question with an ID of  %s to have the best current answer set to Tester3's (User ID of %s)  answer, but the user with an ID of %s was detected to have the best current answer", questionID, expectedCurrentAnswerUserID, retrievedCurrentAnswerUserID)
	}
}

func TestAssessAnswersWithSameUpvotesAndCurrentAnswer(t *testing.T) {

	questionID := "0a24c4cd-4c73-42e4-bcca-3844d088de85"
	expectedCurrentAnswerUserID := "baeee18f-45db-4e68-81c4-25671beaab5f"

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)
	if err != nil {
		t.Error(err)
	}

	retrievedCurrentAnswerUserID := findCurrentAnswerUserID(questionID)

	if retrievedCurrentAnswerUserID != expectedCurrentAnswerUserID {
		t.Errorf("Expected the question withan ID of %s to have a current answer that was posted by Tester6, who has a User ID of %s. Tester6's answer is still supposed to be the current answer despite the fact that a different answer has the same amount of upvotes, but the user ID of %s as the user with the best possible current answer", questionID, expectedCurrentAnswerUserID, retrievedCurrentAnswerUserID)
	}
}

func TestAssessAnswersWithSameUpvotesAndNoCurrentAnswer(t *testing.T) {

	questionID := "bf8111f3-e75f-40d7-8d5a-813ce3a429fe"
	expectedCurrentAnswerUserID := "0c1b2b91-9164-4d52-87b0-9c4b444ee62d"

	err, _ := GlobalAnswerStore.AssessAnswers(questionID)
	if err != nil {
		t.Error(err)
	}

	retrievedCurrentAnswerUserID := findCurrentAnswerUserID(questionID)

	if retrievedCurrentAnswerUserID != expectedCurrentAnswerUserID {
		t.Errorf("Expected the question withan ID of %s to have a current answer that was posted by Tester1, who has a User ID of %s, but the user with the ID of %s was detected to have the current answer", questionID, expectedCurrentAnswerUserID, retrievedCurrentAnswerUserID)
	}
}
