package main

import (
	"Monkey/repl"
	"fmt"
	"os"
	"runtime"
	"time"
)

func main() {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	fmt.Printf("Monkey Beta %s\n", currentTime)
	fmt.Printf("[%v] on %v %v\n", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	repl.Start(os.Stdin, os.Stdout)
}
