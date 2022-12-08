package ast

import (
	"Monkey/token"
	"bytes"
	"strings"
)

type CallExpression struct {
	Token     token.Token // ( 词法单元
	Function  Expression  // 标识符或函数字面量
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {

}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	var out bytes.Buffer

	var args []string
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}

	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")

	return out.String()
}
