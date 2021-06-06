package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/davealexis/seesv"
	"github.com/eiannone/keyboard"
	"github.com/olekukonko/tablewriter"
)

func main() {
	var csvFile seesv.DelimitedFile

	DisplayRows := 20

	err := csvFile.Open("/path/to/large.csv", 0, true)
	if err != nil {
		log.Fatal(err)
	}
	defer csvFile.Close()

	fmt.Println()
	log.Println(csvFile.RowCount, " rows")

	currentRow := int64(0)
	displayRows(&csvFile, currentRow, DisplayRows, 0)

	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer func() {
		_ = keyboard.Close()
	}()

	stop := false
	columnShift := 0
	columnCount := len(csvFile.Headers) - 1

	for stop == false {
		fmt.Print(">> ")
		char, key, err := keyboard.GetKey()
		if err != nil {
			panic(err)
		}

		switch key {
		case keyboard.KeyPgup:
			if currentRow < int64(DisplayRows) {
				currentRow = 0
			} else {
				currentRow -= int64(DisplayRows)
			}
		case keyboard.KeyPgdn:
			currentRow += int64(DisplayRows)

			if currentRow >= csvFile.RowCount {
				currentRow = csvFile.RowCount - int64(DisplayRows)
			}
		case keyboard.KeyArrowRight:
			columnShift += 2
			if columnShift+15 > columnCount {
				columnShift = columnCount - 15
			}
		case keyboard.KeyArrowLeft:
			columnShift -= 2
			if columnShift < 0 {
				columnShift = 0
			}
		case keyboard.KeyEsc:
			fmt.Print("\033[H\033[2J")
			stop = true
		case 0:
			switch char {
			case 'g':
				currentRow = 0
			case 'G':
				currentRow = csvFile.RowCount - int64(DisplayRows)
			case '/':
				fmt.Print("/")
				var input string
				fmt.Scan(&input)
				input = strings.TrimSuffix(input, "\r\n")
				var line int64
				line, err := strconv.ParseInt(strings.TrimSuffix(input, "\r\n"), 10, 64)
				fmt.Println("Go to:", line)

				if err == nil {
					currentRow = int64(line)
					if currentRow >= csvFile.RowCount {
						currentRow = csvFile.RowCount - int64(DisplayRows)
					}
				}
			}
		}

		if stop == false {
			displayRows(&csvFile, currentRow, DisplayRows, columnShift)
		}
	}
}

func displayRows(csv *seesv.DelimitedFile, start int64, rows int, colShift int) {
	rowNum := start
	end := start + int64(rows)
	colCount := 15

	fmt.Print("\033[H\033[2J")
	fmt.Println("Press ESC to quit")
	fmt.Println(start, "->", end, "of", csv.RowCount)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetRowSeparator("-")
	table.SetRowLine(true)
	rowHeaders := append([]string{"Row"}, csv.Headers[colShift:colCount+colShift]...)

	table.SetHeader(rowHeaders)

	for v := range csv.Rows(start, -1) {
		table.Append(append([]string{strconv.FormatInt(rowNum, 10)}, v[colShift:colCount+colShift]...))

		rowNum++

		if rowNum >= end {
			break
		}
	}

	table.Render()
}
