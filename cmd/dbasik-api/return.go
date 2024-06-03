package main

import (
	"archive/zip"
	"fmt"
	"path/filepath"

	"github.com/tealeg/xlsx/v3"
)

type FilePreparer interface {
	Prepare() ([]string, error)
}

type FileSource struct {
	FilePath string
}

type DirectoryFilePackage struct {
	FileSource
}

type ReturnLine struct {
	Sheet   string
	CellRef string
	Value   string
}

type Return struct {
	Name        string
	ReturnLines []ReturnLine
}

type ZipFilePackage struct {
	FileSource
}

// NewReturnLine creates a new ReturnLine object
func NewReturnLine(sheet, cellRef, value string) (*ReturnLine, error) {
	if err := validateInputs(sheet, cellRef, value); err != nil {
		return nil, err
	}

	if !validateSpreadsheetCell(cellRef) {
		return nil, fmt.Errorf("cellRef must be A1 format")
	}

	return &ReturnLine{
		Sheet:   sheet,
		CellRef: cellRef,
		Value:   value,
	}, nil
}

func validateInputs(sheet, cellRef, value string) error {
	if sheet == "" {
		return fmt.Errorf("sheet parameter is required")
	}
	if cellRef == "" {
		return fmt.Errorf("cellRef parameter is required")
	}
	if value == "" {
		return fmt.Errorf("value parameter is required")
	}
	return nil
}

func NewReturn(name string, dm *Datamap, returnLines []ReturnLine) (*Return, error) {
	if len(returnLines) == 0 {
		return nil, fmt.Errorf("ReturnLines must contain at least one ReturnLine")
	}

	return &Return{
		Name:        name,
		ReturnLines: returnLines,
	}, nil
}

func ParseXLSX(filePath string, dm *Datamap) (*Return, error) {
	// Use tealeg/xlsx to parse the Excel file
	wb, err := xlsx.OpenFile(filePath)
	if err != nil {
		return nil, err
	}

	// Get the set of sheets from the Datamap
	sheets := GetSheetsFromDM(*dm)

	// Loop through all DatamapLines
	returnLines := []ReturnLine{}
	for _, dml := range dm.DMLs {
		// Check if the sheet for this DatamapLine is in the set of sheets
		if !contains(sheets, dml.Sheet) {
			continue
		}

		sh, ok := wb.Sheet[dml.Sheet]
		if !ok {
			return nil, fmt.Errorf("sheet %s not found in Excel file", dml.Sheet)
		}

		col, row, err := xlsx.GetCoordsFromCellIDString(dml.CellRef)
		if err != nil {
			return nil, err
		}
		cell, err := sh.Cell(row, col)
		if err != nil {
			return nil, err
		}
		returnLines = append(returnLines, ReturnLine{
			Sheet:   dml.Sheet,
			CellRef: dml.CellRef,
			Value:   cell.Value, // or cell.FormattedValue() if you need formatted values
		})
	}

	// Here we create a new Return object with the name of the Excel file and the ReturnLines slice
	// that we just populated
	rtn, err := NewReturn(filepath.Base(filePath), dm, returnLines)
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}

func PrepareFiles(fp FilePreparer) ([]string, error) {
	ch := make(chan string, 100)

	go func() {
		defer close(ch)
		files, err := fp.Prepare()
		if err != nil {
			ch <- err.Error()
		}

		for _, f := range files {
			ch <- f
		}
	}()

	var files []string
	for f := range ch {
		files = append(files, f)
	}

	return files, nil
}

func (fp *DirectoryFilePackage) Prepare() ([]string, error) {
	files, err := filepath.Glob(fp.FilePath + "/*")
	if err != nil {
		return nil, err
	}
	return files, nil
}

func (fp *ZipFilePackage) Prepare() ([]string, error) {
	files, err := zip.OpenReader(fp.FilePath)
	if err != nil {
		return nil, err
	}
	defer files.Close()
	out := []string{}
	for _, file := range files.File {
		out = append(out, file.Name)
	}
	return out, nil
}

// NewDirectoryFilePackage creates a new DirectoryFilePackage object with the given filePath to the directory
func NewDirectoryFilePackage(filePath string) *DirectoryFilePackage {
	return &DirectoryFilePackage{FileSource{FilePath: filePath}}
}

// NewZipFilePackage creates a new ZipFilePackage object with the given filePath to the zip file
func NewZipFilePackage(filePath string) *ZipFilePackage {
	return &ZipFilePackage{FileSource{FilePath: filePath}}
}
