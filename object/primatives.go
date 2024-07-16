package object

import (
	"bytes"
	"hash/fnv"
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

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ArrayType }
func (a *Array) Inspect() string {
	bb := new(bytes.Buffer)

	elems := []string{}

	for _, e := range a.Elements {
		elems = append(elems, e.Inspect())
	}

	bb.WriteByte('[')
	bb.WriteString(strings.Join(elems, ", "))
	bb.WriteByte(']')

	return bb.String()
}

type HashKey struct {
	Type  ObjectType
	Value uint64
}

func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))

	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

type HashPair struct {
	Key   Object
	Value Object
}

type Hash struct {
	Pairs map[HashKey]HashPair
}

func (h *Hash) Type() ObjectType {
	return HashType
}

func (h *Hash) Inspect() string {
	bb := new(bytes.Buffer)

	pairs := make([]string, len(h.Pairs))
	for _, pair := range h.Pairs {
		pairs = append(pairs, pair.Key.Inspect()+":"+pair.Value.Inspect())
	}

	bb.WriteRune('{')
	bb.WriteString(strings.Join(pairs, ", "))
	bb.WriteRune('}')

	return bb.String()
}
