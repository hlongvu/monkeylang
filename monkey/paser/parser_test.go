package paser

import (
	"fmt"
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T) {
	input := `	
		let x = 5;
		let y = 10;
		let foobar = 100019;
`

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatal("ParseProgram return nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Statements does not contain 3 statements, got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

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

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return 10000;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements not contain 3 statements. Got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt is not ReturnStatement. Got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral is not return, got %s", returnStmt.TokenLiteral())
		}

	}
}

func TestIdentifierExpresstion(t *testing.T) {
	input := "foobar;"
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

	testIdentifier(t, stmt.Expression, "foobar")
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"
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

	testIntegerLiteral(t, stmt.Expression, 5)
}

func TestBooleanExpression(t *testing.T) {
	input := "true;"
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

	testBoolean(t, stmt.Expression, true)
}

func TestParsingPrefixExpression(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5", "!", 5},
		{"-15", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
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

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("Operator is not %s, got %s", exp.Operator, tt.operator)
		}

		testLiteralExpression(t, exp.Right, tt.value)
	}
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

func TestParsingInfixExpression(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
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

		testInfixExpression(t, stmt.Expression, tt.leftValue, tt.operator, tt.rightValue)
	}

}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},


		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},


		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Expected %q, got %q", tt.expected, actual)
		}
	}
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)

	if !ok {
		t.Fatalf("exp is not Identifier, got %T", exp)
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
