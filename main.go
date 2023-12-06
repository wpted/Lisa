package main

import (
	"Lisa/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Printf("This is the Lisa programming language! Type 'QUIT' to quit.\n")
	repl.Start(os.Stdin, os.Stdout)
}
