package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"hilmi.dag/internal/data"

	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  string
	}
}

type application struct {
	config config
	logger *log.Logger
	models data.Models
}

func main() {
	var conf config

	flag.IntVar(&conf.port, "port", 8080, "REST API SERVER PORT")
	//test case olduğu için credentialleri env ye almadım
	flag.StringVar(&conf.db.dsn, "db-dsn", "postgres://postgres:1234@localhost/users?sslmode=disable", "PostgreSQL DNS")

	// çok işlem olmadığı için openconsu manuel verip biraz  performans kasalım:D
	flag.IntVar(&conf.db.maxOpenConns, "db-max-open-conns", 10, "PostgreSQL max open connections number")
	flag.IntVar(&conf.db.maxIdleConns, "db-max-idle-conns", 10, "PostgreSQL max open connections number")
	flag.StringVar(&conf.db.maxIdleTime, "db-max-idle-time", "10m", "PostgreSQL max connection idle time")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDatabase(conf)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	logger.Printf("database pool establish")

	app := &application{
		config: conf,
		logger: logger,
		models: data.NewModels(db),
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.port),
		Handler:      app.router(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting server on %s", server.Addr)
	err = server.ListenAndServe()
	logger.Fatal(err)

}

func openDatabase(conf config) (*sql.DB, error) {
	db, err := sql.Open("postgres", conf.db.dsn)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(conf.db.maxOpenConns)
	db.SetMaxIdleConns(conf.db.maxOpenConns)

	duration, err := time.ParseDuration(conf.db.maxIdleTime)
	if err != nil {
		return nil, err
	}

	db.SetConnMaxIdleTime(duration)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}
	return db, nil
}
