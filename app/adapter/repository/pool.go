package repository

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/grantjforrester/go-ticket/pkg/config"
)

func NewSqlConnectionPool(config config.Provider) *sql.DB {
	connectionString := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		config.GetString("db_host"), config.GetInt("db_port"),
		config.GetString("db_username"), config.GetString("db_password"),
		config.GetString("db_database"))

	sqlDb, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Panicln(err)
	}
	err = sqlDb.Ping()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("Database connected")

	return sqlDb
}
