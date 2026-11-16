package main

import (
	"fmt"
	"os"
	"strings"

	c "github.com/TwiN/go-color"
)

const v = "0.0.0"

func main() {
	// This says "database" in unicode art I promise
	fmt.Print(c.InRed("      ╭─▅▅"), "                             ", c.InPurple("╭─▅▅")+"\n")
	fmt.Print(c.InRed("      │ ██"), "           ", c.InGreen("  ╭─██  "), "          ", c.InPurple("│ ██ Version "), c.InRed(v)+"\n")
	fmt.Print(c.InRed(" ╭─▅███▅██"), c.InYellow(" ╭─▅████▄█▌"), c.InGreen("╭─██████"), c.InBlue("╭─▅████▄█▌"), c.InPurple("│ ██▅███▄ "), c.InRed(" ╭─▅████▄█▌"), c.InYellow("╭▗▇█🮅█▆▖"), c.InGreen("╭─▅█🮅🮅█▅")+"\n")
	fmt.Print(c.InRed("╭╯██▘─╮▝██"), c.InYellow("╭╯▐█▛─╮ ▜█▌"), c.InGreen("╰─╮ ██─┼"), c.InBlue("╯▐█▛─╮ ▜█▌"), c.InPurple("│ ██▘─╮▝██"), c.InRed("╭╯▐█▛─╮ ▜█▌"), c.InYellow("│▝█▙▄▃╮"), c.InGreen(" ╭╯██▄▄▄▄█▋")+"\n")
	fmt.Print(c.InRed("│ ██▖ │▗██"), c.InYellow("│ ▐█▙ │ ▟█▌"), c.InGreen("  │ ██ "), c.InBlue("│ ▐█▙ │ ▟█▌"), c.InPurple("│ ██▖ │▗██"), c.InRed("│ ▐█▙ │ ▟█▌"), c.InYellow("╰╮ ▀▀▜█▖"), c.InGreen("│ ██🮂🮂🮂🮂🮂")+"\n")
	fmt.Print(c.InRed("╰╮ ▀███▀██"), c.InYellow("╰╮ ▀████▀█▌"), c.InGreen("  │ ██ "), c.InBlue("╰╮ ▀████▀█▌"), c.InPurple("│ ██▀███▀╯"), c.InRed("╰╮ ▀████▀█▌")+c.InYellow("╭▝🮅▆▆█▀"), c.InGreen("╰╮ 🮄█▆▆▆🮅")+"\n")
	fmt.Print(c.InRed(" ╰───╨──╯ "), c.InYellow(" ╰─────╨─╯ "), c.InGreen("  ╰──╯  "), c.InBlue("╰─────╨─╯"), c.InPurple(" ╰──╨───╯  "), c.InRed(" ╰─────╨─╯ "), c.InYellow("╰─────╯"), c.InGreen("  ╰──────╯")+"\n")

	if len(os.Args) < 2 {
		fmt.Println(c.InRed("No command specified. Run 'database help' for more information."))
		os.Exit(1)
	}

	switch strings.ToLower(os.Args[1]) {
	case "--version", "-v", "v", "version":
		fmt.Println("Version", v)
	case "--help", "-h", "h", "help":
		help(v)
	case "--start", "-s", "s", "start":
		start()
	default:
		fmt.Println(c.InRed("Unknown command. Run 'database help' for more information."))
		os.Exit(1)
	}
}
