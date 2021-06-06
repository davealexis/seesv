package main

import (
	"fmt"
	"log"

	"github.com/davealexis/seesv"
)

func main() {
	var csvFile seesv.DelimitedFile
	err := csvFile.Open("c:/data/billing.202103.csv", 1, true)
	if err != nil {
		log.Fatal("Failed to open file")
	}
	defer csvFile.File.Close()

	fmt.Println("\nThe file has", csvFile.RowCount, "rows")

	fmt.Println("--- Headers -------------------------")
	fmt.Println(csvFile.Headers[:4])

	fmt.Println("--- First 3 rows --------------------")
	log.Println(csvFile.Row(0)[0:4])
	log.Println(csvFile.Row(1)[0:4])
	log.Println(csvFile.Row(2)[0:4])

	fmt.Println("--- Last row ------------------------")
	log.Println(csvFile.Row(csvFile.RowCount - 1))

	fmt.Println("--- Last 10 rows --------------------")

	r := csvFile.RowCount - 10

	c := 0

	for row := range csvFile.Rows(r, -1) {
		log.Println(row[0:4])
		c++
	}

	fmt.Println("Returned: ", c)

	fmt.Println("--- Return 5 rows starting at the 10th to last row ---")

	c = 0

	for row := range csvFile.Rows(r, 5) {
		log.Println(row[:4])
		c++
	}

	fmt.Println("Returned: ", c)
}
