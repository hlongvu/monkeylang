package paser

import (
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestFunctionLiteral(t *testing.T){
	input := `fn(x,y){x+y;};`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)


	if len(program.Statements) != 1 {
		t.Fatalf("Wrong number of statement, got %d", len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ExpressionStatement, got %T", program.Statements[0])
	}

	fn, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not FunctionLiteral, got %T", stmt.Expression)
	}

	if len(fn.Parameters) != 2{
		t.Fatalf("Function Parameters wrong, want 2, got %d", len(fn.Parameters))
	}

	testLiteralExpression(t, fn.Parameters[0], "x")
	testLiteralExpression(t, fn.Parameters[1], "y")

	if len(fn.Body.Statements) != 1{
		t.Fatalf("Function Body Statments wrong, want 1, got %d", len(fn.Body.Statements))
	}

	bodyStmt, ok := fn.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok{
		t.Fatalf("Body is not ExpressionStatement, got %T", fn.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x","+","y")
}

func TestFunctionParametersParsing(t *testing.T){
	tests := []struct {
		input string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z, t) {};", expectedParams: []string{"x", "y", "z", "t"}},
	}

	for _, tt:= range tests{
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		fn := stmt.Expression.(*ast.FunctionLiteral)

		if len(fn.Parameters) != len(tt.expectedParams){
			t.Errorf("Function params length wrong, got %d, expected %d", len(fn.Parameters), len(tt.expectedParams))
		}

		for i, ident := range tt.expectedParams{
			testLiteralExpression(t, fn.Parameters[i], ident)
		}
	}
}