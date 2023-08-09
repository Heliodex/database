package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	c "github.com/TwiN/go-color"
)

type Definition struct {
	Name     string
	Category string
	Fields   []string
}

func readSchema() []Definition {
	data, err := os.ReadFile("SCHEMA")
	if err != nil {
		log.Fatalln(c.InRed("Failed to read database schema: "), err)
	}

	lines := strings.Split(strings.ReplaceAll(string(data), "\r\n", "\n"), "\n")

	currentLine := -1

	schemaError := func(error string) {
		e := max(0, currentLine-2)

		// Show 5 lines of context
		for i := 0; i < 5; i++ {
			log.Println(c.InCyan(fmt.Sprint(e+1)) + "\t" + c.InRed(lines[e]))
			if (e) == currentLine {
				log.Println("\t" + c.InRed(strings.Repeat("^", len(lines[e]))))
			}
			e++
		}

		log.Fatalln(c.InRed("Schema error: " + c.InUnderline(error)))
	}

	// Returns the next non-empty line
	next := func() string {
		currentLine++

		for currentLine < len(lines)-1 && strings.TrimSpace(lines[currentLine]) == "" {
			currentLine++
		}
		if currentLine >= len(lines) {
			return ""
		}

		return lines[currentLine]
	}
	peek := func() string {
		e := currentLine
		s := next()
		currentLine = e
		return s
	}

	defs := []Definition{}

	// Read lines from file
	for {
		line := next()
		if line == "" {
			break
		}

		if line[0] != '\t' {
			// Line is the beginning of a definition
			// Could be either a type or an enum

			name := strings.TrimSpace(line)
			if strings.Contains(name, " ") {
				log.Println(c.InRed("Invalid Type or Enum name:"), c.InYellow(name))
				schemaError("Type or Enum name cannot contain spaces")
			}

			def := Definition{
				Name:     name,
				Category: "unknown",
			}

			// Read the next line
			line = peek()

			if line == "" || line[0] != '\t' {
				// The definition is empty
				// An empty Enum is invalid, so it must be an empty Type
				def.Category = "Type"
			} else {
				if strings.Contains(strings.TrimSpace(line), " ") {
					// The definition is an object Type, with fields
					def.Category = "Type"
				} else {
					// It's a list of enum values
					def.Category = "Enum"
				}

				// Save the contents of the block to parse later
				for peek() != "" && peek()[0] == '\t' {
					def.Fields = append(def.Fields, next()[1:])
				}
			}

			defs = append(defs, def)
		}
	}

	return defs
}
