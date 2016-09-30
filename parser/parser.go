package parser

import "go/token"

type Parser func([]*Token) []*Token

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

func TypeName(t []*Token) []*Token {
	return Or(QualifiedIdent, Basic(token.IDENT))(t)
}

func QualifiedIdent(t []*Token) []*Token {
	return And(PackageName, Basic(token.PERIOD), Basic(token.IDENT))(t)
}

func PackageName(t []*Token) []*Token {
	if p := pop(&t); p == nil || p.tok != token.IDENT || p.lit == "_" {
		return nil
	}
	return t
}

func Klein(p Parser) Parser {
	return func(t []*Token) []*Token {
		for newT := p(t); newT != nil; newT = p(t) {
			t = newT
		}
		return t
	}
}

func And(ps ...Parser) Parser {
	return func(t []*Token) []*Token {
		for _, p := range ps {
			newT := p(t)
			if newT == nil {
				return nil
			}
			t = newT
		}
		return t
	}
}

func Maybe(p Parser) Parser {
	return func(t []*Token) []*Token {
		if newT := p(t); newT != nil {
			return newT
		}
		return t
	}
}

func Or(ps ...Parser) Parser {
	return func(t []*Token) []*Token {
		for _, p := range ps {
			if newT := p(t); newT != nil {
				return newT
			}
		}
		return nil
	}
}

func Basic(tok token.Token) Parser {
	return func(t []*Token) []*Token {
		if p := pop(&t); p == nil || p.tok != tok {
			return nil
		}
		return t
	}
}
