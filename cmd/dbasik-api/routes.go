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
package main

import "net/http"

func (app *application) routes() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /v1/healthcheck", app.healthcheckHandler)
	mux.HandleFunc("GET /v1/getdatamap/{id}", app.getJSONForDatamap) // TODO: not yet implemented
	mux.HandleFunc("POST /v1/datamap", app.createDatamapHandler)
	mux.HandleFunc("POST /v1/datamapline", app.createDatamapLine)
	mux.HandleFunc("GET /v1/datamaps/{id}", app.showDatamapHandler)
	return mux
}
