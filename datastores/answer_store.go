package datastores

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/patelndipen/AP1/models"
)

type AnswerStoreServices interface {
	IsAnswerSlotAvailable(string) (bool, error)
	StoreAnswer(string, string, string, int) (error, int)
	CastVote(string, int) (string, error, int)
	AssessAnswers(string) (error, int)
}

type AnswerStore struct {
	DB *sql.DB
}

func (store *AnswerStore) IsAnswerSlotAvailable(questionID string) (bool, error) {

	var pending int

	row := store.DB.QueryRow(`SELECT pending_count FROM question WHERE id = $1`, questionID)

	err := row.Scan(&pending)
	if err != nil {
		log.Fatal(err)
		return false, InternalErr
	}

	return pending < 5, nil
}

func (store *AnswerStore) StoreAnswer(questionID, userID, content string, reqUpvotes int) (error, int) {

	row, err := store.DB.Query(`SELECT id FROM answer WHERE question_id = $1::uuid AND user_id = $2::uuid AND content = $3 AND required_upvotes = $4`, questionID, userID, content, reqUpvotes)
	if err != nil {
		evaluateSQLError(err)
	} else if row.Next() {
		return errors.New("Question already exists"), http.StatusBadRequest
	}

	return transact(store.DB, func(tx *sql.Tx) (error, int) {
		_, err := tx.Exec(`INSERT INTO answer(question_id, user_id, content, required_upvotes) values($1::uuid, $2::uuid, $3, $4)`, questionID, userID, content, reqUpvotes)
		if err != nil {
			log.Fatal(err)
			return evaluateSQLError(err)
		}

		_, err = tx.Exec(`UPDATE question SET pending_count = pending_count + 1 WHERE id = $1::uuid`, questionID)
		if err != nil {
			log.Fatal(err)
			return evaluateSQLError(err)
		}

		return nil, http.StatusOK
	})

}

func (store *AnswerStore) CastVote(answerID string, vote int) (string, error, int) {

	var userID string

	row, err := store.DB.Query(`SELECT user_id FROM answer WHERE id = $1::uuid`, answerID)
	if err != nil {
		log.Fatal(err)
		return "", InternalErr, http.StatusInternalServerError
	} else if !row.Next() {
		return "", errors.New("No answer exists with the provided answer id"), http.StatusBadRequest
	}

	err = row.Scan(&userID)
	if err != nil {
		log.Fatal(err)
		return "", InternalErr, http.StatusInternalServerError
	}

	err, statusCode := transact(store.DB, func(tx *sql.Tx) (error, int) {
		_, err := tx.Exec(`UPDATE answer SET upvotes = upvotes + $1 WHERE id = $2`, vote, answerID)
		if err != nil {
			log.Fatal(err)
			return evaluateSQLError(err)
		}

		return nil, http.StatusOK
	})
	if err != nil {
		return "", err, statusCode
	}

	return userID, nil, http.StatusOK
}

// AssessAnswers determines the answer that is most qualified to be considered the current answer
func (store *AnswerStore) AssessAnswers(questionID string) (error, int) {

	var qualifiedAnswers []*models.Answer
	var isCurrentAnswerExistant bool = false

	_, err := store.DB.Exec(`DELETE FROM answer WHERE upvotes = 0`)
	if err != nil {
		log.Fatal(err)
		return InternalErr, http.StatusInternalServerError
	}

	rows, err := store.DB.Query(`SELECT id, user_id, is_current_answer, upvotes, required_upvotes FROM answer WHERE question_id = $1 ORDER BY upvotes DESC, is_current_answer DESC, last_edited_at ASC`, questionID)
	if err != nil {
		log.Fatal(err)
		return InternalErr, http.StatusInternalServerError
	}

	for rows.Next() {
		tempAnswer := new(models.Answer)
		err := rows.Scan(&tempAnswer.ID, &tempAnswer.UserID, &tempAnswer.IsCurrentAnswer, &tempAnswer.Upvotes, &tempAnswer.ReqUpvotes)
		if err != nil {
			log.Fatal(err)
			return InternalErr, http.StatusInternalServerError
		}
		//Appends all answers that have satisfied their calculated required upvotes
		if tempAnswer.Upvotes >= tempAnswer.ReqUpvotes {
			qualifiedAnswers = append(qualifiedAnswers, tempAnswer)
		}

		// Breaks loop after scanning the current answer in order to avoid scanning answers that have less upvotes that the current answer
		if tempAnswer.IsCurrentAnswer == true {
			isCurrentAnswerExistant = true
			break
		}
	}

	//Either none of the answers satisfy their required amount of upvotes or the current answer is still the best candidate
	if len(qualifiedAnswers) == 0 || (len(qualifiedAnswers) == 1 && isCurrentAnswerExistant == true) {
		return nil, http.StatusOK

	}

	return transact(store.DB, func(tx *sql.Tx) (error, int) {

		_, err := tx.Exec(`UPDATE answer SET is_current_answer = 'true' WHERE id = $1`, qualifiedAnswers[0].ID)
		if err != nil {
			return evaluateSQLError(err)
		}

		if isCurrentAnswerExistant == true {
			_, err = tx.Exec(`UPDATE answer SET is_current_answer = 'false' WHERE id = $1`, qualifiedAnswers[len(qualifiedAnswers)-1].ID)
			if err != nil {
				return evaluateSQLError(err)
			}
			_, err = tx.Exec(`UPDATE question SET pending_count = edit_count + 1`)
			if err != nil {
				return evaluateSQLError(err)
			}
		}

		_, err = tx.Exec(`UPDATE question SET edit_count = edit_count + 1`)
		if err != nil {
			return evaluateSQLError(err)
		}
		return nil, http.StatusOK
	})

}
