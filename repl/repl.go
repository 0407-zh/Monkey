package repl

import (
	"Monkey/lexer"
	"Monkey/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		_, err := fmt.Fprintf(out, PROMPT)
		if err != nil {
			return
		}
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		if _, err = io.WriteString(out, program.String()); err != nil {
			return
		}
		if _, err = io.WriteString(out, "\n"); err != nil {
			return
		}
	}
}

func printParseErrors(out io.Writer, errors []string) {
	if _, err := io.WriteString(out, " parser errors:\n"); err != nil {
		return
	}
	for _, msg := range errors {
		if _, err := io.WriteString(out, "\t"+msg+"\n"); err != nil {
			return
		}
	}
}
