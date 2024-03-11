package main

import (
	"fmt"
	"net/http"
	"strconv"
)

func (app *application) createDatamapHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create new datamaps page")
}

func (app *application) showDatamapHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	app.logger.Info("the id requested", "id", id)
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil || id_int < 1 {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details for datamap %d\n", id_int)
}
