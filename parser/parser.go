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
