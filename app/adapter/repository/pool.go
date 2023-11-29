package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/grantjforrester/go-ticket/pkg/config"
)

func NewSQLConnectionPool(config config.Provider) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.GetString("db_host"), config.GetInt("db_port"),
		config.GetString("db_username"), config.GetString("db_password"),
		config.GetString("db_database"))

	sqlDB, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicln(err)
	}
	err = sqlDB.Ping()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Database connected")

	return sqlDB
}
