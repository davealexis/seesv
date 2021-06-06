package seesv_test

import (
	"fmt"
	"log"
	"testing"

	"github.com/davealexis/seesv"
)

func TestOpen(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test.csv", 1, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	if csvFile.RowCount != 5_000 {
		t.Errorf("Expected 5,000 data rows in test.csv. Got %d", csvFile.RowCount)
	}

	hCount := len(csvFile.Headers)
	if hCount != 8 {
		t.Errorf("Expected 8 header columns. Got %d", hCount)
	}

	if csvFile.Headers[0] != "ID" {
		t.Error("Headers seem to be incorrectly parsed.")
	}

	row := csvFile.Row(0)
	if len(row) != len(csvFile.Headers) {
		t.Error("Header count does not match data column count")
	}
}

func TestOpenSkipLines(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test_3_headers.csv", 2, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	if csvFile.RowCount != 4 {
		t.Errorf("Expected 4 data rows in test.csv. Got %d", csvFile.RowCount)
	}

	hCount := len(csvFile.Headers)
	if hCount != 3 {
		t.Errorf("Expected 3 header columns. Got %d", hCount)
	}

	if csvFile.Headers[0] != "ID" {
		t.Error("Headers seem to be incorrectly parsed.")
	}
}

func TestOpenBadFile(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test_bad_file.csv", 0, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	if csvFile.RowCount != 4 {
		t.Errorf("Expected 4 data rows in test.csv. Got %d", csvFile.RowCount)
	}

	hCount := len(csvFile.Headers)
	if hCount != 4 {
		t.Errorf("Expected 4 header columns. Got %d", hCount)
	}

	if csvFile.Headers[0] != "ID" {
		t.Error("Headers seem to be incorrectly parsed.")
	}
}

func TestNoDataInFile(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test_no_data.csv", 0, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	if csvFile.RowCount != 0 {
		t.Errorf("Expected 0 data rows in test.csv. Got %d", csvFile.RowCount)
	}

	hCount := len(csvFile.Headers)
	if hCount != 3 {
		t.Errorf("Expected 3 header columns. Got %d", hCount)
	}
}

func TestScanRows(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test.csv", 1, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	rCount := 0

	var value string

	for row := range csvFile.Rows(0, -1) {
		value = row[0]
		rCount++
	}

	fmt.Println(value)

	if rCount != int(csvFile.RowCount) {
		t.Errorf("Row scan did not produce expected row count. Expected %d. Got %d", csvFile.RowCount, rCount)
	}
}

func TestGetInvalidRow(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test.csv", 0, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	invalidRowNumber := csvFile.RowCount + 10
	row := csvFile.Row(invalidRowNumber)

	// Expect row to be nil
	if row != nil {
		t.Error("Should have returned empty row")
	}

	for row := range csvFile.Rows(invalidRowNumber, -1) {
		// Should not get here
		t.Errorf("Should not have gotten a row: %v", row)
	}
}

func TestGetLastRow(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test.csv", 1, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.Close()

	if csvFile.RowCount != 5000 {
		t.Errorf("Expected %d row. Got %d", 5000, csvFile.RowCount)
	}

	rowToFetch := csvFile.RowCount - 1
	row := csvFile.Row(rowToFetch)

	// Expect row to be nil
	if row == nil {
		t.Error("Should have returned a row")
	}

	if row[0] != "11-1111111" {
		t.Error("Got incorrect data for last row. Expected 11-1111111. Got", row[0])
	}
}
