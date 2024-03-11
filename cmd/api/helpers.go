package main

import (
	"encoding/json"
	"net/http"
)

// writeJSON() helper for sending responses. This takes the destination http.ResponseWriter, the
// HTTP status code to sned, the data to encode to JSON and a header map containing any additional
// HTTP headers we want to include in the response.
func (app *application) writeJSON(w http.ResponseWriter, status int, data any, headers http.Header) error {
	// Encode the data to JSON, returing the error if there was one.
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Append a newline to make it easier tro view in terminal applications
	js = append(js, '\n')

	// We know now that we won't encounter any more errors before writing the response,
	// so it's safe to add any headers that we want to include. We loop through the
	// header map and add each header to the http.ResponseWriter header map.
	// Note that it's okay if the provided header map is nil. Go doesn't throw an error
	// if you wanmt to try to range over (or more generally reader from) a nil map.
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
