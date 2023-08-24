package repl

import (
	"bufio"
	"fmt"
	"io"

	"git.tigh.dev/tigh-latte/monkeyscript/evaluator"
	"git.tigh.dev/tigh-latte/monkeyscript/lexer"
	"git.tigh.dev/tigh-latte/monkeyscript/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)

		if next := scanner.Scan(); !next {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if p.Errors() != nil {
			fmt.Println(p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program)
		if evaluated == nil {
			return
		}
		io.WriteString(out, evaluated.Inspect())
		io.WriteString(out, "\n")
	}
}
