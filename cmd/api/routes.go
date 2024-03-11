package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("POST /v1/datamaps", app.createDatamapHandler)
	mux.HandleFunc("GET /v1/datamaps/{id}", app.showDatamapHandler)
	return mux
}
