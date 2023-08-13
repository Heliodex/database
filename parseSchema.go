package main

import (
	"errors"
	"fmt"
	"log"
	"os"
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

var unknownType = Type{Name: " unknown ", Primitive: false}

type Field struct {
	Name         string
	TypeName     string
	Type         Type
	DefaultValue string
	FunctionName string

	Optional bool
	Unique   bool
	List     bool
	Set      bool
	Bind     bool

	Fields []Field
}

var DefaultFields = []Field{
	{
		Name:         "id",
		TypeName:     "String",
		Type:         PrimitiveTypes["String"],
		FunctionName: "uuid",
		Unique:       true,
	},
	{
		Name:         "created",
		TypeName:     "Date",
		Type:         PrimitiveTypes["Date"],
		FunctionName: "now",
	},
	{
		Name:         "updated",
		TypeName:     "Date",
		Type:         PrimitiveTypes["Date"],
		FunctionName: "now",
	},
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

		primitiveType, primitive := PrimitiveTypes[currentField.TypeName]

		if primitive {
			currentField.Type = primitiveType
		} else {
			currentField.Type = unknownType
		}

		var linkFields []string

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
				if primitive {
					// Primitive types can't have bindings either
					log.Println(c.InRed("Invalid Type postfix:"), c.InYellow(string(postfix)))
					log.Fatalln(c.InRed("Schema error: " + c.InUnderline(c.InYellow(currentField.TypeName)+c.InUnderline(c.InRed(" Type cannot have bindings, as it is not an Object type")))))
				}
				currentField.Bind = true
			}
		}

		if len(fieldSplit) > 2 && fieldSplit[2] != "" {
			if fieldSplit[2][0] == '&' {
				// This is a generator function, stating how to
				// generate the default value for this field
				currentField.FunctionName = fieldSplit[2][1:]
			} else {
				currentField.DefaultValue = fieldSplit[2]
			}
		}

		parsedFields = append(parsedFields, currentField)
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

	getTypeByName := func(name string) (Type, error) {
		for _, t := range parsedSchema {
			if t.Name == name {
				return t, nil
			}
		}
		return Type{}, errors.New("Type not found")
	}

	var addTypes func([]Field, int)

	addTypes = func(fields []Field, depth int) {
		for _, f := range fields {
			if f.Type.Name == unknownType.Name {
				// This is a link to another type, find it and add it
				t, err := getTypeByName(f.TypeName)
				if err != nil {
					log.Println(c.InRed("Schema error: " +
						c.InUnderline(c.InRed("Type ")) +
						c.InUnderline(c.InYellow(f.TypeName)) +
						c.InUnderline(c.InRed(" not found "))))

					log.Println(c.InRed("Valid Types are:"))
					for _, t := range parsedSchema {
						log.Println(c.InYellow("\t" + t.Name))
					}
					os.Exit(1)
				}
				f.Type = t

				for _, v := range DefaultFields {
					// Add default Fields to Link, if they don't already exist
					var exists bool
					for _, f := range f.Fields {
						if f.Name == v.Name {
							exists = true
							break
						}
					}
					if !exists {
						f.Fields = append(f.Fields, v)
					}
				}
			}

			if len(f.Type.Values) > 0 && f.DefaultValue != "" {
				// This is an Enum, make sure the default value is valid
				var valid bool
				for _, v := range f.Type.Values {
					if v == f.DefaultValue {
						valid = true
						break
					}
				}
				if !valid {
					log.Println(c.InRed("Schema error: " +
						c.InUnderline(c.InRed("Value ")) +
						c.InUnderline(c.InYellow(f.DefaultValue)) +
						c.InUnderline(c.InRed(" is not a valid default value for Enum Type ")) +
						c.InUnderline(c.InYellow(f.TypeName))))

					log.Println(c.InRed("Valid Values are:"))
					for _, v := range f.Type.Values {
						log.Println(c.InYellow("\t" + v))
					}
					os.Exit(1)
				}
			}

			fmt.Println(strings.Repeat("\t", depth) + f.Name)

			if len(f.Fields) > 0 {
				addTypes(f.Fields, depth+1)
			}
		}
	}

	for _, t := range parsedSchema {
		fmt.Println(t.Name)

		if len(t.Values) == 0 {
			for _, v := range DefaultFields {
				// Add default Fields to Link, if they don't already exist
				var exists bool
				for _, f := range t.Fields {
					if f.Name == v.Name {
						exists = true
						break
					}
				}
				if !exists {
					t.Fields = append(t.Fields, v)
				}
			}
		}

		addTypes(t.Fields, 1)
	}

	return parsedSchema
}
