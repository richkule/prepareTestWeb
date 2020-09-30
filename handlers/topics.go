package handlers

import (
	"errors"
	db "github.com/richkule/prepareTestWeb/DBWorker"
	in "github.com/richkule/prepareTestWeb/init"
	"net/http"
	"strconv"
)

// Обрабатывает создание новой темы
func NewTopic(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if req.Method == "GET" || sessUs.UsId == in.GuestUserId {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}
	tName := in.TopicName(req.FormValue("TopicName"))
	tDesc := in.TopicDesc(req.FormValue("TopicDesc"))
	strTestId := req.FormValue("TestId")
	intTestId, err := strconv.Atoi(strTestId)
	if err != nil {
		err = errors.New("Ошибка конвертации id теста newTopic " + err.Error())
		return err
	}
	testId := in.TestId(intTestId)
	err = db.CreateTopic(tName, tDesc, testId)
	if err != nil {
		err = errors.New("Ошибка создания темы newTopic " + err.Error())
		return err
	}
	http.Redirect(w, req, `/edit/test/`+strTestId, http.StatusFound)
	return nil
}

// Обрабатывает страницу редактирования темы
func editTopic(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs, topicId in.TopicId) error {
	ok, err := db.CheckAuthorTopic(topicId, sessUs.UsId)
	if err != nil {
		err = errors.New("Ошибка получения записи автора из БД editTest " + err.Error())
		return err
	}
	if !ok {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}
	data := in.DataEditTopic{}
	if data.Header, err = renderHeader(sessUs.UsId); err != nil {
		err = errors.New("Ошибка обработки шаблона шапки editTest " + err.Error())
		return err
	}
	data.Questions, err = db.GetQuestions(topicId)
	if err != nil {
		err = errors.New("Ошибка получения тем editTest " + err.Error())
		return err
	}
	err = renderTemplate(w, in.EditTopicPage, data)
	if err != nil {
		err = errors.New("Ошибка обработки шаблона editTestPage " + err.Error())
	}
	return err
}
