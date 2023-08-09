package main

import (
	// "fmt"
	"log"
)

func start() {
	log.Println("Reading database schema...")
	schema := readSchema()

	log.Println("Parsing database schema...")
	// parsedSchema :=
	parseSchema(schema)

	// fmt.Println(parsedSchema)

	log.Println("Starting database server...")
}
