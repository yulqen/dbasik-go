package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

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

	// // create a new file on the server
	// outFile, err := os.CreateTemp("", "uploaded_csv")
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// // clean up - we have to do this
	// defer os.Remove(outFile.Name())

	// // copy the uploaded file to the server file
	// _, err = io.Copy(outFile, file)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }

	// var b []byte // this doesn't work
	// _, err = outFile.Read(b)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	// fmt.Fprintf(w, string(b))

	// Read the contents of the file
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// Write the file contents to the response
	w.Header().Set("Content-Type", "text/csv") // Set the appropriate content type
	w.Write(fileBytes)
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
