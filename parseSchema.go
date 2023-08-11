package main

import (
	"fmt"
	"log"
	"strings"

	c "github.com/TwiN/go-color"
)

type Type struct {
	Name      string
	Primitive bool
	Fields    []Field
	Values    []string
}

var PrimitiveTypes = map[string]Type{
	"String": {Name: "String", Primitive: true},
	"Int":    {Name: "Int", Primitive: true},
	"Float":  {Name: "Float", Primitive: true},
	"Bool":   {Name: "Bool", Primitive: true},
	"Date":   {Name: "Date", Primitive: true},
}

var unknownType = Type{Name: "Unknown", Primitive: true}

type Field struct {
	Name         string
	TypeName     string
	DefaultValue string
	FunctionName string

	Optional bool
	Unique   bool
	List     bool
	Set      bool
	Bind     bool

	Fields []Field
}

func parseFields(fields []string, indent int) []Field {
	var parsedFields []Field
	var previousField Field

	for i := 0; i < len(fields); i++ {
		currentField := Field{}
		fieldSplit := strings.SplitN(strings.TrimLeft(fields[i], "\t"), " ", 3)

		currentField.Name = fieldSplit[0]

		if len(fieldSplit) < 2 {
			log.Fatalln(c.InRed("Schema error: ") +
				c.InUnderline(c.InRed("Field ")) +
				c.InUnderline(c.InYellow(currentField.Name)) +
				c.InUnderline(c.InRed(" is missing a type, default value, or function")),
			)
		}

		currentField.TypeName = strings.TrimRight(fieldSplit[1], "?!*+#")

		_, primitiveType := PrimitiveTypes[currentField.TypeName]

		var linkFields []string

		fmt.Println(previousField.TypeName, fields[i])

		for len(fields[i])-len(strings.TrimLeft(fields[i], "\t")) > indent {
			// This is the start of a set of Fields on a Link,
			// find them by indentation and parse them into linkFields
			if _, e := PrimitiveTypes[previousField.TypeName]; e {
				// Primitive types can't have links
				log.Fatalln(c.InRed("Schema error: " + c.InUnderline(previousField.TypeName+" Type cannot have links")))
			}

			linkFields = append(linkFields, fields[i])
			i++
		}
		if len(linkFields) > 0 {
			currentField.Fields = parseFields(linkFields, indent+1)
		}

		postfixes := strings.TrimLeft(fieldSplit[1], currentField.TypeName)

		for _, postfix := range postfixes {
			checkInvalidPostfix := func(mutuallyExclusivePostfixes ...string) {
				if strings.ContainsAny(postfixes, strings.Join(mutuallyExclusivePostfixes, "")) {
					log.Println(c.InRed("Invalid Type postfix:"), c.InYellow(string(postfix)))
					log.Fatalln(c.InRed("Schema error: " + c.InUnderline(string(postfix)+" postfix cannot be used with "+strings.Join(mutuallyExclusivePostfixes, " or ")+" postfixes")))
				}
			}

			switch postfix {
			case '?':
				checkInvalidPostfix("+", "*")
				currentField.Optional = true
			case '!':
				checkInvalidPostfix("+", "*")
				currentField.Unique = true
			case '*':
				checkInvalidPostfix("?", "!")
				currentField.List = true
			case '+':
				checkInvalidPostfix("?", "!")
				currentField.Set = true
			case '#':
				if primitiveType {
					// Primitive types can't have bindings either
					log.Println(c.InRed("Invalid Type postfix:"), c.InYellow(string(postfix)))
					log.Fatalln(c.InRed("Schema error: " + c.InUnderline(c.InYellow(currentField.TypeName)+c.InUnderline(c.InRed(" Type cannot have bindings, as it is not an Object type")))))
				}
				currentField.Bind = true
			}
		}
		
		if len(fieldSplit) > 2 {
			if fieldSplit[2][0] == '&' {
				// This is a generator function, stating how to
				// generate the default value for this field
				currentField.FunctionName = fieldSplit[2][1:]
				fmt.Println(currentField.FunctionName)
			} else {
				currentField.DefaultValue = fieldSplit[2]
			}
		}

		previousField = currentField
	}

	return parsedFields
}

func parseSchema(schema []Definition) []Type {
	// Slice of object Types and Enums
	var parsedSchema []Type

	for _, def := range schema {
		if def.Category == "Type" {
			parsedSchema = append(parsedSchema, Type{
				Name:   def.Name,
				Fields: parseFields(def.Fields, 1),
			})
		} else {
			for i, v := range def.Fields {
				def.Fields[i] = strings.TrimLeft(v, "\t") // Remove the leading tab
			}

			parsedSchema = append(parsedSchema, Type{
				Name:   def.Name,
				Values: def.Fields,
			})
		}
	}

	return parsedSchema
}
