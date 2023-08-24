package evaluator_test

import (
	"testing"

	"git.tigh.dev/tigh-latte/monkeyscript/evaluator"
	"git.tigh.dev/tigh-latte/monkeyscript/lexer"
	"git.tigh.dev/tigh-latte/monkeyscript/object"
	"git.tigh.dev/tigh-latte/monkeyscript/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		input:    "5",
		expected: 5,
	}, {
		input:    "10",
		expected: 10,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{
		input:    "true",
		expected: true,
	}, {
		input:    "false",
		expected: false,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func testEval(input string) object.Object {
	return evaluator.Eval(parser.New(lexer.New(input)).ParseProgram())
}

func testIntegerObject(t *testing.T, o object.Object, expected int64) bool {
	result, ok := o.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", o, o)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, o object.Object, expected bool) bool {
	result, ok := o.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", o, o)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
