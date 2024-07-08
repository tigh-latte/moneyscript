package object

import (
	"bytes"
	"strconv"
	"strings"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
)

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string {
	return strconv.FormatInt(i.Value, 10)
}

func (i *Integer) Type() ObjectType {
	return IntegerType
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BooleanType
}

func (b *Boolean) Inspect() string {
	return strconv.FormatBool(b.Value)
}

type Null struct{}

func (null *Null) Type() ObjectType {
	return NullType
}

func (null *Null) Inspect() string {
	return "null"
}

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType {
	return ReturnValueType
}

func (r *ReturnValue) Inspect() string {
	return r.Value.Inspect()
}

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType {
	return ErrorType
}

func (e *Error) Inspect() string {
	return "ERROR: " + e.Message
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType {
	return FunctionType
}

func (f *Function) Inspect() string {
	var bb bytes.Buffer

	params := make([]string, len(f.Parameters))
	for i, param := range f.Parameters {
		params[i] = param.String()
	}

	bb.WriteString("fn (")
	bb.WriteString(strings.Join(params, ", "))
	bb.WriteString(") {\n")
	bb.WriteString(f.Body.String())
	bb.WriteString("\n}")

	return bb.String()
}

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return StringType }
func (s *String) Inspect() string  { return s.Value }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BuiltinType }
func (b *Builtin) Inspect() string  { return "builtin function" }
