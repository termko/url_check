package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"ozon_service/repository"
	"ozon_service/server"
	"ozon_service/service"
	"strconv"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/kelseyhightower/envconfig"
	"github.com/robfig/cron/v3"
)

type Specification struct {
	PgHost     string `envconfig:"PGHOST"`
	Port       int    `envconfig:"PORT"`
	PgUser     string `envconfig:"PGUSER"`
	PgPassword string `envconfig:"PGPASSWORD"`
	Database   string `envconfig:"PGDATABASE"`
	PgPort     int    `envconfig:"PGPORT"`
}

func main() {
	var s Specification
	err := envconfig.Process("api", &s)

	if err != nil {
		log.Fatal(err)
	}
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		s.PgHost, s.PgPort, s.PgUser, s.PgPassword, s.Database)
	DB, err := sql.Open("pgx", psqlInfo)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer DB.Close()
	statement := `CREATE TABLE IF NOT EXISTS schedule (
		id SERIAL PRIMARY KEY,
		url TEXT UNIQUE,
		jobID INTEGER
		);`
	_, err = DB.ExecContext(context.Background(), statement)
	if err != nil {
		log.Fatal(err)
	}
	statement = `CREATE TABLE IF NOT EXISTS ping (
		id SERIAL PRIMARY KEY,
		url TEXT,
		pingtime TIMESTAMP,
		statuscode INTEGER,
		available BOOLEAN);`
	_, err = DB.ExecContext(context.Background(), statement)
	if err != nil {
		log.Fatal(err)
	}

	repo := repository.New(DB)
	cr := cron.New()
	cr.Start()
	svc := service.New(repo, cr)

	server := &http.Server{
		Addr:    ":" + strconv.Itoa(s.Port),
		Handler: server.New(svc),
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
