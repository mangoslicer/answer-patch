package datastores

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/patelndipen/AP1/models"
)

type QuestionStoreServices interface {
	FindPostByID(string) (*models.Question, *models.Answer, error, int)
	FindQuestionsByFilter(string, string) ([]*models.Question, error, int)
	SortQuestions(string, string, string, string) ([]*models.Question, error, int)
	StoreQuestion(string, string, string, string) (error, int)
}

type QuestionStore struct {
	DB *sql.DB
}

func (store *QuestionStore) FindPostByID(questionID string) (*models.Question, *models.Answer, error, int) {

	row, err := store.DB.Query(`SELECT q.id, q.user_id, u.username, c.category_name, q.title, q.content, q.upvotes, q.edit_count, q.pending_count, q.submitted_at FROM question q INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id WHERE q.id =$1`, questionID)
	if err != nil {
		log.Fatal(err)
		return nil, nil, InternalErr, http.StatusInternalServerError
	} else if !row.Next() { // row.Next returns false, if 0 rows were returned by the query
		return nil, nil, errors.New("No question exists with the id of " + questionID), http.StatusBadRequest
	}

	question := new(models.Question)
	err = row.Scan(&question.ID, &question.UserID, &question.Username, &question.Category, &question.Title, &question.Content, &question.Upvotes, &question.EditCount, &question.PendingCount, &question.SubmittedAt)
	if err != nil {
		log.Fatal(err)
		return nil, nil, InternalErr, http.StatusInternalServerError
	}

	row, err = store.DB.Query(`SELECT a.id, a.question_id, a.user_id, u.username, a.is_current_answer, a.content, a.upvotes, a.required_upvotes, a.last_edited_at FROM answer a INNER JOIN ap_user u ON a.user_id = u.id WHERE a.question_id = $1 AND is_current_answer = 'true'`, questionID)
	if err != nil {
		log.Fatal(err)
		return nil, nil, InternalErr, http.StatusInternalServerError
	} else if !row.Next() {
		return question, nil, nil, http.StatusOK // Returns only a question, if the question lacks any valid answer at the current moment
	}

	answer := models.NewAnswer()
	err = row.Scan(&answer.ID, &answer.QuestionID, &answer.UserID, &answer.Username, &answer.IsCurrentAnswer, &answer.Content, &answer.Upvotes, &answer.ReqUpvotes, &answer.LastEditedAt)
	if err != nil {
		log.Fatal(err)
		return nil, nil, InternalErr, http.StatusInternalServerError
	}

	return question, answer, nil, http.StatusOK
}

func (store *QuestionStore) FindQuestionsByFilter(filter, val string) ([]*models.Question, error, int) {

	queryStmt := `SELECT q.id, q.user_id, u.username, c.category_name, q.title, q.content, q.upvotes, q.edit_count, q.pending_count, q.submitted_at FROM question q`

	switch {
	case filter == "posted-by":
		queryStmt += ` INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id WHERE u.username = $1 ORDER BY q.upvotes DESC`
	case filter == "answered-by":
		queryStmt += ` JOIN ap_user answer_author ON answer_author.username = $1 JOIN answer a ON (answer_author.id = a.user_id AND a.question_id=q.id AND a.is_current_answer='true') INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id ORDER BY q.upvotes DESC`
	case filter == "category":
		queryStmt += ` INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id WHERE c.category_name = $1 ORDER BY q.upvotes DESC`
	}

	rows, err := store.DB.Query(queryStmt, val)
	if err != nil {
		log.Fatal(err)
		return nil, InternalErr, http.StatusInternalServerError
	}

	return scanQuestions(rows)
}

func (store *QuestionStore) SortQuestions(postComponent, filter, order, offset string) ([]*models.Question, error, int) {

	var ok bool

	// The following maps convert the  param "filter" into a valid database column name
	questionFilters := map[string]string{
		"upvotes": "q.upvotes",
		"date":    "q.submitted_at",
		"edits":   "q.edit_count",
	}
	answerFilters := map[string]string{
		"upvotes": "a.upvotes",
		"date":    "a.last_edited_at",
	}
	queryStmt := `SELECT q.id, q.user_id, u.username, c.category_name, q.title, q.content, q.upvotes, q.edit_count, q.pending_count, q.submitted_at FROM question q`
	if postComponent == "question" {
		queryStmt += ` INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id`
		filter, ok = questionFilters[filter]
	} else if postComponent == "answer" {
		queryStmt += ` JOIN answer a ON (a.question_id=q.id AND a.is_current_answer='true') INNER JOIN ap_user u ON q.user_id = u.id INNER JOIN category c ON q.category_id = c.id`
		filter, ok = answerFilters[filter]
	}
	if !ok { // Return nil if the url param "filter" can not be converted into a valid database column name
		return nil, errors.New("Could not recognize the sorting criteria"), http.StatusBadRequest
	}

	queryStmt += ` ORDER BY ` + filter + ` ` + strings.ToUpper(order) + ` LIMIT 10 OFFSET $1`

	rows, err := store.DB.Query(queryStmt, offset)
	if err != nil {
		log.Fatal(err)
		return nil, InternalErr, http.StatusInternalServerError
	}

	return scanQuestions(rows)
}

func (store *QuestionStore) StoreQuestion(userID, categoryID, title, content string) (error, int) {

	return transact(store.DB, func(tx *sql.Tx) (error, int) {
		_, err := tx.Exec(`INSERT INTO question(user_id, category_id, title, content) values($1::uuid, $2::uuid, $3, $4)`, userID, categoryID, title, content)
		if err != nil {
			return evaluateSQLError(err)
		}

		return nil, http.StatusOK
	})

}

func scanQuestions(rows *sql.Rows) ([]*models.Question, error, int) {
	var questions []*models.Question

	for rows.Next() {
		tempQuestion := new(models.Question)
		err := rows.Scan(&tempQuestion.ID, &tempQuestion.UserID, &tempQuestion.Username, &tempQuestion.Category, &tempQuestion.Title, &tempQuestion.Content, &tempQuestion.Upvotes, &tempQuestion.EditCount, &tempQuestion.PendingCount, &tempQuestion.SubmittedAt)
		if err != nil {
			log.Fatal(err)
			return nil, InternalErr, http.StatusInternalServerError
		}
		questions = append(questions, tempQuestion)
	}

	if len(questions) == 0 {
		return nil, errors.New("No question(s) founds"), http.StatusBadRequest
	}
	return questions, nil, http.StatusOK
}
