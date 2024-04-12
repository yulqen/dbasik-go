package main

import (
	"slices"
	"testing"
	"time"
)

func TestReadDML(t *testing.T) {
	dm := datamap{
		ID:          0,
		Name:        "Test Name",
		Description: "Test description",
		Created:     time.Now(),
		DMLs: []datamapLine{
			{
				ID:       1,
				Key:      "Test Key",
				Sheet:    "Test Sheet",
				DataType: "TEXT",
				Cellref:  "A10",
			},
			{
				ID:       2,
				Key:      "Test Key 2",
				Sheet:    "Test Sheet",
				DataType: "TEXT",
				Cellref:  "A11",
			},
			{
				ID:       3,
				Key:      "Test Key 3",
				Sheet:    "Test Sheet 2",
				DataType: "TEXT",
				Cellref:  "A12",
			},
		},
	}

	got := GetSheetsFromDM(dm)
	if !slices.Contains(got, "Test Sheet") {
		t.Errorf("expected to find Test Sheet in %v but didn't find it", got)
	}
	if !slices.Contains(got, "Test Sheet 2") {
		t.Errorf("expected to find Test Sheet in %v but didn't find it", got)
	}
	if slices.Contains(got, "Test Sheet 3") {
		t.Errorf("expected to find Test Sheet in %v but didn't find it", got)
	}
}
