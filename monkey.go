package main

import (
	"Monkey/repl"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Monkey Beta")
	repl.Start(os.Stdin, os.Stdout)
}
