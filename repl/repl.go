package repl

import (
	"Monkey/evaluator"
	"Monkey/lexer"
	"Monkey/object"
	"Monkey/parser"
	"bufio"
	"fmt"
	"io"
)

const PROMPT = ">>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

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
		l := lexer.New(line) // 词法分析器
		p := parser.New(l)   // 语法分析器

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParseErrors(out, p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env) // 求值器
		if evaluated != nil {
			if _, err = io.WriteString(out, evaluated.Inspect()); err != nil {
				return
			}
			if _, err = io.WriteString(out, "\n"); err != nil {
				return
			}
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
