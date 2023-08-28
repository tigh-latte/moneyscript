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
	}, {
		input:    "-5",
		expected: -5,
	}, {
		input:    "-10",
		expected: -10,
	}, {
		input:    "5 + 5 + 5 + 5 -10",
		expected: 10,
	}, {
		input:    "2 * 2 * 2 * 2 * 2",
		expected: 32,
	}, {
		input:    "-50 + 100 + -50",
		expected: 0,
	}, {
		input:    "5 * 2 + 10",
		expected: 20,
	}, {
		input:    "5 + 2 * 10",
		expected: 25,
	}, {
		input:    "20 + 2 * -10",
		expected: 0,
	}, {
		input:    "50 / 2 * 2 + 10",
		expected: 60,
	}, {
		input:    "2 * (5 + 10)",
		expected: 30,
	}, {
		input:    "3 * 3 * 3 + 10",
		expected: 37,
	}, {
		input:    "3 * (3 * 3) + 10",
		expected: 37,
	}, {
		input:    "(5 + 10 * 2 + 15 / 3) * 2 + -10",
		expected: 50,
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
	}, {
		input:    "1 < 2",
		expected: true,
	}, {
		input:    "1 > 2",
		expected: false,
	}, {
		input:    "1 < 1",
		expected: false,
	}, {
		input:    "1 > 1",
		expected: false,
	}, {
		input:    "1 == 1",
		expected: true,
	}, {
		input:    "1 != 1",
		expected: false,
	}, {
		input:    "1 == 2",
		expected: false,
	}, {
		input:    "1 != 2",
		expected: true,
	}, {
		input:    "true == true",
		expected: true,
	}, {
		input:    "false == false",
		expected: true,
	}, {
		input:    "true == false",
		expected: false,
	}, {
		input:    "true != false",
		expected: true,
	}, {
		input:    "false != true",
		expected: true,
	}, {
		input:    "(1 < 2) == true",
		expected: true,
	}, {
		input:    "(1 < 2) == false",
		expected: false,
	}, {
		input:    "(1 > 2) == true",
		expected: false,
	}, {
		input:    "(1 > 2) == false",
		expected: true,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestExclaimOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{{
		input:    "!true",
		expected: false,
	}, {
		input:    "!false",
		expected: true,
	}, {
		input:    "!5",
		expected: false,
	}, {
		input:    "!!true",
		expected: true,
	}, {
		input:    "!!false",
		expected: false,
	}, {
		input:    "!!5",
		expected: true,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testBooleanObject(t, evaluated, test.expected)
	}
}

func TestIfElseExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{{
		input:    "if (true) { 10 }",
		expected: 10,
	}, {
		input:    "if (false) { 10 }",
		expected: nil,
	}, {
		input:    "if (1) { 10 }",
		expected: 10,
	}, {
		input:    "if (1 < 2) { 10 }",
		expected: 10,
	}, {
		input:    "if (1 > 2) { 10 }",
		expected: nil,
	}, {
		input:    "if (1 > 2) { 10 } else { 20 }",
		expected: 20,
	}, {
		input:    "if (1 < 2) { 10 } else { 20 }",
		expected: 10,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		integer, ok := test.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		input:    "return 10;",
		expected: 10,
	}, {
		input:    "return 10; 9;",
		expected: 10,
	}, {
		input:    "return 2 * 5; 9;",
		expected: 10,
	}, {
		input:    "9; return 2 * 5; 9;",
		expected: 10,
	}, {
		input:    "if (10 > 1) { if (10 > 1) { return 10; } return 1; }",
		expected: 10,
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		testIntegerObject(t, evaluated, test.expected)
	}
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.Null {
		t.Errorf("object is not Null. got=%T (%#v)", obj, obj)
		return false
	}

	return true
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
		t.Errorf("object is not Boolean. got=%T (%+v)", o, o)
		return false
	}

	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}
