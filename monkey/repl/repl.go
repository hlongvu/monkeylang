package repl

import (
	"bufio"
	"fmt"
	"github.com/hlongvu/monkeylang/monkey/evaluator"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"github.com/hlongvu/monkeylang/monkey/object"
	"github.com/hlongvu/monkeylang/monkey/paser"
	"io"
)

const PROMPT string = ">> "
func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for{
		fmt.Printf(PROMPT)
		scanned := scanner.Scan()
		if !scanned{
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := paser.New(l)
		program := p.ParseProgram()

		if len(p.Errors()) != 0{
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string){
	for _, msg := range errors{
		io.WriteString(out, "\t" + msg + "\n")
	}
}
