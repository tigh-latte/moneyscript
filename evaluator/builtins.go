package evaluator

import "git.tigh.dev/tigh-latte/monkeyscript/object"

var builtins = map[string]*object.Builtin{
	"len": {
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newErrorf("wrong number of arguments. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			default:
				return newErrorf("argument to `len` not supported, got %s", arg.Type())
			}
		},
	},
}
