package main

import (
	"reflect"
	"slices"
	"testing"
	"time"
)

func TestCreateNewReturn(t *testing.T) {
	dm := &Datamap{
		ID:          1,
		Name:        "test name",
		Description: "test description",
		Created:     time.Now(),
		DMLs: []DatamapLine{
			{
				ID:       1,
				Key:      "test key",
				Sheet:    "test sheet",
				DataType: "test datatype",
				CellRef:  "test cellref",
			},
		},
	}
	// Call NewReturn with an empty []ReturnLine slice
	rt, err := NewReturn("test name", dm, []ReturnLine{})
	if err == nil {
		t.Error("Expected an error when passing an empty []ReturnLine slice")
	}

	// Check if the error message is as expected
	expectedErrorMsg := "ReturnLines must contain at least one ReturnLine"
	if err != nil && err.Error() != expectedErrorMsg {
		t.Errorf("Unexpected error message. Expected: %s, Got: %s", expectedErrorMsg, err.Error())
	}

	// Check if the returned Return struct is nil
	if rt != nil {
		t.Error("Expected a nil Return struct when an error occurs")
	}
}

func TestNewReturnLine(t *testing.T) {
	rl, err := NewReturnLine("stabs", "C1", "Knocker")
	if err != nil {
		t.Fatal(err)
	}
	if rl == nil {
		t.Errorf("NewReturnLine() returned nil")
	}
	if rl.Sheet != "stabs" {
		t.Errorf("NewReturnLine() returned wrong sheet")
	}
}

func TestReturnLineCellRefFormat(t *testing.T) {
	_, err := NewReturnLine("stabs", "CC", "Knocker")
	if err != nil {
		if err.Error() != "cellRef must be A1 format" {
			t.Errorf("NewReturnLine() returned wrong error")
		}
	}
}

func TestValidateInputs(t *testing.T) {
	// Happy path
	err := validateInputs("Sheet1", "A1", "value")
	if err != nil {
		t.Errorf("validateInputs failed: %v", err)
	}

	// Missing sheet
	err = validateInputs("", "A1", "value")
	if err == nil {
		t.Error("Expected error for missing sheet")
	}
	if err.Error() != "sheet parameter is required" {
		t.Error("Expected error for missing sheet")
	}
	// Missing cellRef
	err = validateInputs("Sheet1", "", "value")
	if err == nil {
		t.Error("cellRef parameter is required")
	}

	// Missing value
	err = validateInputs("Sheet1", "A1", "")
	if err == nil {
		t.Error("value parameter is required")
	}
}

func TestHelper_validateSpreadsheetCell(t *testing.T) {
	if validateSpreadsheetCell("19") != false {
		t.Errorf("Helper.validateSpreadsheetCell() did not return false")
	}

	if validateSpreadsheetCell("1") != false {
		t.Errorf("Helper.validateSpreadsheetCell() did not return false")
	}

	if validateSpreadsheetCell("A10") != true {
		t.Errorf("Helper.validateSpreadsheetCell() did not return true")
	}
}

func TestParseXLSX(t *testing.T) {
	tests := []struct {
		name     string
		filePath string
		dm       *Datamap
		want     *Return
		wantErr  bool
	}{
		{
			name:     "Valid_Excel_file",
			filePath: "../../testdata/valid_excel.xlsx",
			dm: &Datamap{
				DMLs: []DatamapLine{
					{Sheet: "Sheet1", CellRef: "A1"},
					{Sheet: "Sheet1", CellRef: "B1"},
					{Sheet: "Sheet2", CellRef: "C1"},
				},
			},
			want: &Return{
				Name: "valid_excel.xlsx",
				ReturnLines: []ReturnLine{
					{Sheet: "Sheet1", CellRef: "A1", Value: "Value 1"},
					{Sheet: "Sheet1", CellRef: "B1", Value: "Value 2"},
					{Sheet: "Sheet2", CellRef: "C1", Value: "Value 3"},
				},
			},
			wantErr: false,
		},
		// Add more test cases as needed
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseXLSX(tt.filePath, tt.dm)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseXLSX() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got.Name, tt.want.Name) {
				t.Errorf("ParseXLSX() FileName = %v, want %v", got.Name, tt.want.Name)
			}

			if len(got.ReturnLines) != len(tt.want.ReturnLines) {
				t.Errorf("ParseXLSX() ReturnLines length = %v, want %v", len(got.ReturnLines), len(tt.want.ReturnLines))
				return
			}

			for i := range got.ReturnLines {
				if got.ReturnLines[i].Sheet != tt.want.ReturnLines[i].Sheet {
					t.Errorf("ParseXLSX() ReturnLines[%d].Sheet = %v, want %v", i, got.ReturnLines[i].Sheet, tt.want.ReturnLines[i].Sheet)
				}
				if got.ReturnLines[i].CellRef != tt.want.ReturnLines[i].CellRef {
					t.Errorf("ParseXLSX() ReturnLines[%d].CellRef = %v, want %v", i, got.ReturnLines[i].CellRef, tt.want.ReturnLines[i].CellRef)
				}
				if got.ReturnLines[i].Value != tt.want.ReturnLines[i].Value {
					t.Errorf("ParseXLSX() ReturnLines[%d].Value = %v, want %v", i, got.ReturnLines[i].Value, tt.want.ReturnLines[i].Value)
				}
			}
		})
	}
}

func TestPrepareFiles(t *testing.T) {
	fp := NewDirectoryFilePackage("../../testdata")
	files, err := PrepareFiles(fp)
	if err != nil {
		t.Error(err)
	}
	if !slices.Contains(files, "../../testdata/valid_excel.xlsx") {
		t.Errorf("Prepare() did not return ../../testdata/valid_excel.xlsx")
	}
}

func TestUnzipFiles(t *testing.T) {
	fp := NewZipFilePackage("../../testdata/test.zip")
	files, err := PrepareFiles(fp)
	if err != nil {
		t.Error(err)
	}
	if !slices.Contains(files, "valid_excel.xlsx") {
		t.Errorf("Prepare() did not return test.xlsx")
	}
}
