package main

import (
	"database/sql"
	"fmt"
	"net/http"
	"ozon_service/repository"
	"ozon_service/server"
	"ozon_service/service"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/robfig/cron/v3"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "toor"
	dbname   = "test_db"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	DB, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer DB.Close()
	statement := `DROP TABLE IF EXISTS schedule`
	DB.Exec(statement)
	statement = `DROP TABLE IF EXISTS ping`
	DB.Exec(statement)
	statement = `CREATE TABLE schedule (
		id SERIAL PRIMARY KEY,
		url TEXT UNIQUE,
		jobID INTEGER
		);`
	_, err = DB.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
	}
	statement = `CREATE TABLE ping (
		id SERIAL PRIMARY KEY,
		url TEXT,
		pingtime TIMESTAMP,
		statuscode INTEGER,
		available BOOLEAN);`
	_, err = DB.Exec(statement)
	if err != nil {
		fmt.Println(err.Error())
	}

	repo := repository.New(DB)
	cr := cron.New()
	cr.Start()
	svc := service.New(repo, cr)

	server := &http.Server{
		Addr:    ":8080",
		Handler: server.New(svc),
	}
	err = server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
