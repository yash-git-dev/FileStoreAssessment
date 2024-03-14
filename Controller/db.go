package controller

import (
	"database/sql"
	"fmt"
	"os"
)

func ConnectDB() (*sql.DB, error) {
	database := os.Getenv("DB_DATABASE")
	port := os.Getenv("DB_PORT")
	host := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	connectionReq := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, password, host, port, database)
	fmt.Println(connectionReq)
	db, err := sql.Open("mysql", connectionReq)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
