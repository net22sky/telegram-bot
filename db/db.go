// db/db.go
package db

import (
	"database/sql"
	"errors"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var dbConn *sql.DB

// InitDB инициализирует подключение к базе данных
func InitDB(dataSourceName string) error {
	var err error
	dbConn, err = sql.Open("mysql", dataSourceName+"?parseTime=true")
	if err != nil {
		return err
	}

	err = dbConn.Ping()
	if err != nil {
		return err
	}

	log.Println("Подключено к базе данных")
	return nil
}

// Exec выполняет запрос без возвращаемых результатов
func Exec(query string, args ...interface{}) error {
	_, err := dbConn.Exec(query, args...)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return err
	}
	return nil
}

// QueryRow выполняет запрос и возвращает одну строку
func QueryRow(query string, args ...interface{}) (*sql.Row, error) {
	if dbConn == nil {
		return nil, errors.New("соединение с базой данных не установлено")
	}
	return dbConn.QueryRow(query, args...), nil
}

// Query выполняет запрос и возвращает несколько строк
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	if dbConn == nil {
		return nil, errors.New("соединение с базой данных не установлено")
	}
	rows, err := dbConn.Query(query, args...)
	if err != nil {
		log.Printf("Ошибка при выполнении запроса: %v", err)
		return nil, err
	}
	return rows, nil
}
