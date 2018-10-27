package ast

import "github.com/hlongvu/monkeylang/monkey/token"

type ReturnStatement struct {
	Token token.Token
	ReturnValue Expression
}
