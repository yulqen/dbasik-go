package main

import (
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
)

type datamapLine struct {
	Key     string
	Sheet   string
	Cellref string
}

func (app *application) createDatamapHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10Mb max
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the uploaded file
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// parse the csv
	reader := csv.NewReader(file)
	var datamapLines []datamapLine

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // end of file
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(line) != 3 {
			http.Error(w, "Invalid CSV Format", http.StatusBadRequest)
			return
		}

		datamapLines = append(datamapLines, datamapLine{
			Key:     line[0],
			Sheet:   line[1],
			Cellref: line[2],
		})
	}

	err = app.writeJSON(w, http.StatusOK, datamapLines, nil)
	if err != nil {
		app.logger.Debug("writing out csv", "err", err)
		http.Error(w, "Cannot write output from parsed CSV", http.StatusInternalServerError)
		return
	}

	// fmt.Fprintf(w, "file successfully uploaded")
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
