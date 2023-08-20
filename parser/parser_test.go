package parser_test

import (
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
	testErrs(t, p)
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
	testErrs(t, p)

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

func testErrs(t *testing.T, p *parser.Parser) {
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
