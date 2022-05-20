package config

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"log"
	"os"
)

func GetDBPostgresql() *sql.DB {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	db, err := sql.Open("postgres", os.Getenv("URL_Postgresql"))
	if err = db.Ping(); err != nil {
		fmt.Print(err)
	}
	if err != nil {
		fmt.Print(err)
	}
	return db
}
