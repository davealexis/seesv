package main

import (
	"fmt"
	"log"

	"github.com/davealexis/seesv"
)

func main() {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("/path/to/data/large_file.csv", 1, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.File.Close()

	fmt.Println("\nThe file has", csvFile.RowCount, "rows")
	fmt.Println(csvFile.Headers)

	log.Println(csvFile.Row(0))
	log.Println(csvFile.Row(1))
	log.Println(csvFile.Row(2))

	fmt.Println("--- Return last row -----------------")

	log.Println(csvFile.Row(csvFile.RowCount - 1))

	fmt.Println("--- Return last 10 rows -------------")

	r := csvFile.RowCount - 10

	c := 0

	for row := range csvFile.Rows(r, -1) {
		log.Println(row)
		c++
	}

	fmt.Println("Returned: ", c)

	fmt.Println("--- Return 5 rows starting at the 10th to last row ---")

	c = 0

	for row := range csvFile.Rows(r, 5) {
		log.Println(row[:5])
		c++
	}

	fmt.Println("Returned: ", c)
}
