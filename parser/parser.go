package parser

import (
	"fmt"
	"go/token"
)

type Parser func([][]*Token) [][]*Token

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

func AnonymousField(ts [][]*Token) [][]*Token {
	maybe := make([][]*Token, len(ts))
	for i, t := range ts {
		if p := pop(&t); p == nil || p.tok != token.MUL {
			maybe[i] = ts[i]
		} else {
			maybe[i] = t
		}
	}
	return TypeName(maybe)
}

func Arguments(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LPAREN)
	newTs := ExpressionList(ts) // validate that you understand https://golang.org/ref/spec#Arguments
	temp := make([][]*Token, len(newTs))
	for i, t := range newTs {
		if p := pop(&t); p == nil || p.tok != token.ELLIPSIS {
			temp[i] = newTs[i]
		} else {
			temp[i] = t
		}
	}
	newTs = temp
	temp = make([][]*Token, len(newTs))
	for i, t := range newTs {
		if p := pop(&t); p == nil || p.tok != token.COMMA {
			temp[i] = newTs[i]
		} else {
			temp[i] = t
		}
	}
	ts = append(ts, temp...)
	return tokenParser(ts, token.RPAREN)
}

func ArrayType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	ts = tokenParser(ts, token.RBRACK)
	return Type(ts)
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

func CompositeLit(ts [][]*Token) [][]*Token {
	ts = LiteralType(ts)
	return LiteralValue(ts)
}

func Conversion(ts [][]*Token) [][]*Token {
	ts = Type(ts)
	ts = tokenParser(ts, token.LPAREN)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	temp := make([][]*Token, len(ts))
	for i, t := range ts {
		if p := pop(&t); p != nil && p.tok == token.COMMA {
			temp[i] = t
		} else {
			temp[i] = ts[i]
		}
	}
	return tokenParser(temp, token.RPAREN)
}

func Element(ts [][]*Token) [][]*Token {
	return append(Expression(ts), LiteralValue(ts)...)
}

func ElementList(ts [][]*Token) [][]*Token {
	ts = KeyedElement(ts)
	base := ts
	for len(base) != 0 {
		next := tokenParser(base, token.COMMA)
		next = KeyedElement(next)
		ts = append(ts, next...)
		base = next
	}
	return ts
}

func EllipsisArrayType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	ts = tokenParser(ts, token.ELLIPSIS)
	ts = tokenParser(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func Expression(ts [][]*Token) [][]*Token {
	base := UnaryExpr(ts)
	comp := BinaryOp(base)
	if len(comp) != 0 {
		comp = Expression(comp)
	}
	return append(base, comp...)
}

func ExpressionList(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	for {
		newTs := tokenParser(ts, token.COMMA)
		newTs = Expression(newTs)
		if len(newTs) == 0 {
			break
		}
		ts = newTs
	}
	return ts
}

func FieldDecl(ts [][]*Token) [][]*Token {
	a := IdentifierList(ts)
	if len(a) != 0 {
		a = Type(a)
	}
	ts = append(AnonymousField(ts), a...)
	var withTag [][]*Token
	for _, t := range ts {
		if newT := tokenParser([][]*Token{t}, token.STRING); len(newT) != 0 {
			withTag = append(withTag, newT...)
		}
	}
	return append(ts, withTag...)
}

func FunctionType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.FUNC)
	return Signature(ts)
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

func Index(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	return tokenParser(ts, token.RBRACK)
}

func Key(ts [][]*Token) [][]*Token {
	return append(tokenParser(ts, token.IDENT), Expression(ts)...)
}

func KeyedElement(ts [][]*Token) [][]*Token {
	with := Key(ts)
	with = tokenParser(with, token.COLON)
	return append(Element(ts), Element(with)...)
}

func Literal(ts [][]*Token) [][]*Token {
	return append(BasicLit(ts), CompositeLit(ts)...)
	// FunctionLit
}

func LiteralType(ts [][]*Token) [][]*Token {
	return append(
		append(
			append(StructType(ts), ArrayType(ts)...),
			append(EllipsisArrayType(ts), SliceType(ts)...)...),
		append(MapType(ts), TypeName(ts)...)...)
}

func LiteralValue(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACE)
	if len(ts) != 0 {
		list := ElementList(ts)
		list = append(list, tokenParser(list, token.COMMA)...)
		ts = append(ts, list...)
	}
	return tokenParser(ts, token.RBRACE)
}

func MapType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.MAP)
	ts = tokenParser(ts, token.LBRACK)
	if len(ts) != 0 {
		ts = Type(ts)
	}
	ts = tokenParser(ts, token.RBRACK)
	return Type(ts)
}

func MethodExpr(ts [][]*Token) [][]*Token {
	ts = ReceiverType(ts)
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
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

func PrimaryExpr(ts [][]*Token) [][]*Token {
	base := append(Operand(ts), Conversion(ts)...)
	newBase := base
	for len(newBase) != 0 {
		additions := Selector(newBase)
		additions = append(additions, Index(newBase)...)
		additions = append(additions, Slice(newBase)...)
		additions = append(additions, TypeAssertion(newBase)...)
		additions = append(additions, Arguments(newBase)...)
		base = append(base, additions...)
		newBase = additions
	}
	return base
}

func Operand(ts [][]*Token) [][]*Token {
	xp := tokenParser(ts, token.LPAREN)
	if len(xp) > 0 {
		xp = Expression(xp)
	}
	xp = tokenParser(xp, token.RPAREN)

	return append(
		append(Literal(ts), OperandName(ts)...),
		append(MethodExpr(ts), xp...)...)
}

func OperandName(ts [][]*Token) [][]*Token {
	return append(tokenParser(ts, token.IDENT), QualifiedIdent(ts)...)
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

func ParameterDecl(ts [][]*Token) [][]*Token {
	ts = append(ts, IdentifierList(ts)...)
	ts = append(ts, tokenParser(ts, token.ELLIPSIS)...)
	return Type(ts)
}

func ParameterList(ts [][]*Token) [][]*Token {
	ts = ParameterDecl(ts)
	next := ts
	for len(next) != 0 {
		decl := tokenParser(next, token.COMMA)
		decl = ParameterDecl(decl)
		ts = append(ts, decl...)
		next = decl
	}
	return ts
}

func Parameters(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LPAREN)
	params := ParameterList(ts)
	params = append(params, tokenParser(params, token.COMMA)...)
	ts = append(ts, params...)
	return tokenParser(ts, token.RPAREN)
}

func QualifiedIdent(ts [][]*Token) [][]*Token {
	ts = PackageName(ts)
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
}

func ReceiverType(ts [][]*Token) [][]*Token {
	ptr := tokenParser(ts, token.LPAREN)
	ptr = tokenParser(ptr, token.MUL)
	ptr = TypeName(ptr)
	ptr = tokenParser(ptr, token.RPAREN)

	par := tokenParser(ts, token.LPAREN)
	if len(par) != 0 {
		par = ReceiverType(par)
	}
	par = tokenParser(par, token.RPAREN)

	return append(append(ptr, par...), TypeName(ts)...)
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

func Result(ts [][]*Token) [][]*Token { return append(Parameters(ts), Type(ts)...) }

func Selector(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
}

func Signature(ts [][]*Token) [][]*Token {
	ts = Parameters(ts)
	return append(ts, Result(ts)...)
}

func Slice(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	temp := make([][]*Token, len(ts))
	if len(ts) != 0 {
		for i, t := range ts {
			if newT := Expression([][]*Token{t}); len(newT) == 0 {
				temp[i] = t
			} else {
				temp[i] = newT[0]
			}
		}
	}
	ts = tokenParser(temp, token.COLON)

	// [:]
	if len(ts) != 0 {
		for i, t := range ts {
			if newT := Expression([][]*Token{t}); len(newT) == 0 {
				temp[i] = t
			} else {
				temp[i] = newT[0]
			}
		}
	}
	a := tokenParser(temp, token.RBRACK)

	// [::]
	var b [][]*Token
	if len(ts) != 0 {
		b = Expression(ts)
	}
	b = tokenParser(b, token.COLON)
	if len(b) != 0 {
		b = Expression(b)
	}
	b = tokenParser(b, token.RBRACK)

	return append(a, b...)
}

func SliceType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	ts = tokenParser(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

// spec is wrong here - trailing semicolon is optional not mandatory
func StructType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.STRUCT)
	ts = tokenParser(ts, token.LBRACE)
	fields := FieldDecl(ts)
	next := fields
	for len(next) != 0 {
		field := tokenParser(next, token.SEMICOLON)
		field = FieldDecl(field)
		fields = append(fields, field...)
		next = field
	}
	fields = append(fields, tokenParser(fields, token.SEMICOLON)...)
	return tokenParser(append(ts, fields...), token.RBRACE)
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

func TypeAssertion(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.PERIOD)
	ts = tokenParser(ts, token.LPAREN)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	return tokenParser(ts, token.RPAREN)
}

func TypeName(ts [][]*Token) [][]*Token {
	result := QualifiedIdent(ts)
	return append(result, tokenParser(ts, token.IDENT)...)
}

func UnaryExpr(ts [][]*Token) [][]*Token {
	uo := UnaryOp(ts)
	if len(uo) != 0 {
		uo = UnaryExpr(uo)
	}
	return append(PrimaryExpr(ts), uo...)
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

func print(ts [][]*Token) {
	for _, t := range ts {
		fmt.Println(t)
	}
	fmt.Println("-----")
}
