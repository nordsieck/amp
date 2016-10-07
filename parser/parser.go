package parser

import "go/token"

type Parser func([][]*Token) [][]*Token

func BasicLit(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.INT, p.tok == token.FLOAT, p.tok == token.IMAG,
			p.tok == token.CHAR, p.tok == token.STRING:
			result = append(result, t)
		}
	}
	return result
}

func IdentifierList(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		if p := pop(&t); p == nil || p.tok != token.IDENT {
			continue
		}
		for {
			newT, p := t, (*Token)(nil)
			if p = pop(&newT); p == nil || p.tok != token.COMMA {
				break
			} else if p = pop(&newT); p == nil || p.tok != token.IDENT {
				break
			}
			t = newT
		}
		result = append(result, t)
	}
	return result
}

func TypeName(ts [][]*Token) [][]*Token {
	result := QualifiedIdent(ts)
	return append(result, tokenParser(ts, token.IDENT)...)
}

func QualifiedIdent(ts [][]*Token) [][]*Token {
	ts = PackageName(ts)
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
}

func PackageName(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		if p := pop(&t); p != nil && p.tok == token.IDENT && p.lit != "_" {
			result = append(result, t)
		}
	}
	return result
}

func tokenParser(ts [][]*Token, tok token.Token) [][]*Token {
	var result [][]*Token
	var p *Token
	for _, t := range ts {
		if p = pop(&t); p != nil && p.tok == tok {
			result = append(result, t)
		}
	}
	return result
}
