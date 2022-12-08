package ast

import (
	"Monkey/token"
	"bytes"
)

type BlockStatement struct {
	Token      token.Token // { 词法单元
	Statements []Statement
}

func (bs *BlockStatement) expressionNode() {

}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out bytes.Buffer

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
