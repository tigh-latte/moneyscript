package object

type ObjectType string

const (
	IntegerType = "INTEGER"
	BooleanType = "BOOLEAN"
	NullType    = "NULL"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}
