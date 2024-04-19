package main

import (
	"fmt"
	"github.com/tealeg/xlsx/v3"
	"path/filepath"
)

type ReturnLine struct {
	Sheet   string
	CellRef string
	Value   string
}

type Return struct {
	Name        string
	ReturnLines []ReturnLine
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

//func cellVisitor(c *xlsx.Cell) error {
//	value, err := c.FormattedValue()
//	if err != nil {
//		fmt.Println(err.Error())
//	} else {
//		fmt.Printf("Sheet: %s, Cell value: %s\n", c.Row.Sheet.Name, value)
//	}
//	return err
//}

//func rowVisitor(r *xlsx.Row) error {
//	return r.ForEachCell(cellVisitor, xlsx.SkipEmptyCells)
//}

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

// contains checks if a slice contains a given string
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
