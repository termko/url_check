package main

import (
	"fmt"
	"ozon_service/scheduler"
	"time"

	_ "github.com/jackc/pgx/stdlib"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "toor"
	dbname   = "test_db"
)

func main() {
	// psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
	// 	"password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	// DB, err := sql.Open("pgx", psqlInfo)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// 	return
	// }
	// defer DB.Close()
	// statement := `CREATE TABLE IF NOT EXISTS checks (
	// 	id SERIAL PRIMARY KEY,
	// 	URL TEXT UNIQUE,
	// 	Interval TIMESTAMP
	// 	);`
	// _, err = DB.Exec(statement)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	// statement = `CREATE TABLE IF NOT EXISTS count (
	// 	id SERIAL PRIMARY KEY,
	// 	URL TEXT,
	// 	time TIMESTAMP,
	// 	status INTEGER);`
	// _, err = DB.Exec(statement)
	// if err != nil {
	// 	fmt.Println(err.Error())
	// }
	fmt.Println("Wow")
	go scheduler.Test()
	fmt.Println("Wah")
	time.Sleep(1000 * time.Second)
}
