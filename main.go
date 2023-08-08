package main

import (
	"fmt"
	"os"
	"strings"
)

const v = "0.0.0"

func main() {
	fmt.Println(`
      ╭─▅▅                             ╭─▅▅
      │ ██             ╭─██            │ ██
 ╭─▅███▅██ ╭─▅████▄█▌╭─██████╭─▅████▄█▌│ ██▅███▄  ╭─▅████▄█▌╭▗▇█🮅🮅▆▖ ╭─▅█🮅🮅█▅
╭╯██▘  ▝██╭╯▐█▛   ▜█▌╰─╮ ██─┼╯▐█▛   ▜█▌│ ██▘  ▝██╭╯▐█▛   ▜█▌│▝█▙▄▄▃╯╭╯██▄▄▄▄█▋
│ ██▖  ▗██│ ▐█▙   ▟█▌  │ ██ │ ▐█▙   ▟█▌│ ██▖  ▗██│ ▐█▙   ▟█▌╰╮ 🮃▀▜█▋│ ██🮂🮂🮂🮂🮂
╰╮ ▀███▀██╰╮ 🮄████▀█▌  │ ██ ╰╮ 🮄████▀█▌│ ██▀███▀╯╰╮ 🮄████▀█▌╭▝🮅▆▆█🮅╯╰╮ 🮄█▆▆▆🮅
 ╰───╨──╯  ╰─────╨─╯   ╰──╯  ╰─────╨─╯ ╰──╨───╯   ╰─────╨─╯ ╰─────╯  ╰──────╯
	`)

	if len(os.Args) < 2 {
		fmt.Println("Please specify a command")
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "version", "v", "--version", "-v":
		fmt.Println("Version", v)
	case "help", "h", "--help", "-h":
		fmt.Println("Help")
		help()
	}
}
