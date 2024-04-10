// dbasik provides a service with which to convert spreadsheets containing
// data to JSON for further processing.

// Copyright (C) 2024 M R Lemon

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

// Package main is the entry point for dbasik.
// It provides the main functionality for the dbasik application.
package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"
)

const version = "0.0.1"

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
	flag.IntVar(&cfg.port, "port", 5000, "API server port")
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

	logger.Info("starting server", "addr", srv.Addr, "env", cfg.env)

	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
