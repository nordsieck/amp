package parser

import "go/token"

type Parser func([]*Token) []*Token

type Multi func([][]*Token) [][]*Token

func BasicLit(t [][]*Token) [][]*Token { return Each(basicLit)(t) }
func basicLit(t []*Token) []*Token {
	switch p := pop(&t); true {
	case p == nil:
		return nil
	case p.tok == token.INT, p.tok == token.FLOAT, p.tok == token.IMAG,
		p.tok == token.CHAR, p.tok == token.STRING:
		return t
	}
	return nil
}

func IdentifierList(t [][]*Token) [][]*Token { return Each(identifierList)(t) }
func identifierList(t []*Token) []*Token {
	return and(
		Basic(token.IDENT),
		klein(and(
			Basic(token.COMMA), Basic(token.IDENT))))(t)
}

func TypeName(ts [][]*Token) [][]*Token {
	return Or(QualifiedIdent, Each(Basic(token.IDENT)))(ts)
}

func QualifiedIdent(t [][]*Token) [][]*Token { return Each(qualifiedIdent)(t) }
func qualifiedIdent(t []*Token) []*Token {
	return and(packageName, Basic(token.PERIOD), Basic(token.IDENT))(t)
}

func PackageName(t [][]*Token) [][]*Token { return Each(packageName)(t) }
func packageName(t []*Token) []*Token {
	if p := pop(&t); p == nil || p.tok != token.IDENT || p.lit == "_" {
		return nil
	}
	return t
}
