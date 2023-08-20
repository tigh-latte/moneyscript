package ast

import (
	"bytes"
	"fmt"

	"git.tigh.dev/tigh-latte/monkeyscript/token"
)

type Node interface {
	fmt.Stringer
	TokenLiteral() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) String() string {
	bb := new(bytes.Buffer)

	for _, s := range p.Statements {
		bb.WriteString(s.String())
	}

	return bb.String()
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}

	return ""
}

type LetStatement struct {
	Token token.Token // `token.LET`
	Name  *Identifier
	Value Expression
}

func (l *LetStatement) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString(l.TokenLiteral() + " ")
	bb.WriteString(l.Name.String())
	bb.WriteString(" = ")
	if l.Value != nil {
		bb.WriteString(l.Value.String())
	}
	bb.WriteString(";")

	return bb.String()
}

func (l *LetStatement) statementNode() {}
func (l *LetStatement) TokenLiteral() string {
	return l.Token.Literal
}

type Identifier struct {
	Token token.Token // `token.IDENT`
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token // `token.INT`
	Value int64
}

func (i *IntegerLiteral) expressionNode() {}
func (i *IntegerLiteral) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IntegerLiteral) String() string {
	return i.Token.Literal
}

type ReturnStatement struct {
	Token       token.Token // `token.RETURN`
	ReturnValue Expression
}

func (r *ReturnStatement) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString(r.TokenLiteral())
	if r.ReturnValue != nil {
		bb.WriteString(r.ReturnValue.String())
	}

	return bb.String()
}

func (r *ReturnStatement) statementNode() {}
func (r *ReturnStatement) TokenLiteral() string {
	return r.Token.Literal
}

type ExpressionStatement struct {
	Token      token.Token // first token of the expression
	Expression Expression
}

func (e *ExpressionStatement) statementNode() {}
func (e *ExpressionStatement) TokenLiteral() string {
	return e.Token.Literal
}

func (e ExpressionStatement) String() string {
	if e.Expression != nil {
		return e.Expression.String()
	}
	return ""
}

type PrefixExpression struct {
	Token    token.Token // The prefix token [`token.MINUS`, `token.EXCLAIM`]
	Operator string
	Right    Expression
}

func (p *PrefixExpression) expressionNode() {}
func (p *PrefixExpression) TokenLiteral() string {
	return p.Token.Literal
}

func (p *PrefixExpression) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString("(")
	bb.WriteString(p.Operator)
	bb.WriteString(p.Right.String())
	bb.WriteString(")")

	return bb.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (i *InfixExpression) expressionNode() {}
func (i *InfixExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *InfixExpression) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString("(")
	bb.WriteString(i.Left.String())
	bb.WriteString(" " + i.Operator + " ")
	bb.WriteString(i.Right.String())
	bb.WriteString(")")

	return bb.String()
}
