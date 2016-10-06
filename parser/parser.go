package parser

import "go/token"

type Parser func([]*Token) []*Token

type Multi func([][]*Token) [][]*Token

func BasicLit(t []*Token) []*Token {
	p := pop(&t)
	if p == nil {
		return nil
	}
	switch p.tok {
	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
		return t
	}
	return nil
}

func IdentifierList(t []*Token) []*Token {
	return And(
		Basic(token.IDENT),
		Klein(And(
			Basic(token.COMMA), Basic(token.IDENT))))(t)
}

// func TypeName(t []*Token) []*Token {
// 	return Or(QualifiedIdent, Basic(token.IDENT))(t)
// }

func QualifiedIdent(t []*Token) []*Token {
	return And(PackageName, Basic(token.PERIOD), Basic(token.IDENT))(t)
}

func PackageName(t []*Token) []*Token {
	if p := pop(&t); p == nil || p.tok != token.IDENT || p.lit == "_" {
		return nil
	}
	return t
}
