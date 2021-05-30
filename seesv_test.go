package seesv_test

import (
	"davealexis/seesv"
	"fmt"
	"log"
	"testing"
)

func TestOpen(t *testing.T) {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("testdata/test.csv", 0, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.File.Close()

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
	defer csvFile.File.Close()

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
	defer csvFile.File.Close()

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
	defer csvFile.File.Close()

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
	err := csvFile.Open("testdata/test.csv", 0, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.File.Close()

	rCount := 0

	for row := range csvFile.Rows(0, -1) {
		fmt.Println(row)
		rCount++
	}

	if rCount != int(csvFile.RowCount) {
		t.Errorf("Row scan did not produce expecter row count. Expected %d. Got %d", csvFile.RowCount, rCount)
	}
}