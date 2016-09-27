package datastores

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"regexp"
	"time"

	_ "github.com/lib/pq"
	"github.com/patelndipen/AP1/settings"
)

var (
	InternalErr = errors.New("Internal error")
)

func ConnectToPostgres() *sql.DB {

	dns := settings.GetPostgresDSN()

	db, err := sql.Open("postgres", dns)
	if err != nil {
		log.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
	return db
}

func transact(db *sql.DB, fn func(*sql.Tx) (error, int)) (error, int) {

	tx, err := db.Begin()
	if err != nil {
		return InternalErr, http.StatusInternalServerError
	}

	err, statusCode := fn(tx)

	if err != nil {
		tx.Rollback()
		return err, statusCode
	} else {
		tx.Commit()
	}

	return nil, statusCode
}

func standardizeTime(x *time.Time, y *time.Time) {
	*x = *y
}

func initializePostgres(db *sql.DB) {

	// Run 'create extension "uuid-ossp";' in your psql shell in order for uuid_generate_v4() to be a valid value for columns

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS ap_user(id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), username varchar(20) NOT NULL UNIQUE, hashed_password char(60) NOT NULL UNIQUE, created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'))`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS category (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), category_name varchar(15) NOT NULL, user_id uuid REFERENCES ap_user NOT NULL, created_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'))`)
	if err != nil {
		log.Fatal(err)

	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS question (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), user_id uuid REFERENCES ap_user NOT NULL, category_id uuid REFERENCES category NOT NULL, title varchar(255) NOT NULL UNIQUE, content text NOT NULL, upvotes integer DEFAULT 0, edit_count integer DEFAULT 0, pending_count integer DEFAULT 0, submitted_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'))`)
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS answer (id uuid PRIMARY KEY DEFAULT uuid_generate_v4(), question_id uuid REFERENCES question ON DELETE CASCADE NOT NULL, user_id uuid REFERENCES ap_user NOT NULL, content text, upvotes integer DEFAULT 0, required_upvotes integer DEFAULT 0, is_current_answer boolean DEFAULT false, last_edited_at TIMESTAMP WITHOUT TIME ZONE DEFAULT (CURRENT_TIMESTAMP AT TIME ZONE 'UTC'))`)
	if err != nil {
		log.Fatal(err)
	}
}

func dropPostgresTables(db *sql.DB) {

	var err error
	tables := [4]string{"answer", "question", "category", "ap_user"}

	for _, t := range tables {

		_, err = db.Exec(`DROP TABLE IF EXISTS ` + t)
		if err != nil {
			log.Fatal(err)
		}
	}

}

func isCategoryRegistered(DB *sql.DB, category string) (bool, error) {

	row, err := DB.Query(`SELECT category_name WHERE category_name=$1`, category)
	if err != nil {
		return false, InternalErr
	}

	return row.Next(), nil
}

func evaluateSQLError(err error) (error, int) {

	var r *regexp.Regexp

	matched, _ := regexp.MatchString("violates foreign key constraint", err.Error())

	if matched == true {
		r = regexp.MustCompile("user_id|category_id|question_id")
		return errors.New("The provided " + r.FindString(err.Error()) + " does not exist"), http.StatusBadRequest
	}

	matched, _ = regexp.MatchString("duplicate key value violates unique constraint", err.Error())

	if matched == true {
		r = regexp.MustCompile("username|title")
		return errors.New("The provided " + r.FindString(err.Error()) + " is not unique"), http.StatusConflict
	}

	return InternalErr, http.StatusInternalServerError
}

//Populates DB with questions, answers, and users for unit testing
func populatePostgres(db *sql.DB) {

	var err error

	//Users

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES('{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}'::uuid, 'Tester1', '$2a$10$lWqqb7MhwH7YryO4DyjdeOsFQ9hK7qxZ8PPcm6qjuNlM47KNInHMK')`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES ('{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid, 'Tester2', '$2a$10$16XDQDyDfQxvil6dqC7fV.tWlf/lc1kD9sA9/8qONoGEm9GiJz6vS')`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES('{85c3bdbc-5882-4571-aaee-e46a32713e91}'::uuid, 'Tester3', '$2a$10$DtKbFWQgT0Htcu9dQUwSZ.Mu2Wp0YQDiyufWVfxDp30Bi/Fpru3A2')`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES('{df38ea24-e67b-43c6-92bf-184cecee3003}'::uuid, 'Tester4', '$2a$10$5uDl2b6dNMYsJ/G3AnnM3.4UkHmiAAW7t2u4CphLQnsmDpUec9wfe')`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES('{61633349-89f3-43c9-ac91-653b3229ecf7}'::uuid, 'Tester5', '$2a$10$LT3lw8NX7Ybnt0QJn511zuuP9JDBSW8g3/YSOasqN6L1b41SYvRDa')`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO ap_user(id, username, hashed_password) VALUES('{baeee18f-45db-4e68-81c4-25671beaab5f}'::uuid, 'Tester6', '$2a$10$mks.PzG/45yOGtmmj9N..e.uUg8tdp9psNqK9Xw5kCRMy2ZFWuR6e')`); err != nil {
		log.Fatal(err)
	}

	//Categories

	if _, err = db.Exec(`INSERT INTO category(id, category_name, user_id) VALUES ('{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}'::uuid,'Gains','{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid)`); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`INSERT INTO category(id, category_name, user_id) VALUES ('{7d2b570d-54c6-48b1-8f46-68304f163d6a}'::uuid, 'City Dining', '{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid)`); err != nil {
		log.Fatal(err)
	}
	if _, err = db.Exec(`INSERT INTO category(id, category_name, user_id) VALUES ('{cb996d64-bd2d-414c-bbdc-81faba62cdc2}'::uuid, 'Balling', '{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid)`); err != nil {
		log.Fatal(err)
	}

	//Questions

	if _, err = db.Exec(`INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{38681976-4d2d-4581-8a68-1e4acfadcfa0}'::uuid,'{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}'::uuid, '{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}'::uuid, 'What should my squat to bench ratio be?', 'I need gains', 13, 4, 1)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(` INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{526c4576-0e49-4e90-b760-e6976c698574}'::uuid,'{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}'::uuid, '{7d2b570d-54c6-48b1-8f46-68304f163d6a}'::uuid, 'Where is the best sushi place?', 'I have cravings', 15, 5, 2)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(` INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{0a24c4cd-4c73-42e4-bcca-3844d088de85}'::uuid,'{85c3bdbc-5882-4571-aaee-e46a32713e91}'::uuid, '{cb996d64-bd2d-414c-bbdc-81faba62cdc2}'::uuid, 'Can Jordans make me a sick baller?', 'I need to improve my game', 10, 1, 5)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(` INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{b19dc050-5ab2-417b-931c-d02445c27aca}'::uuid,'{df38ea24-e67b-43c6-92bf-184cecee3003}'::uuid, '{33f6b77a-4564-4aa9-8cc8-50bb01c6a609}'::uuid, 'How can I convince people to skip leg day?', 'Please', 15, 2, 2)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(` INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{28a12532-bc7a-427c-8f55-b72b18df7c02}'::uuid,'{61633349-89f3-43c9-ac91-653b3229ecf7}'::uuid, '{7d2b570d-54c6-48b1-8f46-68304f163d6a}'::uuid, 'Should I sign up for a Groupon account?', 'I like to dine at new places', 5, 3, 1)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(` INSERT INTO question(id, user_id, category_id, title, content, upvotes, edit_count, pending_count) VALUES('{bf8111f3-e75f-40d7-8d5a-813ce3a429fe}'::uuid,'{baeee18f-45db-4e68-81c4-25671beaab5f}'::uuid, '{cb996d64-bd2d-414c-bbdc-81faba62cdc2}'::uuid, 'Is ball really life?', 'I am having an existential crisis', 10, 2, 1)`); err != nil {
		log.Fatal(err)
	}

	//Answers

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{f46fd5c9-ea9b-4677-ba8a-433b27fc097c}'::uuid, '{38681976-4d2d-4581-8a68-1e4acfadcfa0}'::uuid, '{61633349-89f3-43c9-ac91-653b3229ecf7}'::uuid, 'false', 'Always to never', 20, 25)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{150aebd1-a381-4ba5-a612-cee110f771f0}'::uuid, '{38681976-4d2d-4581-8a68-1e4acfadcfa0}'::uuid, '{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid, 'false', 'It depends on the amount of cardio you do before leg day', 10, 25)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{4bf6d7e3-681b-4ec3-9353-34490aba965b}'::uuid, '{38681976-4d2d-4581-8a68-1e4acfadcfa0}'::uuid, '{85c3bdbc-5882-4571-aaee-e46a32713e91}'::uuid, 'false', 'Why would you disgrace the bench?', 14, 15)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{c6f753ea-8b55-468f-9eb2-3ac03f6ed179}'::uuid, '{526c4576-0e49-4e90-b760-e6976c698574}'::uuid,'{df38ea24-e67b-43c6-92bf-184cecee3003}'::uuid, 'true', 'Not Utah', 40, 15)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{7253b7cd-0783-4b29-a11c-90bbc5d09c0e}'::uuid, '{526c4576-0e49-4e90-b760-e6976c698574}'::uuid,'{61633349-89f3-43c9-ac91-653b3229ecf7}'::uuid, 'false', 'Not Massachusetts', 10, 15)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{b50f0224-3fda-435b-a8a6-8257fcbf5aa7}'::uuid, '{0a24c4cd-4c73-42e4-bcca-3844d088de85}'::uuid,'{baeee18f-45db-4e68-81c4-25671beaab5f}'::uuid, 'true', 'Yeah, get the ones with the neon laces', 25, 20)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{fbd3d2ac-df1f-4861-8e46-9dd902f6f071}'::uuid, '{0a24c4cd-4c73-42e4-bcca-3844d088de85}'::uuid,'{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}'::uuid, 'false', 'Yeah, get the shinny ones', 25, 20)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{d8e17edf-d58c-49d9-81a7-72badb2786e3}'::uuid, '{b19dc050-5ab2-417b-931c-d02445c27aca}'::uuid,'{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid, 'true', 'Convince them that small calfs are genetic', 40, 4)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{c0392c0e-bd4e-41f7-94e7-ddd03ae58416}'::uuid, '{b19dc050-5ab2-417b-931c-d02445c27aca}'::uuid,'{85c3bdbc-5882-4571-aaee-e46a32713e91}'::uuid, 'false', 'Be honest about the gain loss', 45, 10)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{3b745a45-d085-476d-b909-11a1164dddb2}'::uuid, '{28a12532-bc7a-427c-8f55-b72b18df7c02}'::uuid,'{df38ea24-e67b-43c6-92bf-184cecee3003}'::uuid, 'false', 'Yes, Groupon can save you a ton of money', 25, 25)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{924310dc-9f18-447d-9dab-301653aed3bf}'::uuid, '{28a12532-bc7a-427c-8f55-b72b18df7c02}'::uuid,'{baeee18f-45db-4e68-81c4-25671beaab5f}'::uuid, 'false', 'Only if can stand disappointing food', 1, 25)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{43573548-9d27-4bc9-a0b1-5d72c8f6d5a5}'::uuid, '{bf8111f3-e75f-40d7-8d5a-813ce3a429fe}'::uuid,'{0c1b2b91-9164-4d52-87b0-9c4b444ee62d}'::uuid, 'false', 'Well, this is has puzzled many philosophers. In fact, Albert Camus thought that this question was the only true philosophical question. There is not simple answer.', 30, 20)`); err != nil {
		log.Fatal(err)
	}

	if _, err = db.Exec(`INSERT INTO answer(id, question_id, user_id, is_current_answer, content, upvotes, required_upvotes) VALUES ('{7e0dca3b-0477-42c0-a501-05a6f89288c8}'::uuid, '{bf8111f3-e75f-40d7-8d5a-813ce3a429fe}'::uuid,'{95954f28-a8c3-4e76-8c80-18de07931639}'::uuid, 'false', 'Yes.', 30, 20)`); err != nil {
		log.Fatal(err)
	}

}
