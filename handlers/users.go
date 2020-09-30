package handlers

import (
	"errors"
	db "github.com/richkule/prepareTestWeb/DBWorker"
	hf "github.com/richkule/prepareTestWeb/helpFun"
	in "github.com/richkule/prepareTestWeb/init"
	"net/http"
)

// Выводит страницу входа в систему /login
func Login(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if sessUs.UsId != in.GuestUserId {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}
	err := renderTemplate(w, in.LogPage, nil)
	if err != nil {
		err = errors.New("Ошибка обработки шаблона авторизации login " + err.Error())
	}
	return err
}

// Выводит страницу выхода из системы /logout
func Logout(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if sessUs.UsId == in.GuestUserId {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}
	err := db.UpdateSessActivity(sessUs.SessId)
	if err != nil {
		err = errors.New("Ошибка деактивации сессии logout " + err.Error())
		return err
	}
	_, err = hf.CreateAndSetSess(w, req, nil, in.GuestUserId)
	if err != nil {
		err = errors.New("Ошибка генерации или установки сессии makeHandler " + err.Error())
		return err
	}
	http.Redirect(w, req, `/`, http.StatusFound)
	return nil
}

// Обрабатывает авторизацию
func Auto(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if sessUs.UsId != in.GuestUserId {
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	}

	// Функция, выводящая информацию об ошибке пользователю
	wrongFunc := func(wrong string) error {
		data := in.DataLogin{}
		data.AutoWrong = in.AutoWrong(wrong)
		err := renderTemplate(w, in.LogPage, data)
		if err != nil {
			err = errors.New("Ошибка обработки шаблона авторизации auto " + err.Error())
			return err
		}
		return nil
	}

	if req.Method == "GET" {
		http.Redirect(w, req, `/login`, http.StatusFound)
		return nil
	}
	uName := req.FormValue("logUserName")
	pass := req.FormValue("logPassword")
	if !in.RegPass.MatchString(pass) {
		err := wrongFunc(`Неверный пароль. Пароль может содержать латинские буквы и цифры, а также спецсивмолы -_\.@#\$% и состоять не менее чем из 8 и не более 20 символов`)
		return err
	}
	if !in.RegLog.MatchString(uName) {
		err := wrongFunc("Неверное имя пользователя. Логин должен содержать от 5 до 20 латинских символов или цифр. Также возможны спецсимволы -_.")
		return err
	}
	id, getPass, err := db.GetUserData(in.UserName(uName))
	if err != nil {
		err := wrongFunc("Неверные данные пользователя")
		return err
	}
	hexPass, err := hf.HexMD5(pass)
	if err != nil {
		err := errors.New("Ошибка хеширования auto " + err.Error())
		return err
	}
	if hexPass == getPass {
		err := db.UpdateSessActivity(sessUs.SessId)
		if err != nil {
			err := errors.New("Ошибка деактивации сессии auto " + err.Error())
			return err
		}
		_, err = hf.CreateAndSetSess(w, req, nil, id)
		if err != nil {
			err = errors.New("Ошибка генерации или установки сессии auto " + err.Error())
			return err
		}
		http.Redirect(w, req, `/`, http.StatusFound)
		return nil
	} else {
		err = wrongFunc("Неверные данные пользователя")
	}
	return err
}

// Обрабатывает регистрацию
func Reg(w http.ResponseWriter, req *http.Request, sessUs *in.SessUs) error {
	if sessUs.UsId != in.GuestUserId {
		err := Index(w, req, sessUs)
		return err
	}

	// Функция, выводящая информацию об ошибке пользователю
	wrongFunc := func(wrong string) error {
		data := in.DataLogin{}
		data.RegWrong = in.RegWrong(wrong)
		err := renderTemplate(w, in.LogPage, data)
		if err != nil {
			err = errors.New("Ошибка обработки шаблона авторизации reg " + err.Error())
			return err
		}
		return nil
	}
	strName := req.FormValue("username")
	strPass := req.FormValue("password")
	cPass := req.FormValue("confpassword")
	if !in.RegPass.MatchString(strPass) {
		err := wrongFunc(`Неверный пароль. Пароль может содержать латинские буквы и цифры, а также спецсивмолы -_\.@#\$% и состоять не менее чем из 8 и не более 20 символов`)
		return err
	}
	if !in.RegLog.MatchString(strName) {
		err := wrongFunc("Неверное имя пользователя. Логин должен содержать от 5 до 20 латинских символов или цифр. Также возможны спецсимволы -_.")
		return err
	}
	uName := in.UserName(strName)
	exUs, err := db.CheckUser(uName)
	if err != nil {
		err = errors.New("Ошибка проверки существования пользователя reg " + err.Error())
		return err
	}
	if exUs {
		err = wrongFunc("Пользователь существует")
		return err
	}
	if strPass != cPass {
		err = wrongFunc("Пароли не совпдают")
		return err
	}
	hashPass, err := hf.HexMD5(strPass)
	if err != nil {
		err = errors.New("Ошибка хэширования reg " + err.Error())
		return err
	}
	if err = db.NewUser(uName, hashPass); err != nil {
		err = errors.New("Ошибка создания нового пользователя reg " + err.Error())
		return err
	}
	id, _, err := db.GetUserData(uName)
	if err != nil {
		err = errors.New("Ошибка получения id пользователя пользователя reg " + err.Error())
		return err
	}
	err = db.UpdateSessActivity(sessUs.SessId)
	if err != nil {
		err = errors.New("Ошибка деактивации сессии reg " + err.Error())
		return err
	}
	_, err = hf.CreateAndSetSess(w, req, nil, id)
	if err != nil {
		err = errors.New("Ошибка генерации или установки сессии reg " + err.Error())
		return err
	}
	http.Redirect(w, req, `/`, http.StatusFound)
	return nil
}
