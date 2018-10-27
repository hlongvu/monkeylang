package paser

import (
	"github.com/hlongvu/monkeylang/monkey/ast"
	"github.com/hlongvu/monkeylang/monkey/lexer"
	"testing"
)

func TestLetStatement(t *testing.T)  {
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

	if len(program.Statements) != 3{
		t.Fatalf("Statements does not contain 3 statements, got %d", len(program.Statements))
	}

	tests := []struct{
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i,tt := range tests{
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier){
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

func checkParserErrors(t *testing.T, p *Parser){
	errors := p.Errors()
	if len(errors) == 0{
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, err := range errors{
		t.Errorf("parser error: %q", err)
	}

	t.FailNow()
}


func TestReturnStatement(t *testing.T){
	input := `
	return 5;
	return 10;
	return 10000;
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3{
		t.Fatalf("program.Statements not contain 3 statements. Got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements{
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