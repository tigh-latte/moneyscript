package object

type ObjectType string

const (
	IntegerType     = "INTEGER"
	BooleanType     = "BOOLEAN"
	NullType        = "NULL"
	ReturnValueType = "RETURN_VALUE"
	ErrorType       = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
