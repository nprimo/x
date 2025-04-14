package main

import (
	"encoding/csv"
	"fmt"
	"os"

	csvtag "github.com/nprimo/csv-tag"
)

func main() {
	f, err := os.Open("example.csv")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rows, err := csv.NewReader(f).ReadAll()
	if err != nil {
		panic(err)
	}

	type A struct {
		Id   int     `csv:"a"`
		Word string  `csv:"b"`
		Num  float32 `csv:"c"`
	}
	var data []A
	if err := csvtag.Unmarshal(rows, &data); err != nil {
		panic(err)
	}

	fmt.Println("-- Original")
	for i, el := range data {
		fmt.Printf("row %d: %+v\n", i, el)
	}

	data[0].Word = "guarda"
	data[1].Word = "come"
	fmt.Println("-- Modified")
	for i, el := range data {
		fmt.Printf("row %d: %+v\n", i, el)
	}

	rows, err = csvtag.Marshall(data)
	if err != nil {
		panic(err)
	}

	f1, err := os.Create("example_modified.csv")
	if err != nil {
		panic(err)
	}
	defer f1.Close()
	if err := csv.NewWriter(f1).WriteAll(rows); err != nil {
		panic(err)
	}
}
