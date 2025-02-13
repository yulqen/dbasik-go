// dbasik provides a service with which to convert spreadsheets containing
// data to JSON for further processing.

// Copyright (C) 2024 M R Lemon

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.
package main

import (
	"database/sql"
	"errors"
	"time"
)

// ErrRecordNotFound A custom err to return from our Get() method when looking up a Datamap
// that doesn't exist
var (
	ErrRecordNotFound = errors.New("record not found")
)

// A Models struct wraps the DatmapModel. We can add other models to this as
// we progress
type Models struct {
	Datamaps     datamapModel
	DatamapLines datamapLineModel
}

// DatamapLine holds the data parsed from each line of a submitted Datamap CSV file.
// The fields need to be exported otherwise they won't be included when encoding
// the struct to json.
type DatamapLine struct {
	ID       int64  `json:"id"`
	Key      string `json:"key"`
	Sheet    string `json:"sheet"`
	DataType string `json:"datatype"`
	CellRef  string `json:"cellref"`
}

type datamapLineModel struct {
	DB *sql.DB
}

// Datamap includes a slice of DatamapLine objects alongside header metadata
type Datamap struct {
	ID          int64         `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Created     time.Time     `json:"created"`
	DMLs        []DatamapLine `json:"datamap_lines"`
}

type datamapModel struct {
	DB *sql.DB
}

// GetSheetsFromDM extracts a set of sheet names from a Datamap struct
func GetSheetsFromDM(dm Datamap) []string {
	// this is basically how sets are done in Go - see https://www.sohamkamani.com/golang/sets/
	sheets := map[string]struct{}{}
	for _, dml := range dm.DMLs {
		sheets[dml.Sheet] = struct{}{}
	}
	var out []string
	for sheet := range sheets {
		out = append(out, sheet)
	}
	return out
}

func NewModels(db *sql.DB) Models {
	return Models{
		Datamaps:     datamapModel{DB: db},
		DatamapLines: datamapLineModel{DB: db},
	}
}

func (m *datamapLineModel) Insert(dm Datamap, dmls []DatamapLine) (int, error) {
	var datamapID int64
	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	err = m.DB.QueryRow(`INSERT INTO datamaps (name, description, created)
		 VALUES ($1, $2, CURRENT_TIMESTAMP)
		 RETURNING id`, dm.Name, dm.Description).Scan(&datamapID)
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare(`INSERT INTO datamap_lines
				(datamap_id, key, sheet, data_type, cellref)
				VALUES ($1, $2, $3, $4, $5)`)
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	for _, line := range dmls {
		_, err := stmt.Exec(
			int64(datamapID),
			line.Key,
			line.Sheet,
			line.DataType,
			line.CellRef)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}
	tx.Commit()
	return int(datamapID), nil
}
