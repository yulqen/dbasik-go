package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
)

// datamapLine holds the data parsed from each line of a submitted datamap CSV file.
// The fields need to be exported otherwise they won't be included when encoding
// the struct to json.
type datamapLine struct {
	Key      string `json:"key"`
	Sheet    string `json:"sheet"`
	DataType string `json:"datatype"`
	Cellref  string `json:"cellref"`
}

// datamap includes a slice of datamapLine objects alongside header metadata
type datamap struct {
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Created     time.Time     `json:"created"`
	DMLs        []datamapLine `json:"datamap_lines"`
}

func (app *application) createDatamapHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10Mb max
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Get form values
	dmName := r.FormValue("name")
	app.logger.Info("obtain value from form", "name", dmName)
	dmDesc := r.FormValue("description")
	app.logger.Info("obtain value from form", "description", dmDesc)

	// Get the uploaded file and name
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// parse the csv
	reader := csv.NewReader(file)
	var dmls []datamapLine
	var dm datamap

	for {
		line, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break // end of file
			}
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(line) != 4 {
			http.Error(w, "Invalid CSV Format", http.StatusBadRequest)
			return
		}

		dmls = append(dmls, datamapLine{
			Key:      line[0],
			Sheet:    line[1],
			DataType: line[2],
			Cellref:  line[3],
		})
	}
	dm = datamap{Name: dmName, Description: dmDesc, Created: time.Now(), DMLs: dmls}

	err = app.writeJSONPretty(w, http.StatusOK, envelope{"datamap": dm}, nil)
	if err != nil {
		app.logger.Debug("writing out csv", "err", err)
		app.serverErrorResponse(w, r, err)
	}

	// fmt.Fprintf(w, "file successfully uploaded")
}

func (app *application) showDatamapHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	app.logger.Info("the id requested", "id", id)
	id_int, err := strconv.ParseInt(id, 10, 64)
	if err != nil || id_int < 1 {
		app.notFoundResponse(w, r)
	}
	fmt.Fprintf(w, "show the details for datamap %d\n", id_int)
}

func (app *application) createDatamapLine(w http.ResponseWriter, r *http.Request) {
	var input datamapLine
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}
	fmt.Fprintf(w, "%v\n", input)
}
