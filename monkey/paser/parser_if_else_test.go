package paser

import (
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestIfExpression(t *testing.T){
	input := `if (x < y) { x }`
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

	ifexp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok{
		t.Fatalf("stmt.Expression is not IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, ifexp.Condition, "x","<","y"){
		return
	}

	if len(ifexp.Consequence.Statements) != 1{
		t.Fatalf("Consequence is not 1 statement, got %d", len(ifexp.Consequence.Statements))
	}

	con, ok:= ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifexp.Consequence.Statements[0] is not ExpressionStatement, got %T", ifexp.Consequence.Statements[0])
	}

	if !testIdentifier(t, con.Expression,"x"){
		return
	}

	if ifexp.Alternative != nil{
		t.Errorf("Alternative is not nil, got %+v", ifexp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T){
	input := `if (x < y) { x } else {y}`
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

	ifexp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok{
		t.Fatalf("stmt.Expression is not IfExpression, got %T", stmt.Expression)
	}

	if !testInfixExpression(t, ifexp.Condition, "x","<","y"){
		return
	}

	if len(ifexp.Consequence.Statements) != 1{
		t.Fatalf("Consequence is not 1 statement, got %d", len(ifexp.Consequence.Statements))
	}

	con, ok:= ifexp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("ifexp.Consequence.Statements[0] is not ExpressionStatement, got %T", ifexp.Consequence.Statements[0])
	}

	if !testIdentifier(t, con.Expression,"x"){
		return
	}

	if len(ifexp.Alternative.Statements) != 1{
		t.Fatalf("Alternative is not 1 statement, got %d", len(ifexp.Alternative.Statements))
	}

	alt, ok2 := ifexp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok2 {
		t.Fatalf("ifexp.Alternative.Statements[0] is not ExpressionStatement, got %T", ifexp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alt.Expression,"y"){
		return
	}

}
