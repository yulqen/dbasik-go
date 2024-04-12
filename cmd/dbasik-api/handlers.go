package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/tealeg/xlsx/v3"
)

func (app *application) createReturnHandler(w http.ResponseWriter, r *http.Request) {
	// Parse the multipart form
	err := r.ParseMultipartForm(10 << 20) // 10Mb max
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	// Get form values
	dmName := r.FormValue("name")
	dmDesc := r.FormValue("description")

	// Get the return file and save it to a file
	returnFile, handler, err := r.FormFile("returnfile")
	app.logger.Info("got excel file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer returnFile.Close()
	dst, err := os.Create(handler.Filename)
	if err != nil {
		http.Error(w, "Cannot create new file object from uplaoded file.", http.StatusInternalServerError)
		return
	}
	defer dst.Close()
	if _, err := io.Copy(dst, returnFile); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the uploaded csv file and name
	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Missing file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// parse the csv
	reader := csv.NewReader(file)
	var dmls []datamapLine

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
	dm := datamap{Name: dmName, Description: dmDesc, Created: time.Now(), DMLs: dmls}

	// Parse the XLSX file based on the datamap...
	// open an existing file
	wb, err := xlsx.OpenFile(dst.Name())
	if err != nil {
		http.Error(w, "Cannot open Excel file", http.StatusBadRequest)
		return
	}
	// wb now contains a reference to the workbook
	// show all the sheets in the workbook
	fmt.Println("Sheets in this file:")
	for i, sh := range wb.Sheets {
		fmt.Println(i, sh.Name)
	}
	fmt.Println("----")

	//Here is where we parse our files.

	err = app.writeJSONPretty(w, http.StatusOK, envelope{"datamap": dm}, nil)
	if err != nil {
		app.logger.Debug("writing out csv", "err", err)
		app.serverErrorResponse(w, r, err)
	}

	fmt.Fprintf(w, "file successfully uploaded")

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

	fmt.Fprintf(w, "file successfully uploaded")
}

func (app *application) saveDatamapHandler(w http.ResponseWriter, r *http.Request) {
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
			ID:       0,
			Key:      line[0],
			Sheet:    line[1],
			DataType: line[2],
			Cellref:  line[3],
		})

	}
	dm = datamap{Name: dmName, Description: dmDesc, Created: time.Now(), DMLs: dmls}

	// save to the database
	_, err = app.models.DatamapLines.Insert(dm, dmls)
	if err != nil {
		http.Error(w, "Cannot save to database", http.StatusBadRequest)
		return
	}

}

func (app *application) getJSONForDatamap(w http.ResponseWriter, r *http.Request) {
	// Get the DM out of the database
	// dm = datamap{Name: dmName, Description: dmDesc, Created: time.Now(), DMLs: dmls}

	// err = app.writeJSONPretty(w, http.StatusOK, envelope{"datamap": dm}, nil)
	// if err != nil {
	// 	app.logger.Debug("writing out csv", "err", err)
	// 	app.serverErrorResponse(w, r, err)
	// }

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
