package paser

import (
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue interface{}
	}{
		{"return 5;", 5},
		{"return x;", "x"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements not contain 1 statements. Got %d", len(program.Statements))
		}

		stmt := program.Statements[0]

		if !testReturnStatement(t, stmt, tt.expectedValue){
			return
		}
	}
}
