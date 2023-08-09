package main

import (
	"fmt"
	"log"
)

func start() {
	schema := schema()

	log.Print("Starting database server...")
	fmt.Println(schema)
}
