package paser

import (
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	tests := []struct {
		input              string
		expectedIdentifier string
		expectedValue      interface{}
	}{
		{"let x = 5;", "x", 5},
		{"let x = true;", "x", true},
		{"let foobar = y", "foobar", "y"},
	}

	for _, tt := range tests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if program == nil {
			t.Fatal("ParseProgram return nil")
		}

		if len(program.Statements) != 1{
			t.Fatalf("Statements count not is 1, got %d", len(program.Statements))
		}
		stmt := program.Statements[0]

		if !testLetStatement(t, stmt,tt.expectedIdentifier){
			return
		}

		val := stmt.(*ast.LetStatement).Value

		if !testLiteralExpression(t, val, tt.expectedValue){
			return
		}

	}
}
