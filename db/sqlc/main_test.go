package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
	"Bankstore/utils"
)

// Декларация переменных
var testQueries *Queries
var testDB *pgxpool.Pool

func TestMain(m *testing.M) {
	// Загружаем настройки из файла app.env
	config, err := utils.LoadConfig("../..") // "../.." - Bankstore directory
	if err != nil {
		log.Fatal("can not read config file", err)
	}
	// Соединение с БД
	testDB, err = pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("can not connect to db", err)
	}
	// Закрываем соединение
	defer testDB.Close()

	// Вызываем конструктор для создания экземпляра типа данных Queries
	testQueries = New(testDB)

	// Запускаем subtest(тесты) и итоговый код выполнения передаем в Exit()
	os.Exit(m.Run())
}
