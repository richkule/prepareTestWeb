package DBWorker

import (
	"context"
	in "github.com/richkule/prepareTestWeb/init"
)

// Получает все темы теста
func GetQuestions(id in.TopicId) ([]in.Question, error) {
	sql := `select id, html, q_type
	from questions
	where topic_id = $1`
	qArr := make([]in.Question, 0)
	res, err := db.Query(context.Background(), sql, id)
	if err != nil {
		return nil, err
	}
	for res.Next() {
		q := in.Question{}
		err := res.Scan(&q.QuestionId, q.Body, q.QuestionType)
		if err != nil {
			return nil, err
		}
		qArr = append(qArr, q)
	}
	return qArr, nil
}

// Возвращает true, если пользователь является автором темы
// В разработке
func CheckAuthorQuestion(testId in.TestId, id *in.SessUs) (bool, error) {
	sql := `select * from test_author where test_id = $1 and us_id = $2`
	res, err := db.Exec(context.Background(), sql, testId, id.UsId)
	if err != nil {
		return false, err
	}
	if res[len(res)-1] == '0' {
		return false, nil
	}
	return true, nil
}
