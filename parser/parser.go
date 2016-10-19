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
	ts = append(ts, tokenParser(ts, token.MUL)...)
	return TypeName(ts)
}

// spec is wrong, maybe
func Arguments(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LPAREN)
	newTs := ExpressionList(ts)
	newTs = append(newTs, tokenParser(newTs, token.ELLIPSIS)...)
	newTs = append(newTs, tokenParser(newTs, token.COMMA)...)
	return tokenParser(append(ts, newTs...), token.RPAREN)
}

func ArrayType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	ts = tokenParser(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func Assignment(ts [][]*Token) [][]*Token {
	ts = ExpressionList(ts)
	ts = AssignOp(ts)
	return ExpressionList(ts)
}

func AssignOp(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		switch p := pop(&t); true {
		case p == nil:
		case p.tok == token.ADD_ASSIGN, p.tok == token.SUB_ASSIGN, p.tok == token.MUL_ASSIGN, p.tok == token.QUO_ASSIGN,
			p.tok == token.REM_ASSIGN, p.tok == token.AND_ASSIGN, p.tok == token.OR_ASSIGN, p.tok == token.XOR_ASSIGN,
			p.tok == token.SHL_ASSIGN, p.tok == token.SHR_ASSIGN, p.tok == token.AND_NOT_ASSIGN, p.tok == token.ASSIGN:
			result = append(result, t)
		}
	}
	return result
}

func BasicLit(ts [][]*Token) [][]*Token {
	return append(
		append(append(tokenParser(ts, token.INT), tokenParser(ts, token.FLOAT)...),
			append(tokenParser(ts, token.IMAG), tokenParser(ts, token.CHAR)...)...),
		tokenParser(ts, token.STRING)...)
}

func BinaryOp(ts [][]*Token) [][]*Token {
	return append(
		append(append(tokenParser(ts, token.LAND), tokenParser(ts, token.LOR)...),
			append(RelOp(ts), AddOp(ts)...)...),
		MulOp(ts)...)
}

func Block(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACE)
	ts = StatementList(ts)
	return tokenParser(ts, token.RBRACE)
}

func BreakStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.BREAK)
	return append(ts, tokenParser(ts, token.IDENT)...)
}

func ChannelType(ts [][]*Token) [][]*Token {
	plain := tokenParser(ts, token.CHAN)
	after := tokenParser(plain, token.ARROW)
	before := tokenParser(ts, token.ARROW)
	before = tokenParser(before, token.CHAN)
	together := append(append(plain, after...), before...)
	if len(together) == 0 {
		return nil
	}
	return Type(together)
}

func CompositeLit(ts [][]*Token) [][]*Token {
	ts = LiteralType(ts)
	return LiteralValue(ts)
}

// spec is wrong
func ConstDecl(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.CONST)
	paren := tokenParser(ts, token.LPAREN)
	paren = ConstSpec(paren)
	next := paren
	for len(next) != 0 {
		current := tokenParser(next, token.SEMICOLON)
		current = ConstSpec(current)
		paren = append(paren, current...)
		next = current
	}
	paren = append(paren, tokenParser(paren, token.SEMICOLON)...)
	paren = tokenParser(paren, token.RPAREN)
	return append(paren, ConstSpec(ts)...)
}

// spec is wrong
func ConstSpec(ts [][]*Token) [][]*Token {
	ts = IdentifierList(ts)
	ts = append(ts, Type(ts)...)
	ts = tokenParser(ts, token.ASSIGN)
	return ExpressionList(ts)
}

func ContinueStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.CONTINUE)
	return append(ts, tokenParser(ts, token.IDENT)...)
}

func Conversion(ts [][]*Token) [][]*Token {
	ts = Type(ts)
	ts = tokenParser(ts, token.LPAREN)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	ts = append(ts, tokenParser(ts, token.COMMA)...)
	return tokenParser(ts, token.RPAREN)
}

func Declaration(ts [][]*Token) [][]*Token {
	return append(append(ConstDecl(ts), TypeDecl(ts)...), VarDecl(ts)...)
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

func EmptyStmt(ts [][]*Token) [][]*Token { return ts }

func ExprCaseClause(ts [][]*Token) [][]*Token {
	ts = ExprSwitchCase(ts)
	ts = tokenParser(ts, token.COLON)
	return StatementList(ts)
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
	next := ts
	for len(next) != 0 {
		current := tokenParser(next, token.COMMA)
		current = Expression(current)
		ts = append(ts, current...)
		next = current
	}
	return ts
}

func ExpressionStmt(ts [][]*Token) [][]*Token { return Expression(ts) }

func ExprSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenParser(ts, token.CASE)
	cas = ExpressionList(cas)
	return append(tokenParser(ts, token.DEFAULT), cas...)
}

func ExprSwitchStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.SWITCH)
	stmt := SimpleStmt(ts)
	ts = append(ts, tokenParser(stmt, token.SEMICOLON)...)
	if len(ts) == 0 {
		return nil
	}
	ts = append(ts, Expression(ts)...)
	ts = tokenParser(ts, token.LBRACE)
	next := ts
	for len(next) != 0 {
		next = ExprCaseClause(next)
		ts = append(ts, next...)
	}
	return tokenParser(ts, token.RBRACE)
}

func FallthroughStmt(ts [][]*Token) [][]*Token { return tokenParser(ts, token.FALLTHROUGH) }

func FieldDecl(ts [][]*Token) [][]*Token {
	a := IdentifierList(ts)
	if len(a) != 0 {
		a = Type(a)
	}
	ts = append(AnonymousField(ts), a...)
	return append(ts, tokenParser(ts, token.STRING)...)
}

func FunctionType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.FUNC)
	return Signature(ts)
}

func GoStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.GO)
	return Expression(ts)
}

func GotoStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.GOTO)
	return tokenParser(ts, token.IDENT)
}

func IdentifierList(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.IDENT)
	more := ts
	for len(more) != 0 {
		next := tokenParser(more, token.COMMA)
		next = tokenParser(next, token.IDENT)
		ts = append(ts, next...)
		more = next
	}
	return ts
}

func IfStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.IF)
	simple := SimpleStmt(ts)
	ts = append(ts, tokenParser(simple, token.SEMICOLON)...)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	ts = Block(ts)
	els := tokenParser(ts, token.ELSE)
	if len(els) == 0 {
		return ts
	}
	els = append(IfStmt(els), Block(els)...)
	return append(ts, els...)
}

func IncDecStmt(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	return append(tokenParser(ts, token.INC), tokenParser(ts, token.DEC)...)
}

func Index(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	return tokenParser(ts, token.RBRACK)
}

// bad spec
func InterfaceType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.INTERFACE)
	ts = tokenParser(ts, token.LBRACE)
	list := MethodSpec(ts)
	next := list
	for len(next) != 0 {
		ms := tokenParser(next, token.SEMICOLON)
		ms = MethodSpec(ms)
		list = append(list, ms...)
		next = ms
	}
	list = append(list, tokenParser(list, token.SEMICOLON)...)
	ts = append(ts, list...)
	return tokenParser(ts, token.RBRACE)
}

func Key(ts [][]*Token) [][]*Token {
	return append(tokenParser(ts, token.IDENT), Expression(ts)...)
}

func KeyedElement(ts [][]*Token) [][]*Token {
	with := Key(ts)
	with = tokenParser(with, token.COLON)
	return append(Element(ts), Element(with)...)
}

func LabeledStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.IDENT)
	ts = tokenParser(ts, token.COLON)
	if len(ts) == 0 {
		return nil
	}
	return Statement(ts)
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
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func MethodExpr(ts [][]*Token) [][]*Token {
	ts = ReceiverType(ts)
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
}

func MethodSpec(ts [][]*Token) [][]*Token {
	sig := nonBlankIdent(ts)
	sig = Signature(sig)
	return append(sig, TypeName(ts)...)
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

func PackageName(ts [][]*Token) [][]*Token { return nonBlankIdent(ts) }

func ParameterDecl(ts [][]*Token) [][]*Token {
	ts = append(ts, IdentifierList(ts)...)
	ts = append(ts, tokenParser(ts, token.ELLIPSIS)...)
	if len(ts) == 0 {
		return nil
	}
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

func PointerType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.MUL)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
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

func Result(ts [][]*Token) [][]*Token {
	var t [][]*Token
	if len(ts) != 0 {
		t = Type(ts)
	}
	return append(Parameters(ts), t...)
}

func ReturnStmt(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.RETURN)
	return append(ts, ExpressionList(ts)...)
}

func Selector(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.PERIOD)
	return tokenParser(ts, token.IDENT)
}

func SendStmt(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	ts = tokenParser(ts, token.ARROW)
	return Expression(ts)
}

func ShortVarDecl(ts [][]*Token) [][]*Token {
	ts = IdentifierList(ts)
	ts = tokenParser(ts, token.DEFINE)
	return ExpressionList(ts)
}

func Signature(ts [][]*Token) [][]*Token {
	ts = Parameters(ts)
	return append(ts, Result(ts)...)
}

func SimpleStmt(ts [][]*Token) [][]*Token {
	return append(
		append(append(EmptyStmt(ts), ExpressionStmt(ts)...),
			append(SendStmt(ts), IncDecStmt(ts)...)...),
		append(Assignment(ts), ShortVarDecl(ts)...)...)
}

func Slice(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	ts = append(ts, Expression(ts)...)
	ts = tokenParser(ts, token.COLON)

	a := append(ts, Expression(ts)...)

	b := Expression(ts)
	b = tokenParser(b, token.COLON)
	b = Expression(b)

	return tokenParser(append(a, b...), token.RBRACK)
}

func SliceType(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.LBRACK)
	ts = tokenParser(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func Statement(ts [][]*Token) [][]*Token {
	// switch
	// select
	// for
	// defer

	return append(
		append(append(append(Declaration(ts), LabeledStmt(ts)...), append(SimpleStmt(ts), GoStmt(ts)...)...),
			append(append(ReturnStmt(ts), BreakStmt(ts)...), append(ContinueStmt(ts), GotoStmt(ts)...)...)...),
		append(append(FallthroughStmt(ts), Block(ts)...), IfStmt(ts)...)...)
}

func StatementList(ts [][]*Token) [][]*Token {
	if len(ts) == 0 {
		return nil
	}
	list := Statement(ts)
	next := list
	for len(next) != 0 {
		current := tokenParser(next, token.SEMICOLON)
		current = Statement(current)
		list = append(list, current...)
		next = current
	}
	list = append(list, tokenParser(list, token.SEMICOLON)...)
	return append(ts, list...)
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
	name := TypeName(ts)
	lit := TypeLit(ts)
	paren := tokenParser(ts, token.LPAREN)
	if len(paren) != 0 {
		paren = Type(paren)
	}
	paren = tokenParser(paren, token.RPAREN)
	return append(append(name, paren...), lit...)
}

func TypeAssertion(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.PERIOD)
	ts = tokenParser(ts, token.LPAREN)
	if len(ts) != 0 {
		ts = Expression(ts)
	}
	return tokenParser(ts, token.RPAREN)
}

func TypeDecl(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.TYPE)
	multi := tokenParser(ts, token.LPAREN)
	multi = TypeSpec(multi)
	next := multi
	for len(next) != 0 {
		current := tokenParser(next, token.SEMICOLON)
		current = TypeSpec(current)
		multi = append(multi, current...)
		next = current
	}
	multi = append(multi, tokenParser(multi, token.SEMICOLON)...)
	multi = tokenParser(multi, token.RPAREN)
	return append(TypeSpec(ts), multi...)
}

func TypeList(ts [][]*Token) [][]*Token {
	ts = Type(ts)
	next := ts
	for len(next) != 0 {
		current := tokenParser(next, token.COMMA)
		current = Type(current)
		ts = append(ts, current...)
		next = current
	}
	return ts
}

func TypeLit(ts [][]*Token) [][]*Token {
	return append(
		append(append(ArrayType(ts), StructType(ts)...),
			append(PointerType(ts), FunctionType(ts)...)...),
		append(append(InterfaceType(ts), SliceType(ts)...),
			append(MapType(ts), ChannelType(ts)...)...)...)
}

func TypeName(ts [][]*Token) [][]*Token {
	result := QualifiedIdent(ts)
	return append(result, tokenParser(ts, token.IDENT)...)
}

func TypeSpec(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.IDENT)
	return Type(ts)
}

func TypeSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenParser(ts, token.CASE)
	return append(TypeList(cas), tokenParser(ts, token.DEFAULT)...)
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

func VarDecl(ts [][]*Token) [][]*Token {
	ts = tokenParser(ts, token.VAR)
	paren := tokenParser(ts, token.LPAREN)
	paren = VarSpec(paren)
	next := paren
	for len(next) != 0 {
		current := tokenParser(next, token.SEMICOLON)
		current = VarSpec(current)
		paren = append(paren, current...)
		next = current
	}
	paren = append(paren, tokenParser(paren, token.SEMICOLON)...)
	paren = tokenParser(paren, token.RPAREN)
	return append(VarSpec(ts), paren...)
}

func VarSpec(ts [][]*Token) [][]*Token {
	ts = IdentifierList(ts)
	typ := Type(ts)
	extra := tokenParser(typ, token.ASSIGN)
	extra = ExpressionList(extra)
	typ = append(typ, extra...)
	assign := tokenParser(ts, token.ASSIGN)
	assign = ExpressionList(assign)
	return append(typ, assign...)
}

func nonBlankIdent(ts [][]*Token) [][]*Token {
	var result [][]*Token
	for _, t := range ts {
		if p := pop(&t); p != nil && p.tok == token.IDENT && p.lit != `_` {
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
