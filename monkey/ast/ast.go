package ast

import (
	"bytes"
	"github.com/hlongvu/monkeylang/monkey/token"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}


// Program is a Node, contains many statements
type Program struct {
	Statements []Statement
}
func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}
func (p *Program) String() string{
	var out bytes.Buffer
	for _, s := range p.Statements{
		out.WriteString(s.String())
	}
	return out.String()
}


// Identifier is an Expression
type Identifier struct {
	Token token.Token
	Value string
}
func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}
func (i *Identifier) String() string{
	return i.Value
}


// IntegerLiteral is an Expression
type IntegerLiteral struct{
	Token token.Token
	Value int64
}
func (i *IntegerLiteral) expressionNode(){}
func (i *IntegerLiteral) TokenLiteral() string{
	return i.Token.Literal
}
func (i *IntegerLiteral) String() string{
	return i.Token.Literal
}


//PrefixExpression is an Expresstion
type PrefixExpression struct{
	Token token.Token
	Operator string
	Right Expression
}
func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string{
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string{
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

//InfixExpression
type InfixExpression struct {
	Token token.Token
	Left Expression
	Operator string
	Right Expression
}

func (ie *InfixExpression) expressionNode(){}
func (ie *InfixExpression) TokenLiteral() string{
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string{
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("" + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}
