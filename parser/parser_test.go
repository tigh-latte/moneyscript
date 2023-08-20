package parser_test

import (
	"fmt"
	"testing"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
	"git.tigh.dev/tigh-latte/monkeyscript/lexer"
	"git.tigh.dev/tigh-latte/monkeyscript/parser"
)

func TestLetStatements(t *testing.T) {
	input := `
	let x = 5;
	let y = 10;
	let foobar = 838383;
	`

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

	p := parser.New(lexer.New(input))

	prog := p.ParseProgram()
	checkParserErrors(t, p)
	if prog == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(prog.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(prog.Statements))
	}

	tests := []struct {
		expIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := prog.Statements[i]
		if !testLetStatement(t, stmt, tt.expIdent) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	input := `
	return 5;
	return 10;
	return993322;
	`

	p := parser.New(lexer.New(input))

	prog := p.ParseProgram()
	checkParserErrors(t, p)

	if len(prog.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements got=%d", len(prog.Statements))
	}

	for _, stmt := range prog.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}

		if retStmt.TokenLiteral() != "return" {
			t.Errorf("retStmt.TokenLiteral not 'return', got %q", retStmt.TokenLiteral())
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
		input        string
		operator     string
		integerValue int64
	}{
		{
			input:        "!5",
			operator:     "!",
			integerValue: 5,
		},
		{
			input:        "-15",
			operator:     "-",
			integerValue: 15,
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

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt.Expression is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != test.operator {
			t.Fatalf("exp.Operator1 is not '%s'. got=%s", test.operator, exp.Operator)
		}

		if !checkIntegerLiteral(t, exp.Right, test.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{{
		input:      "5 + 5",
		leftValue:  5,
		operator:   "+",
		rightValue: 5,
	}, {
		input: "5 + 5",

		leftValue:  5,
		operator:   "+",
		rightValue: 5,
	}, {
		input: "5 * 5",

		leftValue:  5,
		operator:   "*",
		rightValue: 5,
	}, {
		input: "5 / 5",

		leftValue:  5,
		operator:   "/",
		rightValue: 5,
	}, {
		input: "5 / 5",

		leftValue:  5,
		operator:   "/",
		rightValue: 5,
	}, {
		input: "5 > 5",

		leftValue:  5,
		operator:   ">",
		rightValue: 5,
	}, {
		input: "5 < 5",

		leftValue:  5,
		operator:   "<",
		rightValue: 5,
	}, {
		input: "5 == 5",

		leftValue:  5,
		operator:   "==",
		rightValue: 5,
	}, {
		input: "5 != 5",

		leftValue:  5,
		operator:   "!=",
		rightValue: 5,
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

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if !checkIntegerLiteral(t, exp.Left, test.leftValue) {
			return
		}

		if exp.Operator != test.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", test.operator, exp.Operator)
		}

		if !checkIntegerLiteral(t, exp.Right, test.rightValue) {
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

func checkIntegerLiteral(t *testing.T, il ast.Expression, v int64) bool {
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
