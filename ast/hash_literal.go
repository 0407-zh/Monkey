package ast

import (
	"Monkey/token"
	"bytes"
	"strings"
)

type HashLiteral struct {
	Token token.Token // '{'词法单元
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode() {

}

func (hl *HashLiteral) TokenLiteral() string {
	return hl.Token.Literal
}

func (hl *HashLiteral) String() string {
	var out bytes.Buffer

	var pairs []string
	for key, value := range hl.Pairs {
		pairs = append(pairs, key.String()+":"+value.String())
	}

	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")

	return out.String()
}
