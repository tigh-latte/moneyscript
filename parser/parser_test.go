package parser_test

import (
	"fmt"
	"testing"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
	"git.tigh.dev/tigh-latte/monkeyscript/lexer"
	"git.tigh.dev/tigh-latte/monkeyscript/parser"
)

func TestLetStatements(t *testing.T) {
	tests := []struct {
		input    string
		expIdent string
		expVal   any
	}{{
		input:    "let x = 5;",
		expIdent: "x",
		expVal:   5,
	}, {
		input:    "let y = true;",
		expIdent: "y",
		expVal:   true,
	}, {
		input:    "let foobar = y;",
		expIdent: "foobar",
		expVal:   "y",
	}}

	testLetStatement := func(t *testing.T, s ast.Statement, name string) bool {
		if s.TokenLiteral() != "let" {
			t.Errorf("s.TokenLiteral not 'let'. got=%q", s.TokenLiteral())
			return false
		}

		letStmt, ok := s.(*ast.LetStatement)
		if !ok {
			t.Errorf("s not *ast.LetStatement, got=%T", s)
			return false
		}

		if letStmt.Name.Value != name {
			t.Errorf("letStmt.Name.Value not '%s'. got=%s", name, letStmt.Name.Value)
			return false
		}

		if letStmt.Name.TokenLiteral() != name {
			t.Errorf("letStmt.Name.TokenLiteral() not '%s'. got=%s", name, letStmt.Name.TokenLiteral())
		}

		return true
	}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		if !testLetStatement(t, stmt, test.expIdent) {
			return
		}

		val := stmt.(*ast.LetStatement).Value
		if !testLiteralExpression(t, val, test.expVal) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	tests := []struct {
		input  string
		expVal any
	}{
		{input: "return 5;", expVal: 5},
		{input: "return true;", expVal: true},
		{input: "return foobar;", expVal: "foobar"},
	}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Fatalf("stmt not *ast.ReturnStatement. got=%T", stmt)
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Fatalf("returnStmt.TokenLiteral not 'return', got %q",
				returnStmt.TokenLiteral())
		}

		if testLiteralExpression(t, returnStmt.ReturnValue, test.expVal) {
			return
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := parser.New(lexer.New(input))

	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", prog.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerExpression(t *testing.T) {
	input := "5;"

	p := parser.New(lexer.New(input))
	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if len(prog.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d", len(prog.Statements))
	}

	stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", prog.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %d. got=%s", 5, literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	tests := []struct {
		input       string
		expOperator string
		expVal      any
	}{
		{
			input:       "!5",
			expOperator: "!",
			expVal:      5,
		},
		{
			input:       "-15",
			expOperator: "-",
			expVal:      15,
		},
		{
			input:       "!foobar",
			expOperator: "!",
			expVal:      "foobar",
		},
		{
			input:       "-foobar",
			expOperator: "-",
			expVal:      "foobar",
		},
		{
			input:       "!true",
			expOperator: "!",
			expVal:      true,
		},
		{
			input:       "!false",
			expOperator: "!",
			expVal:      false,
		},
	}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain %d statements. got=%d", 1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statement[0] is not ast.ExpressionStatement. got=%T", prog.Statements[0])
		}

		expression, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if expression.Operator != test.expOperator {
			t.Fatalf("exp.Operator1 is not '%s'. got=%s", test.expOperator, expression.Operator)
		}

		if !testLiteralExpression(t, expression.Right, test.expVal) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  any
		operator   string
		rightValue any
	}{{
		input:      "5 + 5",
		leftValue:  5,
		operator:   "+",
		rightValue: 5,
	}, {
		input:      "5 + 5",
		leftValue:  5,
		operator:   "+",
		rightValue: 5,
	}, {
		input:      "5 * 5",
		leftValue:  5,
		operator:   "*",
		rightValue: 5,
	}, {
		input:      "5 / 5",
		leftValue:  5,
		operator:   "/",
		rightValue: 5,
	}, {
		input:      "5 / 5",
		leftValue:  5,
		operator:   "/",
		rightValue: 5,
	}, {
		input:      "5 > 5",
		leftValue:  5,
		operator:   ">",
		rightValue: 5,
	}, {
		input:      "5 < 5",
		leftValue:  5,
		operator:   "<",
		rightValue: 5,
	}, {
		input:      "5 == 5",
		leftValue:  5,
		operator:   "==",
		rightValue: 5,
	}, {
		input:      "5 != 5",
		leftValue:  5,
		operator:   "!=",
		rightValue: 5,
	}, {
		input:      "foobar + barfoo;",
		leftValue:  "foobar",
		operator:   "+",
		rightValue: "barfoo",
	}, {
		input:      "foobar - barfoo;",
		leftValue:  "foobar",
		operator:   "-",
		rightValue: "barfoo",
	}, {
		input:      "foobar * barfoo;",
		leftValue:  "foobar",
		operator:   "*",
		rightValue: "barfoo",
	}, {
		input:      "foobar / barfoo;",
		leftValue:  "foobar",
		operator:   "/",
		rightValue: "barfoo",
	}, {
		input:      "foobar > barfoo;",
		leftValue:  "foobar",
		operator:   ">",
		rightValue: "barfoo",
	}, {
		input:      "foobar < barfoo;",
		leftValue:  "foobar",
		operator:   "<",
		rightValue: "barfoo",
	}, {
		input:      "foobar == barfoo;",
		leftValue:  "foobar",
		operator:   "==",
		rightValue: "barfoo",
	}, {
		input:      "foobar != barfoo;",
		leftValue:  "foobar",
		operator:   "!=",
		rightValue: "barfoo",
	}, {
		input:      "true == true;",
		leftValue:  true,
		operator:   "==",
		rightValue: true,
	}, {
		input:      "true != false;",
		leftValue:  true,
		operator:   "!=",
		rightValue: false,
	}, {
		input:      "false == false;",
		leftValue:  false,
		operator:   "==",
		rightValue: false,
	}}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		if len(prog.Statements) != 1 {
			t.Fatalf("prog.Statements does not contain %d statemetns. got=%d", 1, len(prog.Statements))
		}

		stmt, ok := prog.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("prog.Statements[0] is not ast.ExpressionStatement. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, stmt.Expression, test.leftValue, test.operator, test.rightValue) {
			return
		}
	}
}

func TestOperatorPrecendenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{{
		input:    "-a * b",
		expected: "((-a) * b)",
	}, {
		input:    "!-a",
		expected: "(!(-a))",
	}, {
		input:    "a + b + c",
		expected: "((a + b) + c)",
	}, {
		input:    "a + b - c",
		expected: "((a + b) - c)",
	}, {
		input:    "a * b * c",
		expected: "((a * b) * c)",
	}, {
		input:    "a * b / c",
		expected: "((a * b) / c)",
	}, {
		input:    "a + b / c",
		expected: "(a + (b / c))",
	}, {
		input:    "a + b * c + d / e - f",
		expected: "(((a + (b * c)) + (d / e)) - f)",
	}, {
		input:    "3 + 4; -5 * 5",
		expected: "(3 + 4)((-5) * 5)",
	}, {
		input:    "5 > 4 == 3 < 4",
		expected: "((5 > 4) == (3 < 4))",
	}, {
		input:    "5 < 4 != 3 > 4",
		expected: "((5 < 4) != (3 > 4))",
	}, {
		input:    "3 + 4 * 5 == 3 * 1 + 4 * 5",
		expected: "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
	}, {
		input:    "true",
		expected: "true",
	}, {
		input:    "false",
		expected: "false",
	}, {
		input:    "3 > 5 == false",
		expected: "((3 > 5) == false)",
	}, {
		input:    "3 < 5 == true",
		expected: "((3 < 5) == true)",
	}, {
		input:    "1 + (2 + 3) + 4",
		expected: "((1 + (2 + 3)) + 4)",
	}, {
		input:    "(5 + 5) * 2",
		expected: "((5 + 5) * 2)",
	}, {
		input:    "2 / (5 + 5)",
		expected: "(2 / (5 + 5))",
	}, {
		input:    "(5 + 5) * 2 * (5 + 5)",
		expected: "(((5 + 5) * 2) * (5 + 5))",
	}, {
		input:    "-(5 + 5)",
		expected: "(-(5 + 5))",
	}, {
		input:    "!(true == true)",
		expected: "(!(true == true))",
	}, {
		input:    "a + add(b * c) + d",
		expected: "((a + add((b * c))) + d)",
	}, {
		input:    "add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
		expected: "add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
	}, {
		input:    "add(a + b + c * d / f + g)",
		expected: "add((((a + b) + ((c * d) / f)) + g))",
	}}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		prog := p.ParseProgram()
		checkParserErrors(t, p)

		actual := prog.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input   string
		expBool bool
	}{{
		input:   "true;",
		expBool: true,
	}, {
		input:   "false;",
		expBool: false,
	}}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != test.expBool {
			t.Errorf("boolean.Value not %t. got=%t", test.expBool,
				boolean.Value)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T", stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fn(x, y) { x + y; }`

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fn() {};", expectedParams: []string{}},
		{input: "fn(x) {};", expectedParams: []string{"x"}},
		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		p := parser.New(lexer.New(tt.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	p := parser.New(lexer.New(input))
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestCallExpressionParameterParsing(t *testing.T) {
	tests := []struct {
		input         string
		expectedIdent string
		expectedArgs  []string
	}{
		{
			input:         "add();",
			expectedIdent: "add",
			expectedArgs:  []string{},
		},
		{
			input:         "add(1);",
			expectedIdent: "add",
			expectedArgs:  []string{"1"},
		},
		{
			input:         "add(1, 2 * 3, 4 + 5);",
			expectedIdent: "add",
			expectedArgs:  []string{"1", "(2 * 3)", "(4 + 5)"},
		},
	}

	for _, test := range tests {
		p := parser.New(lexer.New(test.input))
		program := p.ParseProgram()
		checkParserErrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		exp, ok := stmt.Expression.(*ast.CallExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
				stmt.Expression)
		}

		if !testIdentifier(t, exp.Function, test.expectedIdent) {
			return
		}

		if len(exp.Arguments) != len(test.expectedArgs) {
			t.Fatalf("wrong number of arguments. want=%d, got=%d",
				len(test.expectedArgs), len(exp.Arguments))
		}

		for i, arg := range test.expectedArgs {
			if exp.Arguments[i].String() != arg {
				t.Errorf("argument %d wrong. want=%q, got=%q", i,
					arg, exp.Arguments[i].String())
			}
		}
	}
}

func checkParserErrors(t *testing.T, p *parser.Parser) {
	errs := p.Errors()
	if p.Errors() == nil {
		return
	}

	v, ok := errs.(interface {
		Unwrap() []error
	})
	if !ok {
		t.Errorf("unexpected error type")
	}
	t.Errorf("parser has %d errors", len(v.Unwrap()))

	for _, err := range v.Unwrap() {
		t.Errorf("parser error: %s", err)
	}
	t.FailNow()
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected any) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, v int64) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if i.Value != v {
		t.Errorf("i.Value not %d. got=%d", i.Value, v)
		return false
	}

	if i.TokenLiteral() != fmt.Sprintf("%d", v) {
		t.Errorf("i.TokenLiteral not %d. got=%s", v, i.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value,
			ident.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left any, operator string, right any) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}
