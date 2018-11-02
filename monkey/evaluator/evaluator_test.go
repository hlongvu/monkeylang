package evaluator

import (
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"github.com/hlongvu/monkeylang/monkey/object"
	"github.com/hlongvu/monkeylang/monkey/paser"
)

func testEval(input string) object.Object{
	l := lexer.New(input)
	p := paser.New(l)
	program := p.ParseProgram()
	return Eval(program)
}