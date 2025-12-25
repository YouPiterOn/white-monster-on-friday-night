package main

import (
	"fmt"
	"os"
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
		Run(os.Args[2])
	case "repl":
		REPL()
	default:
		fmt.Println("usage: \n\tcli run <file>\n\tcli repl")
		os.Exit(1)
	}
}
