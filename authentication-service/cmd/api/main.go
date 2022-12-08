package main

import (
	"authentication/cmd/api/data"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	//db drivers
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

const webPort = "80"

var counts int8

type Config struct {
	DB     *sql.DB     // connect to database
	Models data.Models // data models
}

func main() {
	log.Println("Starting authentication service")

	// TODO connect to DB       using pgconn pgx, PostgreSQL database drivers.
	conn := connectToDB()
	if conn == nil {
		log.Panic("Can't connect to Postgres!")
	}

	// set up config
	app := Config{
		DB:     conn,
		Models: data.New(conn),
	}

	srv := &http.Server{
		Addr:    ":" + webPort,
		Handler: app.routes(),
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
}

// open the database connection
func openDB(dsn string) (*sql.DB, error) { //dsn connection string for the database (take from env)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping() // if connection error try pinging
	if err != nil {
		return nil, err
	}

	return db, nil
}

func connectToDB() *sql.DB {
	dsn := os.Getenv("DSN")

	for { // try connecting until connect
		connection, err := openDB(dsn)
		if err != nil {
			log.Println("Postgres not yet ready ...")
			counts++
		} else {
			log.Println("Connected to Postgres!")
			return connection
		}

		if counts > 10 { // if continuously getting error connecting return
			log.Println(err)
			return nil
		}

		log.Println("Backing off for two seconds....")
		time.Sleep(3 * time.Second)
		continue
	}
}
