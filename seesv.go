package seesv

// ................................................................................................

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

/*
DelimitedFile holds metadata about the file being worked with, and exposes methods to work
with the file.
*/
type DelimitedFile struct {
	File     *os.File
	RowCount int64
	RowIndex []int64
	Headers  []string
	Size     int64
}

// ................................................................................................

/*
Open initializes access to a delimited file (e.g. CSV file). This includes
stats (file size, number of rows), optionally skipping irrelevant lines at the top of the file,
parsing the column headers, and creating an index of row positions to enable O(1) access to
any part of the file.

Example:
    var csvFile seesv.DelimitedFile
    err := csvFile.Open("/path/to/file.csv", 1, true)

This opens the specified CSV file and specified that the 1st line of the file should be ignored (skipped)
and that the file contains a column header line.

    for row := csvFile.Rows(csvFile.RowCount - 10) {
            ...
    }

This returns the last 10 rows of the file.
*/
func (df *DelimitedFile) Open(filePath string, linesToSkip int, hasHeader bool) error {
	df.RowIndex = make([]int64, 0, 2_000_000)

	stat, err := os.Stat(filePath)
	if err != nil {
		return err
	}

	df.Size = stat.Size()

	df.File, err = os.Open(filePath)

	if err != nil {
		return errors.New("Failed to open source file")
	}

	var offset int64

	// Skip 1st line(s) if required.
	// Some files contain extra, non-header information at the top of the file
	// e.g. file summary info. We don't need these, as they are not part of the
	// row data.
	// The number of lines to skip is supplied in the call to Open(), and can be
	// 0 to N.
	if linesToSkip > 0 {
		offset = skipLines(df.File, linesToSkip)
	}

	// Get headers if required.
	// Some files may not have a column header row. In this case, readHeader() returns
	// no headers and the same offset as supplied.
	if hasHeader {
		df.Headers, offset, err = readHeader(df.File, offset)
		if err != nil {
			return err
		}
	}

	// Scan file looking for line ends to populate row index
	// We're going to use a 1MB buffer to minimize the number of I/O trips to disk.
	bufferSize := 1024 * 1024
	buffer := make([]byte, bufferSize)
	bytesRead := bufferSize
	df.RowCount = 0

	var pos int64 = offset
	var i int64 = 0
	var lastLineBreak int64 = pos
	df.RowIndex = append(df.RowIndex, pos)
	df.File.Seek(pos, 0)

	for {
		bytesRead, err = df.File.Read(buffer)

		if err == io.EOF || bytesRead == 0 {
			break
		}

		i = 0

		for {
			if buffer[i] == '\n' {
				lastLineBreak = pos + i
				df.RowIndex = append(df.RowIndex, lastLineBreak+1)
				df.RowCount++

				// TODO: Send progress to progress channel when implemented
				if df.RowCount == 1 || df.RowCount%100_000 == 0 {
					os.Stdout.WriteString(fmt.Sprintf("\r%dK", df.RowCount/1000))
				}
			}

			i++

			if i >= int64(bytesRead) {
				break
			}
		}

		pos += int64(bytesRead)

		if bytesRead < bufferSize {
			break
		}

	}

	// We need to ensure we count the last line if it does not have a carriage return.
	if pos-lastLineBreak > 2 {
		df.RowCount++
	}

	return nil
}

// ................................................................................................

/*
Row returns the row specified by rowNumber. Row numbers are zero-based, so .Row(0) returns the
first data row of the file.

Row data is returned as an array of strings.

Nil is returned if the specified row number is greater than the number of rows in the file or
there is an error parsing the data.
*/
func (df DelimitedFile) Row(rowNumber int64) []string {
	if rowNumber >= df.RowCount {
		return nil
	}

	scanner := bufio.NewScanner(df.File)
	rowPosition := df.RowIndex[rowNumber]
	df.File.Seek(rowPosition, 0)
	scanner.Scan()
	row, err := parseCsv(scanner.Text())
	if err != nil {
		return nil
	}

	return row
}

// ................................................................................................

/*
Rows returns a stream of data rows starting from line `rowNumber` in the file.

If `rowCount` is specified, that number of rows will be returned, unless the end of the file
is reached.

If -1 is specified for `rowCount` then the stream will start at `rowNumber` and continue until
the end of the file.
*/
func (df DelimitedFile) Rows(rowNumber int64, rowCount int64) <-chan []string {
	rowsChan := make(chan []string)

	if rowNumber >= df.RowCount {
		close(rowsChan)
		return rowsChan
	}

	rowPosition := df.RowIndex[rowNumber]
	df.File.Seek(rowPosition, 0)
	rowsReturned := 0

	reader := bufio.NewReader(df.File)
	csvReader := csv.NewReader(reader)
	csvReader.FieldsPerRecord = len(df.Headers)

	go func() {
		for {
			row, err := csvReader.Read()
			if err != nil {
				break
			}

			rowsChan <- row
			rowsReturned++

			if rowCount > 0 && rowsReturned >= int(rowCount) {
				break
			}
		}

		close(rowsChan)
	}()

	return rowsChan
}

// ................................................................................................
func (df *DelimitedFile) Close() {
	df.File.Close()
	df.File = nil
	df.RowCount = 0
	df.RowIndex = make([]int64, 0)
	df.Headers = make([]string, 0)
	df.Size = 0
}

// ................................................................................................

/*
skipLines moves the file pointer past the end of the number of lines specified so that those
lines are ignored by the rest of the file processing.
*/
func skipLines(file *os.File, linesToSkip int) int64 {
	var pos int64
	lines := linesToSkip
	reader := bufio.NewReader(file)

	for lines > 0 {

		for {
			b, err := reader.ReadByte()
			pos++

			if err == io.EOF {
				return pos
			}

			if b == '\n' {
				break
			}
		}

		lines--
	}

	return pos
}

// ................................................................................................

/*
readHeader parses the column header row starting from the byte position in `startPosition`.
This would be 0 if the file has no extra top lines to be ignored.

If the file has lines to be ignored, the call to `skipLines()` will move past them and return
byte position at which the normal file starts. This position would normally then be passed
to `readHeader()`.
*/
func readHeader(file *os.File, startPosition int64) ([]string, int64, error) {
	pos, err := file.Seek(startPosition, 0)
	if err == io.EOF {
		return nil, startPosition, err
	}

	breader := bufio.NewReader(file)
	buffer := make([]byte, 0)

	for {
		b, err := breader.ReadByte()
		pos++

		if err == io.EOF || b == '\n' {
			break
		}
		buffer = append(buffer, b)
	}

	csvReader := csv.NewReader(strings.NewReader(strings.TrimSpace(string(buffer))))
	headers, err := csvReader.Read()
	if err != nil {
		return nil, pos, err
	}

	return headers, pos, nil
}

// ................................................................................................

// parseCsv converts a raw string to parsed CSV data, and returns a string array.
func parseCsv(rowString string) ([]string, error) {
	csvReader := csv.NewReader(strings.NewReader(rowString))
	row, err := csvReader.Read()

	if err != nil {
		return nil, err
	}

	return row, nil

}
