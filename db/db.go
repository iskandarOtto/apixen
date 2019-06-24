package db

import (
	"log"

	"github.com/jmoiron/sqlx"
)

func DBConnect() *sqlx.DB {
	db, err := sqlx.Connect("postgres", "user=docker dbname=testdb sslmode=disable password=secret")
	if err != nil {
		log.Fatalln(err)
		return nil
	}
	if err = db.Ping(); err != nil {
		log.Fatalln(err)
		return nil
	}
	return db
}
