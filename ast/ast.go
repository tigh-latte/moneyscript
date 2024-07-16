package ast

import (
	"bytes"
	"fmt"
	"strings"

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

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}

func (b *Boolean) String() string {
	return b.Token.Literal
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence *BlockStatement
	Alternative *BlockStatement
}

func (i *IfExpression) expressionNode() {}
func (i *IfExpression) TokenLiteral() string {
	return i.Token.Literal
}

func (i *IfExpression) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString("if")
	bb.WriteString(i.Condition.String())
	bb.WriteString(" ")
	bb.WriteString(i.Consequence.String())

	if i.Alternative != nil {
		bb.WriteString(" else ")
		bb.WriteString(i.Alternative.String())
	}

	return bb.String()
}

type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (b *BlockStatement) statementNode() {}
func (b *BlockStatement) TokenLiteral() string {
	return b.Token.Literal
}

func (b *BlockStatement) String() string {
	bb := new(bytes.Buffer)

	for _, s := range b.Statements {
		bb.WriteString(s.String())
	}

	return bb.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []*Identifier
	Body       *BlockStatement
}

func (f *FunctionLiteral) expressionNode() {}
func (f *FunctionLiteral) TokenLiteral() string {
	return f.Token.Literal
}

func (f *FunctionLiteral) String() string {
	bb := new(bytes.Buffer)

	bb.WriteString(f.TokenLiteral())
	bb.WriteString("(")

	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}

	bb.WriteString(strings.Join(params, ", "))
	bb.WriteString(") {")
	bb.WriteString(f.Body.String())
	bb.WriteString("}")

	return bb.String()
}

type CallExpression struct {
	Token     token.Token // Then '(' token
	Function  Expression  // Identifier or FunctionLiteral
	Arguments []Expression
}

func (c *CallExpression) expressionNode() {}
func (c *CallExpression) TokenLiteral() string {
	return c.Token.Literal
}

func (c *CallExpression) String() string {
	bb := new(bytes.Buffer)

	args := make([]string, len(c.Arguments))
	for i, a := range c.Arguments {
		args[i] = a.String()
	}

	bb.WriteString(c.Function.String())
	bb.WriteRune('(')
	bb.WriteString(strings.Join(args, ", "))
	bb.WriteRune(')')

	return bb.String()
}

type StringLiteral struct {
	Token token.Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return sl.Token.Literal }

type ArrayLiteral struct {
	Token    token.Token
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	bb := new(bytes.Buffer)

	elems := make([]string, 0, len(al.Elements))
	for _, el := range al.Elements {
		elems = append(elems, el.String())
	}

	bb.WriteByte('[')
	bb.WriteString(strings.Join(elems, ", "))
	bb.WriteByte(']')

	return bb.String()
}

type IndexExpression struct {
	Token token.Token
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	bb := new(bytes.Buffer)

	bb.WriteByte('(')
	bb.WriteString(ie.Left.String())
	bb.WriteByte('[')
	bb.WriteString(ie.Index.String())
	bb.WriteString("])")

	return bb.String()
}

type HashLiteral struct {
	Token token.Token // the '{' character
	Pairs map[Expression]Expression
}

func (hl *HashLiteral) expressionNode()      {}
func (hl *HashLiteral) TokenLiteral() string { return hl.Token.Literal }
func (hl *HashLiteral) String() string {
	bb := new(bytes.Buffer)

	pairs := make([]string, len(hl.Pairs))
	for k, v := range hl.Pairs {
		pairs = append(pairs, k.String()+":"+v.String())
	}

	bb.WriteRune('{')
	bb.WriteString(strings.Join(pairs, ", "))
	bb.WriteRune('}')

	return bb.String()
}
