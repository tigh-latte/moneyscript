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

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{
		input:    "5 + true;",
		expected: "type mismatch: INTEGER + BOOLEAN",
	}, {
		input:    "5 + true; 5;",
		expected: "type mismatch: INTEGER + BOOLEAN",
	}, {
		input:    "-true",
		expected: "unknown operator: -BOOLEAN",
	}, {
		input:    "true + false;",
		expected: "unknown operator: BOOLEAN + BOOLEAN",
	}, {
		input:    "5; true + false; 5;",
		expected: "unknown operator: BOOLEAN + BOOLEAN",
	}, {
		input:    "if (10 > 1) { true + false;}",
		expected: "unknown operator: BOOLEAN + BOOLEAN",
	}, {
		input: `
		if (10 > 1) {
			if ( 10 > 1) {
				return true + false;
			}

			return 1;
		}
		`,
		expected: "unknown operator: BOOLEAN + BOOLEAN",
	}, {
		input:    "foobar",
		expected: "identifier not found: foobar",
	}, {
		input:    `"Hello" - "World"`,
		expected: "unknown operator: STRING - STRING",
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)", evaluated, evaluated)
			continue
		}

		if errObj.Message != test.expected {
			t.Errorf("wrong error message. expected=%q got=%q", test.expected, errObj.Message)
		}
	}
}

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		input:    "let a = 5; a;",
		expected: 5,
	}, {
		input:    "let a = 5 * 5; a;",
		expected: 25,
	}, {
		input:    "let a = 5; let b = a; b;",
		expected: 5,
	}, {
		input:    "let a = 5; let b = a; let c = a + b + 5; c;",
		expected: 15,
	}}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fn(x) { x + 2; };"

	evaluated := testEval(input)
	fn, ok := evaluated.(*object.Function)
	if !ok {
		t.Fatalf("object is not Function. got=%T (%+v)", evaluated, evaluated)
	}

	if len(fn.Parameters) != 1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fn.Parameters)
	}

	if fn.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x'. got=%q", fn.Parameters[0])
	}

	if fn.Body.String() != "(x + 2)" {
		t.Fatalf("body is not \"(x + 2)\". got=%q", fn.Body.String())
	}
}

func TestFunctionApplciation(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{{
		input:    "let identity = fn(x) { x; }; identity(5);",
		expected: 5,
	}, {
		input:    "let identity = fn(x) { return x; }; identity(5);",
		expected: 5,
	}, {
		input:    "let double = fn(x) { return x * 2; }; double(5);",
		expected: 10,
	}, {
		input:    "let add = fn(x, y) { return x + y; }; add(5, 2);",
		expected: 7,
	}, {
		input:    "let add = fn(x, y) { return x + y; }; add(5 + 5, add(5, 5));",
		expected: 20,
	}, {
		input:    "fn(x) { x ;}(5);",
		expected: 5,
	}}

	for _, test := range tests {
		testIntegerObject(t, testEval(test.input), test.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	let newAdder = fn(x) {
		fn(y) { x + y };
	};

	let addTwo = newAdder(2);
	addTwo(2);
	`

	testIntegerObject(t, testEval(input), 4)
}

func TestStringLiteral(t *testing.T) {
	input := `"Hello World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestStringConcatentation(t *testing.T) {
	input := `"Hello" + " " + "World!"`

	evaluated := testEval(input)
	str, ok := evaluated.(*object.String)
	if !ok {
		t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
	}

	if str.Value != "Hello World!" {
		t.Errorf("String has wrong value. got=%q", str.Value)
	}
}

func TestBuiltinFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected any
	}{{
		input:    `len("")`,
		expected: 0,
	}, {
		input:    `len("four")`,
		expected: 4,
	}, {
		input:    `len("hello world")`,
		expected: 11,
	}, {
		input:    `len(1)`,
		expected: "argument to `len` not supported, got INTEGER",
	}, {
		input:    `len("one", "two")`,
		expected: "wrong number of arguments. got=2, want=1",
	}}

	for _, test := range tests {
		evaluated := testEval(test.input)
		switch expected := test.expected.(type) {
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case string:
			errObj, ok := evaluated.(*object.Error)
			if !ok {
				t.Errorf("object is not Error. got=%T (%+v)", errObj, errObj)
				continue
			}
			if errObj.Message != expected {
				t.Errorf("wrong error message. expected=%q, got=%q", expected, errObj.Message)
			}
		}
	}
}

func TestArrayLiterals(t *testing.T) {
	input := "[1, 2* 2, 3 + 3]"

	evaluated := testEval(input)
	result, ok := evaluated.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
	}

	if len(result.Elements) != 3 {
		t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
	}

	testIntegerObject(t, result.Elements[0], 1)
	testIntegerObject(t, result.Elements[1], 4)
	testIntegerObject(t, result.Elements[2], 6)
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != evaluator.Null {
		t.Errorf("object is not Null. got=%T (%#v)", obj, obj)
		return false
	}

	return true
}

func testEval(input string) object.Object {
	return evaluator.Eval(parser.New(lexer.New(input)).ParseProgram(), object.NewEnvironment(nil))
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
