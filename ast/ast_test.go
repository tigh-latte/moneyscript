package ast_test

import (
	"testing"

	"git.tigh.dev/tigh-latte/monkeyscript/ast"
	"git.tigh.dev/tigh-latte/monkeyscript/token"
)

func TestString(t *testing.T) {
	prog := &ast.Program{
		Statements: []ast.Statement{
			&ast.LetStatement{
				Token: token.Token{Type: token.LET, Literal: "let"},
				Name: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "myVar"},
					Value: "myVar",
				},
				Value: &ast.Identifier{
					Token: token.Token{Type: token.IDENT, Literal: "anotherVar"},
					Value: "anotherVar",
				},
			},
		},
	}

	if str := prog.String(); str != "let myVar = anotherVar;" {
		t.Errorf("program.String() wrong. got=%q", str)
	}
}
