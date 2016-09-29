package parser

import "go/token"

type Parser func([]*Token) []*Token

func BasicLit(t []*Token) []*Token {
	switch pop(&t).tok {
	case token.INT, token.FLOAT, token.IMAG, token.CHAR, token.STRING:
		return t
	}
	return nil
}

func Klein(p Parser) Parser {
	return func(t []*Token) []*Token {
		for newT := p(t); newT != nil; newT = p(t) {
			t = newT
		}
		return t
	}
}
