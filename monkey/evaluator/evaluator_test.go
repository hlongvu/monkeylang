package evaluator

import (
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"github.com/hlongvu/monkeylang/monkey/object"
	"github.com/hlongvu/monkeylang/monkey/paser"
	"testing"
)

func testEval(input string) object.Object{
	l := lexer.New(input)
	p := paser.New(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	return Eval(program, env)
}

func testNullObject(t *testing.T, object object.Object) bool{
	return object == NULL
}