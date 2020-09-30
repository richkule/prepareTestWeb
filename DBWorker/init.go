package DBWorker

import (
	"context"
	"github.com/jackc/pgx/pgxpool"
)

const (
	dsn = "user=postgres password=postgres host=127.0.0.1 port=5432 dbname=postgres sslmode=disable pool_max_conns=10"
)

var db *pgxpool.Pool // Переменная подключения к БД

// Осуществляющая соединение с БД
func DbConnect() error {
	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return err
	}
	db, err = pgxpool.ConnectConfig(context.Background(), poolConfig)
	if err != nil {
		return err
	}
	return nil
}
