package ast

import "Monkey/token"

type Identifier struct {
	Token token.Token // IDENT词法单元
	Value string
}

func (i *Identifier) expressionNode() {

}

func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}
