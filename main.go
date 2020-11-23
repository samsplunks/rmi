package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/BenJoyenConseil/rmi/index"
)

const FIRST_LINE_OF_DATA int = 2

func main() {
	if len(os.Args) <= 1 {
		log.Fatal("Usage: main.go <search_age>")
	}
	file := "data/people.csv"
	// load the age column and parse values into float64 values
	ageColumn := extractColumn(file, "age")

	// create an index over the age column
	idx := index.New(ageColumn)
	search, _ := strconv.ParseFloat(os.Args[1], 64)
	log.Println("max error is :", idx.MaxErrBound, "; min error is", idx.MinErrBound)

	// search an age and get back its line position inside the file people.csv
	result, err := idx.Lookup(search)
	if err != nil {
		log.Fatalf("There is no entry found for %s inside %s \n", os.Args[1], file)
	}
	lines := []int{}
	for _, l := range result {
		lines = append(lines, l+FIRST_LINE_OF_DATA)
	}
	log.Printf("We found %d entries in the index \n", len(lines))
	log.Printf("People who are %s years old are located at %d inside %s \n", os.Args[1], lines, file)

	// generate plot images and save them
	png, _ := filepath.Abs("assets/plot.png")
	svg, _ := filepath.Abs("assets/plot.svg")
	index.Genplot(idx, ageColumn, png, false)
	index.Genplot(idx, ageColumn, svg, true)
}

func extractColumn(file string, colName string) []float64 {
	csvfile, _ := os.Open(file)
	r := csv.NewReader(csvfile)

	var valuesColumn []float64
	var ageCid int
	var headerLine bool = true
	for {
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if headerLine {
			for i, c := range record {
				if strings.ToLower(c) == colName {
					ageCid = i
				}
			}
			headerLine = false
			continue
		}
		v, _ := strconv.ParseFloat(record[ageCid], 64)
		valuesColumn = append(valuesColumn, v)
	}
	return valuesColumn
}
