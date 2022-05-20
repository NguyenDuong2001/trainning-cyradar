package config

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func GetDBMySql() *sql.DB {
	e := godotenv.Load()
	if e != nil {
		log.Fatal(e)
	}
	db, err := sql.Open("mysql", os.Getenv("URL_MySQL"))
	if err = db.Ping(); err != nil {
		fmt.Print(err)
	}
	if err != nil {
		fmt.Print(err)
	}
	return db
}
