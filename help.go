package main

import (
	"fmt"

	c "github.com/TwiN/go-color"
)

func help(v string) {
	fmt.Println("database", c.InBlackOverGreen(" "+v+" ")+"\n")
	fmt.Println(c.InYellow("Usage"))
	fmt.Println(c.InGreen("	database [command] [arguments]\n"))
	fmt.Println(c.InYellow("Commands"))
	fmt.Println(c.InBlue("	h help") + "     Shows this help message")
	fmt.Println(c.InBlue("	v version") + "  Shows the software version")
	fmt.Println(c.InBlue("	s start") + "    Starts the database server")
}
