package main

import "fmt"

type Enum struct {
	Name   string
	Values []string
}

type Type struct {
	Name   string
	Fields []Field
}

type Field struct {
	Name string
	Type string
}

func parseSchema(schema []Definition) []any {
	fmt.Println("sup")
	// fmt.Println(schema)

	// map of Types and Enums
	parsedSchema := []any{}

	for _, def := range schema {
		if def.Category == "Type" {
			parsedSchema = append(parsedSchema, Type{
				Name: def.Name,
				// todo parse fields
			})
		} else {
			parsedSchema = append(parsedSchema, Enum{
				Name:   def.Name,
				Values: def.Fields,
			})
		}
	}

	for _, def := range parsedSchema {
		fmt.Println(def)
	}

	return parsedSchema
}
