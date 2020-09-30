package DBWorker

import (
	"context"
	in "github.com/richkule/prepareTestWeb/init"
)

// Возвращает идентефикатор пользователя, если его сессия активна и существует
func GetUserId(sid in.SessId) (*in.SessUs, error) {
	var uid in.SessUs
	var isActive bool
	res := db.QueryRow(context.Background(), `select * from session where sess_id = $1 and active = 'true'`, sid)
	err := res.Scan(nil, &uid.SessId, &uid.UsId, nil, &isActive)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, nil
		}
		return nil, err
	}
	if isActive {
		return &uid, nil
	} else {
		return nil, nil
	}
}

// Проверяет существует ли пользователь с таким логином
func CheckUser(name in.UserName) (bool, error) {
	res, err := db.Exec(context.Background(), `select * from login where login = $1`, name)
	if err != nil {
		return false, err
	}
	if res[len(res)-1] == '0' {
		return false, nil
	}
	return true, nil
}

// Создание нового пользователя
func NewUser(name in.UserName, pass string) error {
	_, err := db.Exec(context.Background(), `insert into login values($1,$2)`, name, pass[:])
	return err
}

// Получает данные пользователя
func GetUserData(name in.UserName) (id in.UsId, pass string, err error) {
	res := db.QueryRow(context.Background(), `select * from login where login = $1`, name)
	passArr := make([]byte, 32)
	err = res.Scan(nil, &passArr, &id)
	if err != nil {
		return 0, "", err
	}
	return id, string(passArr[:]), nil
}

// Получает логин пользователя по его id
func GetUserName(id in.UsId) (in.UserName, error) {
	res := db.QueryRow(context.Background(), `select * from login where id = $1`, id)

	var uName in.UserName
	err := res.Scan(&uName, nil, nil)
	if err != nil {
		return "", err
	}
	return uName, nil
}
