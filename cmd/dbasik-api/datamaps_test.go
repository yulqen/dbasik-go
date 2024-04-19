package main

import (
	"slices"
	"testing"
	"time"
)

func TestGetSheetsFromDM(t *testing.T) {
	testCases := []struct {
		name     string
		datamap  Datamap
		expected []string
	}{
		{
			name: "Extract unique sheet names",
			datamap: Datamap{
				ID:          0,
				Name:        "Test Name",
				Description: "Test description",
				Created:     time.Now(),
				DMLs: []DatamapLine{
					{
						ID:       1,
						Key:      "Test Key",
						Sheet:    "Test Sheet",
						DataType: "TEXT",
						CellRef:  "A10",
					},
					{
						ID:       2,
						Key:      "Test Key 2",
						Sheet:    "Test Sheet",
						DataType: "TEXT",
						CellRef:  "A11",
					},
					{
						ID:       3,
						Key:      "Test Key 3",
						Sheet:    "Test Sheet 2",
						DataType: "TEXT",
						CellRef:  "A12",
					},
				},
			},
			expected: []string{"Test Sheet", "Test Sheet 2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := GetSheetsFromDM(tc.datamap)
			if !slices.Equal(got, tc.expected) {
				t.Errorf("GetSheetsFromDM(%v) = %v, expected %v", tc.datamap, got, tc.expected)
			}
		})
	}
}
