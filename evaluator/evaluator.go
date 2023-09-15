package evaluator

import (
	"fmt"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
	"git.tigh.dev/tigh-latte/monkeyscript/object"
)

var (
	Null  = &object.Null{}
	True  = &object.Boolean{Value: true}
	False = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}
	case *ast.Boolean:
		return evalBoolean(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}

		return evalInfixExpression(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.ReturnStatement:
		val := Eval(node.ReturnValue, env)
		if isError(val) {
			return val
		}
		return &object.ReturnValue{Value: val}
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		return env.Set(node.Name.Value, val)
	case *ast.Identifier:
		val, ok := env.Get(node.Value)
		if !ok {
			return newError("identifier not found: " + node.Value)
		}
		return val
	}

	return nil
}

func evalBoolean(b bool) object.Object {
	if b {
		return True
	}
	return False
}

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement, env)

		if result != nil {
			rt := result.Type()
			if rt == object.ReturnValueType || rt == object.ErrorType {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalExclaimOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	default:
		return newError("unknown operator: " + operator + string(right.Type()))
	}
}

func evalExclaimOperatorExpression(right object.Object) object.Object {
	switch right {
	case True:
		return False
	case False:
		return True
	case Null:
		return True
	default:
		return False
	}
}

func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.IntegerType {
		return newError("unknown operator: -" + string(right.Type()))
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.IntegerType && right.Type() == object.IntegerType:
		return evalIntegerInfixExpression(operator, left, right)
	case operator == "==":
		return evalBoolean(left == right)
	case operator == "!=":
		return evalBoolean(left != right)
	case left.Type() != right.Type():
		return newErrorf("type mismatch: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newErrorf("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, lObj, rObj object.Object) object.Object {
	left := lObj.(*object.Integer).Value
	right := rObj.(*object.Integer).Value
	switch operator {
	case "+":
		return &object.Integer{Value: left + right}
	case "-":
		return &object.Integer{Value: left - right}
	case "*":
		return &object.Integer{Value: left * right}
	case "/":
		return &object.Integer{Value: left / right}
	case "<":
		return evalBoolean(left < right)
	case ">":
		return evalBoolean(left > right)
	case "==":
		return evalBoolean(left == right)
	case "!=":
		return evalBoolean(left != right)
	default:
		return newErrorf("unknown operator: %s %s %s", lObj.Type(), operator, rObj.Type())
	}
}

func evalIfExpression(exp *ast.IfExpression, env *object.Environment) object.Object {
	condition := Eval(exp.Condition, env)
	if isError(condition) {
		return condition
	}

	if truthy(condition) {
		return Eval(exp.Consequence, env)
	} else if exp.Alternative != nil {
		return Eval(exp.Alternative, env)
	} else {
		return Null
	}
}

func truthy(o object.Object) bool {
	switch o {
	case Null:
		return false
	case True:
		return true
	case False:
		return false
	default:
		return true
	}
}

func newError(s string) *object.Error {
	return newErrorf(s)
}

func newErrorf(s string, v ...any) *object.Error {
	return &object.Error{Message: fmt.Sprintf(s, v...)}
}

func isError(o object.Object) bool {
	return o != nil && o.Type() == object.ErrorType
}
