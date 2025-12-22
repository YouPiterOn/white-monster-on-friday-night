package main

import (
	"fmt"
	"os"

	"youpiteron.dev/white-monster-on-friday-night/internal/cli"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: \n\tcli run <file>\n\tcli repl")
		os.Exit(1)
	}
	switch os.Args[1] {
	case "run":
		if len(os.Args) < 3 {
			fmt.Println("usage: cli run <file>")
			os.Exit(1)
		}
		cli.Run(os.Args[2])
	case "repl":
		cli.REPL()
	default:
		fmt.Println("usage: \n\tcli run <file>\n\tcli repl")
		os.Exit(1)
	}
}
