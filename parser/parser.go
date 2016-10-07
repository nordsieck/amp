package parser

import "go/token"

type Parser func([][]*Token) [][]*Token

func Expression(ts [][]*Token) [][]*Token {
	return UnaryExpr(ts)
	// Expression binary_op Expression
}

func UnaryExpr(ts [][]*Token) [][]*Token {
	return PrimaryExpr(ts)
	// unary_op UnaryExpr
}

func PrimaryExpr(ts [][]*Token) [][]*Token {
	return Operand(ts)
	// Conversion
	// PrimaryExpr Selector
	// PrimaryExpr Index
	// PrimaryExpr Slice
	// PrimaryExpr TypeAssertion
	// PrimaryExpr Arguments
}

func Operand(ts [][]*Token) [][]*Token {
	xp := tokenParser(ts, token.LPAREN)
	if len(xp) > 0 {
		xp = Expression(xp)
	}
	xp = tokenParser(xp, token.RPAREN)

	return append(
		append(Literal(ts), OperandName(ts)...),
		xp...)
	// MethodExpr

}

func OperandName(ts [][]*Token) [][]*Token {
	return append(tokenParser(ts, token.IDENT), QualifiedIdent(ts)...)
}

func Literal(ts [][]*Token) [][]*Token {
	return BasicLit(ts)
	// CompositeLit
	// FunctionLit
}

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

func Type(ts [][]*Token) [][]*Token {
	a := TypeName(ts)
	// TypeLit
	b := tokenParser(ts, token.LPAREN)
	if len(b) != 0 {
		b = Type(b)
	}
	b = tokenParser(b, token.RPAREN)
	return append(a, b...)
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

func UnaryOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.ADD, p.tok == token.SUB, p.tok == token.NOT, p.tok == token.XOR,
			p.tok == token.MUL, p.tok == token.AND, p.tok == token.ARROW:
			result = append(result, t)
		}
	}
	return result
}

func MulOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.MUL, p.tok == token.QUO, p.tok == token.REM, p.tok == token.SHL,
			p.tok == token.SHR, p.tok == token.AND, p.tok == token.AND_NOT:
			result = append(result, t)
		}
	}
	return result
}

func AddOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.ADD, p.tok == token.SUB, p.tok == token.OR, p.tok == token.AND:
			result = append(result, t)
		}
	}
	return result
}

func RelOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.EQL, p.tok == token.NEQ, p.tok == token.LSS,
			p.tok == token.LEQ, p.tok == token.GTR, p.tok == token.GEQ:
			result = append(result, t)
		}
	}
	return result
}

func BinaryOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.LAND, p.tok == token.LOR:
			result = append(result, t)
		}
	}
	result = append(result, RelOp(ts)...)
	result = append(result, AddOp(ts)...)
	return append(result, MulOp(ts)...)
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
