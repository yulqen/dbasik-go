package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"

// holds the config for our application
// We will read this from the command line flags when we run the application
type config struct {
	port int
	env  string
}

// This application struct holds the dependencies for our HTTP handlers, helpers and
// middleware.
type application struct {
	config config
	logger *slog.Logger
}

func main() {
	// Instance of config
	var cfg config

	// Read the flags into the config struct. Defaults are provided if none given.
	flag.IntVar(&cfg.port, "port", 4000, "API server port")
	flag.StringVar(&cfg.env, "env", "development", "Environment (development|staging|production)")
	flag.Parse()

	// Initialize a new structured logger which writes to stdout
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	// An instance of application struct, containing the config struct and the logger
	app := &application{
		config: cfg,
		logger: logger,
	}

	// Declare an http server which listens provided in the config struct and has
	// sensible timeout settings and writes log messages to the structured logger at
	// Error level.
	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("startings server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
