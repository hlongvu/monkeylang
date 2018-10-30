package paser

import (
	"fmt"
	"github.com/hlongvu/monkeylang/monkey/ast"
	"testing"
)

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not let, got '%s'", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)

	if !ok {
		t.Errorf("s not LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got '%s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s', got '%s'", name, letStmt.Name)
		return false
	}
	return true
}

func testReturnStatement(t *testing.T, s ast.Statement, name interface{})bool{
	if s.TokenLiteral() != "return" {
		t.Errorf("s.TokenLiteral is not return, got '%s'", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.ReturnStatement)

	if !ok {
		t.Errorf("s not ReturnStatement, got %T", s)
		return false
	}

	if !testLiteralExpression(t, letStmt.ReturnValue, name){
		return false
	}
	return true
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, err := range errors {
		t.Errorf("parser error: %q", err)
	}

	t.FailNow()
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il is not IntegerLiteral, got %T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value is not %d, got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral is not %d, got %s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp is not Identifier, got %T : %+v", exp, exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident value is not %s, got %s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral is not %s, got %s", value, ident.Value)
		return false
	}
	return true
}

func testBoolean(t *testing.T, exp ast.Expression, value bool) bool {
	b, ok := exp.(*ast.Boolean)

	if !ok {
		t.Fatalf("exp is not Identifier, got %T", exp)
		return false
	}

	if b.Value != value {
		t.Errorf("ident value is not %t, got %t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("ident.TokenLiteral is not %t, got %s", value, b.TokenLiteral())
		return false
	}
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBoolean(t, exp, v)
	}

	t.Errorf("Type of exp not handled, got %T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opEx, ok := exp.(*ast.InfixExpression)

	if !ok {
		t.Fatalf("exp is not InfixExpression, got %T(%s)", exp, exp)
		return false
	}
	if !testLiteralExpression(t, opEx.Left, left) {
		return false
	}

	if opEx.Operator != operator {
		t.Fatalf("Operator is not %s, got %s", operator, opEx.Operator)
		return false
	}

	if !testLiteralExpression(t, opEx.Right, right) {
		return false
	}
	return true
}
