package evaluator

import (
	"fmt"

	"git.tigh.dev/tigh-latte/monkeyscript/object"
)

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorf("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newErrorf("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
	"first": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorf("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return Null
				}
				return arg.Elements[0]
			default:
				return newErrorf("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"last": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorf("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return Null
				}
				return arg.Elements[len(arg.Elements)-1]
			default:
				return newErrorf("argument to `first` not supported, got %s", arg.Type())
			}
		},
	},
	"rest": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorf("wrong number of arguments. got=%d, want=1", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				if len(arg.Elements) == 0 {
					return Null
				}
				cpy := make([]object.Object, len(arg.Elements)-1)
				copy(cpy, arg.Elements[1:])
				return &object.Array{Elements: cpy}
			default:
				return newErrorf("argument to `rest` not supported, got %s", arg.Type())
			}
		},
	},
	"push": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newErrorf("wrong number of arguments. got=%d, want=2", len(args))
			}

			switch arg := args[0].(type) {
			case *object.Array:
				cpy := make([]object.Object, len(arg.Elements)+1)
				copy(cpy, arg.Elements)
				cpy[len(arg.Elements)] = args[1]
				return &object.Array{Elements: cpy}
			default:
				return newErrorf("argument to `push` not supported, go %s", arg.Type())
			}
		},
	},
	"puts": {
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}

			return Null
		},
	},
}
