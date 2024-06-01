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

import (
	"encoding/json"
	"net/http"
	"regexp"
)

// validateSpreadsheetCell checks that the cellRef is in a valid format
func validateSpreadsheetCell(cellRef string) bool {
	pattern := `^[A-Z]+[1-9][0-9]*$`

	regExp := regexp.MustCompile(pattern)

	return regExp.MatchString(cellRef)
}

// We want this so that our JSON is nested under a key at the top, e.g. "Datamap:"...
type envelope map[string]interface{}

// writeJSON)Pretty() helper for sending responses - pretty prints output. This takes the destination http.ResponseWriter, the
// HTTP status code to send, the data to encode to JSON and a header map containing any additional
// HTTP headers we want to include in the response.
func (app *application) writeJSONPretty(w http.ResponseWriter, status int, data envelope, headers http.Header) error {
	// Encode the data to JSON, returing the error if there was one.
	js, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications
	js = append(js, '\n')

	// We know now that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the
	// header map and add each header to the http.ResponseWriter header map.
	// Note that it's okay if the provided header map is nil. Go doesn't throw an error
	// if you want to try to range over (or more generally reader from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

// writeJSON() helper for sending responses. This takes the destination http.ResponseWriter, the
// HTTP status code to send, the data to encode to JSON and a header map containing any additional
// HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode the data to JSON, returing the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier to view in terminal applications
	js = append(js, '\n')

	// We know now that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the
	// header map and add each header to the http.ResponseWriter header map.
	// Note that it's okay if the provided header map is nil. Go doesn't throw an error
	// if you want to try to range over (or more generally reader from) a nil map.
	for key, value := range headers {
		w.Header()[key] = value
	}
	// Add the "Content-Type: application/json" header, then write the status code and
	// JSON response.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}
