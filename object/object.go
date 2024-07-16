package object

type ObjectType string

const (
	IntegerType     = "INTEGER"
	BooleanType     = "BOOLEAN"
	NullType        = "NULL"
	ReturnValueType = "RETURN_VALUE"
	ErrorType       = "ERROR"
	FunctionType    = "FUNCTION"
	StringType      = "STRING"
	BuiltinType     = "BUILTIN"
	ArrayType       = "ARRAY"
	HashType        = "HASH"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Hashable interface {
	HashKey() HashKey
}
