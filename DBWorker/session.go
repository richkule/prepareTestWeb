package DBWorker

import (
	"context"
	in "github.com/richkule/prepareTestWeb/init"
	"time"
)

// Создает новую запись сессии
func CreateSessId(id *in.SessUs) error {
	sql := `insert into session(sess_id,user_id,last_activity,active) values($1,(select id from login where id = $2),$3,'true')`
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(context.Background(), sql, id.SessId, id.UsId, timeNow)
	return err
}

// Обновляет последнюю активность сессии данный момент времени
func UpdateSessTime(sid in.SessId) error {
	sql := `UPDATE session SET last_activity = $1 WHERE sess_id = $2`
	timeNow := time.Now().Format("2006-01-02 15:04:05")
	_, err := db.Exec(context.Background(), sql, timeNow, sid)
	return err
}

// Деактивирует сессию для данного пользователя
func UpdateSessActivity(sid in.SessId) error {
	sql := `UPDATE session SET active = $1 WHERE sess_id = $2`
	_, err := db.Exec(context.Background(), sql, false, sid)
	return err
}
