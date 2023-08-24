package object

import "strconv"

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
