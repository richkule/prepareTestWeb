package helpFun

import (
	"encoding/hex"
	"errors"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	db "github.com/richkule/prepareTestWeb/DBWorker"
	in "github.com/richkule/prepareTestWeb/init"
	"net/http"
)

// Создает сессию в БД и устанавливает куки у пользователя
func CreateAndSetSess(w http.ResponseWriter, r *http.Request, session *sessions.Session, id in.UsId) (*in.SessUs, error) {
	var err error
	if session == nil {
		session, err = in.Store.Get(r, in.SessionName)
		if err != nil {
			return nil, err
		}
	}
	uid, err := genSessId(id)
	if err != nil {
		return nil, err
	}
	err = db.CreateSessId(uid)
	if err != nil {
		return nil, err
	}
	session.Values[in.CookRowName] = uid
	err = sessions.Save(r, w)
	if err != nil {
		return nil, err
	}
	return uid, nil
}

// Генерирует сессию без коллизий в БД
func genSessId(id in.UsId) (*in.SessUs, error) {
	for true {
		sid := in.SessId(hex.EncodeToString(securecookie.GenerateRandomKey(32)))
		uid, err := db.GetUserId(sid)
		if err != nil {
			return nil, err
		}
		if uid == nil {
			uid = &in.SessUs{SessId: sid, UsId: id}
			p := len(sid)
			_ = p
			return uid, nil
		}
	}
	err := errors.New("Невозможный выход из цикла ")
	return nil, err
}
