package handlers

import (
	"bytes"
	"errors"
	"fmt"
	db "github.com/richkule/prepareTestWeb/DBWorker"
	hf "github.com/richkule/prepareTestWeb/helpFun"
	in "github.com/richkule/prepareTestWeb/init"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"regexp"
	"strconv"
)

// Функция, для создания и обработки ручек ручек
func MakeHandler(fn in.HandlerIdFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer func() {
			if err != nil {
				wrongFun(w, err)
			}
		}()
		session, err := in.Store.Get(r, in.SessionName)
		if err != nil {
			log.Println("Ошибки куки MakeHandler ")
		}
		uid, ok := session.Values[in.CookRowName].(in.SessUs)
		// Сессия получена
		if ok {

			var oldUid *in.SessUs
			oldUid, err = db.GetUserId(uid.SessId)
			if err != nil {
				err = errors.New("Ошибка получения userId из базы makeHandler " + err.Error())
				return
			}

			if oldUid != nil { // Данная сессия уже существует в БД
				err = db.UpdateSessTime(uid.SessId)
				if err != nil {
					err = errors.New("Ошибка обновления время активности сессии makeHandler " + err.Error())
					return
				}
				err = fn(w, r, &uid)
			} else { // Сессии нет в БД
				var uid *in.SessUs
				uid, err = hf.CreateAndSetSess(w, r, session, in.GuestUserId)
				if err != nil {
					err = errors.New("Ошибка генерации или установки сессии makeHandler " + err.Error())
					return
				}
				err = fn(w, r, uid)
			}

		} else { // Куки не удалось правильно прочитать, действия как если сессии нету
			var uid *in.SessUs
			uid, err = hf.CreateAndSetSess(w, r, session, in.GuestUserId)
			if err != nil {
				err = errors.New("Ошибка генерации или установки сессии makeHandler " + err.Error())
				return
			}
			err = fn(w, r, uid)
		}
	}
}

// Обрабатывает шаблон шапки сайта
func renderHeader(userId in.UsId) (template.HTML, error) {
	var err error
	data := in.DataHeader{}
	if userId == in.GuestUserId {
		data.UserName = "Гость"

	} else {
		data.UserName, err = db.GetUserName(userId)
		if err != nil {
			err = errors.New("Ошибка получения имени пользователя " + err.Error())
			return "", err
		}
	}
	buf := bytes.NewBufferString("")
	err = renderTemplate(buf, in.HeaderPath, data)
	if err != nil {
		err = errors.New("Ошибка обработки шапки " + err.Error())
		return "", err
	}
	return template.HTML(buf.String()), nil
}

// Выводит основную страницу пользователя
func Index(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	var err error
	_ = req // Переменная необходима для совместительства с типом HandlerIdFunc
	data := in.DataIndex{}

	if data.Header, err = renderHeader(sessUs.UsId); err != nil {
		err = errors.New("Ошибка обработки шапки index " + err.Error())
		return err
	}
	if err = renderTemplate(w, in.IndexPage, data); err != nil {
		err = errors.New("Ошибка обработки шаблона index " + err.Error())
		return err
	}
	return nil
}

// Обрабатывает ошибку во время исполнения
func wrongFun(w http.ResponseWriter, err error) {
	log.Println(err.Error())
	http.Error(w, "", http.StatusBadGateway)
}

// Возвращает страницу пользователю, обрабатывая заданный шаблон
func renderTemplate(w io.Writer, pagePath string, data interface{}) error {
	pageName := filepath.Base(pagePath)
	err := in.Templates.ExecuteTemplate(w, pageName, data)
	if err != nil {
		return err
	}
	return nil
}

// Обработка страницы редактирования различных элементов /edit/elem/id
func Edit(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if sessUs.UsId == in.GuestUserId {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}

	var path string
	// Функция получающая id элемента из путя с помощью регулярного выражения
	// В случае ошибки вставит название элемента nameId в ошибку
	idFunc := func(regexp *regexp.Regexp, nameId string) (int, error) {
		elemGroup := regexp.FindStringSubmatch(path)

		// Регулярные выражения построенны так, что в первой группе всегда будет необходимый id
		id, err := strconv.Atoi(elemGroup[1])
		if err != nil {
			err := fmt.Errorf("Неправильный id %s edit %s ", nameId, err.Error())
			return 0, err
		}
		return id, nil
	}
	path = req.URL.Path
	switch {
	case in.RegTestEdit.MatchString(path):
		intTestId, err := idFunc(in.RegTestEdit, "теста")
		testId := in.TestId(intTestId)
		if err != nil {
			return err
		}
		return editTest(w, req, sessUs, testId)
	case in.RegTopicEdit.MatchString(path):
		intTopicId, err := idFunc(in.RegTopicEdit, "темы")
		topicId := in.TopicId(intTopicId)
		if err != nil {
			return err
		}
		return editTopic(w, req, sessUs, topicId)
	case in.RegQuesEdit.MatchString(path):
		quesId, err := idFunc(in.RegQuesEdit, "вопроса")
		if err != nil {
			return nil
		}
		return editQuestion(w, req, sessUs, quesId)
	}
	err := errors.New("Неправильный путь для редактирования edit ")
	return err
}
