package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type config struct {
	port int
	env  string
}

type application struct {
	config config
	logger *log.Logger
}

func main() {
	var conf config

	flag.IntVar(&conf.port, "port", 8080, "REST API SERVER PORT")
	flag.Parse()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	app := &application{
		config: conf,
		logger: logger,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/healthcheck", app.healthCheckHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", conf.port),
		Handler:      mux,
		IdleTimeout:  time.Minute,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	logger.Printf("starting server on %s", server.Addr)
	err := server.ListenAndServe()
	logger.Fatal(err)

}
