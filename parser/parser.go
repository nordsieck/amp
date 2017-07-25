package parser

import (
	"bytes"
	"fmt"
	"go/token"
)

var (
	_true  = true
	_false = false
)

type (
	Reader   func([][]*Token) [][]*Token
	Parser   func([][]*Token) ([]Renderer, [][]*Token)
	Renderer interface {
		Render() []byte
	}

	Stator func([]State) []State
	State  struct {
		r []Renderer
		t []*Token
	}

	e struct{} // empty renderer
)

func (_ e) Render() []byte { return nil }

func AddOp(ss []State) []State {
	var result []State
	for _, s := range ss {
		switch p := pop(&s.t); true {
		case p == nil:
		case p.tok == token.ADD, p.tok == token.SUB, p.tok == token.OR, p.tok == token.AND:
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func AnonymousField(ss []State) []State {
	ss = append(ss, tokenParserState(ss, token.MUL)...)
	ss = TypeName(ss)
	for i, s := range ss {
		af := anonymousField{typeName: s.r[len(s.r)-1]}
		prev := s.r[len(s.r)-2]
		if tok, ok := prev.(*Token); ok && tok.tok == token.MUL {
			af.pointer = true
			s.r = s.r[:len(s.r)-2]
		} else {
			s.r = s.r[:len(s.r)-1]
		}
		ss[i].r = rAppend(s.r, 0, af)
	}
	return ss
}

type anonymousField struct {
	pointer  bool
	typeName Renderer
}

func (a anonymousField) Render() []byte {
	if a.pointer {
		return append([]byte(`*`), a.typeName.Render()...)
	}
	return a.typeName.Render()
}

// bad spec
// "(" [ ExpressionList [ "..." ] [ "," ]] ")"
func Arguments(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LPAREN)
	newTs := ExpressionList(ts)
	newTs = append(newTs, tokenReader(newTs, token.ELLIPSIS)...)
	newTs = append(newTs, tokenReader(newTs, token.COMMA)...)
	return tokenReader(append(ts, newTs...), token.RPAREN)
}

func ArrayType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACK)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	ts = tokenReader(ts, token.RBRACK)
	return Type(ts)
}

func ArrayTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = ExpressionState(ss)
	ss = tokenReaderState(ss, token.RBRACK)
	ss = TypeState(ss)

	for i, s := range ss {
		at := arrayType{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, at)
	}
	return ss
}

type arrayType struct{ expr, typ Renderer }

func (a arrayType) Render() []byte {
	return append(append([]byte(`[`), a.expr.Render()...), append([]byte(`]`), a.typ.Render()...)...)
}

func Assignment(ts [][]*Token) [][]*Token {
	ts = ExpressionList(ts)
	ts = fromState(AssignOp(toState(ts)))
	return ExpressionList(ts)
}

func AssignOp(ss []State) []State {
	var result []State
	for _, s := range ss {
		switch p := pop(&s.t); true {
		case p == nil:
		case p.tok == token.ADD_ASSIGN, p.tok == token.SUB_ASSIGN, p.tok == token.MUL_ASSIGN, p.tok == token.QUO_ASSIGN,
			p.tok == token.REM_ASSIGN, p.tok == token.AND_ASSIGN, p.tok == token.OR_ASSIGN, p.tok == token.XOR_ASSIGN,
			p.tok == token.SHL_ASSIGN, p.tok == token.SHR_ASSIGN, p.tok == token.AND_NOT_ASSIGN, p.tok == token.ASSIGN:
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func BasicLit(ss []State) []State {
	return append(
		append(append(tokenParserState(ss, token.INT), tokenParserState(ss, token.FLOAT)...),
			append(tokenParserState(ss, token.IMAG), tokenParserState(ss, token.CHAR)...)...),
		tokenParserState(ss, token.STRING)...)
}

func BinaryOp(ss []State) []State {
	return append(
		append(append(AddOp(ss), RelOp(ss)...),
			append(MulOp(ss), tokenParserState(ss, token.LAND)...)...),
		tokenParserState(ss, token.LOR)...)
}

// bad spec
// "{" StatementList [ ";" ] "}"
// force empty stmt
func Block(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACE)
	ts = StatementList(ts)
	return tokenReader(ts, token.RBRACE)
}

func BreakStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.BREAK)
	return append(ts, tokenReader(ts, token.IDENT)...)
}

func ChannelType(ts [][]*Token) [][]*Token {
	plain := tokenReader(ts, token.CHAN)
	after := tokenReader(plain, token.ARROW)
	before := tokenReader(ts, token.ARROW)
	before = tokenReader(before, token.CHAN)
	together := append(append(plain, after...), before...)
	if len(together) == 0 {
		return nil
	}
	return Type(together)
}

func ChannelTypeState(ss []State) []State {
	plain := tokenReaderState(ss, token.CHAN)
	after := tokenReaderState(plain, token.ARROW)
	before := tokenReaderState(ss, token.ARROW)
	before = tokenReaderState(before, token.CHAN)

	if len(plain) != 0 {
		plain = TypeState(plain)
		for _, s := range plain {
			ct := channelType{nil, s.r[len(s.r)-1]}
			s.r[len(s.r)-1] = ct
		}
	}
	if len(after) != 0 {
		after = TypeState(after)
		for _, s := range after {
			ct := channelType{&_false, s.r[len(s.r)-1]}
			s.r[len(s.r)-1] = ct
		}
	}
	if len(before) != 0 {
		before = TypeState(before)
		for _, s := range before {
			ct := channelType{&_true, s.r[len(s.r)-1]}
			s.r[len(s.r)-1] = ct
		}
	}
	return append(append(plain, after...), before...)
}

type channelType struct {
	leadingArrow *bool // nil = no arrow
	r            Renderer
}

func (c channelType) Render() []byte {
	preamble := `chan<- `
	if c.leadingArrow == nil {
		preamble = `chan `
	} else if *c.leadingArrow {
		preamble = `<-chan `
	}

	return append([]byte(preamble), c.r.Render()...)
}

func CommCase(ts [][]*Token) [][]*Token {
	cas := tokenReader(ts, token.CASE)
	cas = append(SendStmt(cas), RecvStmt(cas)...)
	return append(cas, tokenReader(ts, token.DEFAULT)...)
}

func CommClause(ts [][]*Token) [][]*Token {
	ts = CommCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func CompositeLit(ts [][]*Token) [][]*Token {
	ts = LiteralType(ts)
	return LiteralValue(ts)
}

// bad spec
// "const" ( ConstSpec | "(" [ ConstSpec { ";" ConstSpec } [ ";" ]] ")" )
func ConstDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.CONST)
	paren := tokenReader(ts, token.LPAREN)
	paren = ConstSpec(paren)
	next := paren
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = ConstSpec(current)
		paren = append(paren, current...)
		next = current
	}
	paren = append(paren, tokenReader(paren, token.SEMICOLON)...)
	paren = tokenReader(paren, token.RPAREN)
	return append(paren, ConstSpec(ts)...)
}

// bad spec
// IdentifierLit [ Type ] "=" ExpressionList
func ConstSpec(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	ts = append(ts, Type(ts)...)
	ts = tokenReader(ts, token.ASSIGN)
	return ExpressionList(ts)
}

func ContinueStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.CONTINUE)
	return append(ts, tokenReader(ts, token.IDENT)...)
}

func Conversion(ts [][]*Token) [][]*Token {
	ts = Type(ts)
	ts = tokenReader(ts, token.LPAREN)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	ts = append(ts, tokenReader(ts, token.COMMA)...)
	return tokenReader(ts, token.RPAREN)
}

func Declaration(ts [][]*Token) [][]*Token {
	return append(append(ConstDecl(ts), TypeDecl(ts)...), VarDecl(ts)...)
}

func DeferStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.DEFER)
	return Expression(ts)
}

func Element(ts [][]*Token) [][]*Token {
	return append(Expression(ts), LiteralValue(ts)...)
}

func ElementList(ts [][]*Token) [][]*Token {
	ts = KeyedElement(ts)
	base := ts
	for len(base) != 0 {
		next := tokenReader(base, token.COMMA)
		next = KeyedElement(next)
		ts = append(ts, next...)
		base = next
	}
	return ts
}

func EllipsisArrayType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACK)
	ts = tokenReader(ts, token.ELLIPSIS)
	ts = tokenReader(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func EllipsisArrayTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	ss = tokenReaderState(ss, token.ELLIPSIS)
	ss = tokenReaderState(ss, token.RBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)

	for i, s := range ss {
		eat := ellipsisArrayType{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, eat)
	}
	return ss
}

type ellipsisArrayType struct{ r Renderer }

func (e ellipsisArrayType) Render() []byte { return append([]byte(`[...]`), e.r.Render()...) }

func EmptyStmt(ts [][]*Token) [][]*Token { return ts }

func ExprCaseClause(ts [][]*Token) [][]*Token {
	ts = ExprSwitchCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func Expression(ts [][]*Token) [][]*Token {
	base := UnaryExpr(ts)
	comp := fromState(BinaryOp(toState(base)))
	if len(comp) == 0 {
		return base
	}
	comp = Expression(comp)
	return append(base, comp...)
}

func ExpressionState(ss []State) []State {
	ss = UnaryExprState(ss)
	// Expression binary_op Expression
	return ss
}

func ExpressionList(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	next := ts
	for len(next) != 0 {
		current := tokenReader(next, token.COMMA)
		current = Expression(current)
		ts = append(ts, current...)
		next = current
	}
	return ts
}

func ExpressionStmt(ts [][]*Token) [][]*Token { return Expression(ts) }

func ExprSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenReader(ts, token.CASE)
	cas = ExpressionList(cas)
	return append(tokenReader(ts, token.DEFAULT), cas...)
}

// bad spec
// "switch" [ SimpleStmt ";" ] [ Expression ] "{" [ ExprCaseClause { ";" ExprCaseClause } [ ";" ]] "}"
func ExprSwitchStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.SWITCH)
	stmt := SimpleStmt(ts)
	ts = append(ts, tokenReader(stmt, token.SEMICOLON)...)
	if len(ts) == 0 {
		return nil
	}
	ts = append(ts, Expression(ts)...)
	ts = tokenReader(ts, token.LBRACE)
	list := ExprCaseClause(ts)
	next := list
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = ExprCaseClause(current)
		list = append(list, current...)
		next = current
	}
	list = append(list, tokenReader(list, token.SEMICOLON)...)
	ts = append(ts, list...)
	return tokenReader(ts, token.RBRACE)
}

func FallthroughStmt(ss []State) []State { return tokenParserState(ss, token.FALLTHROUGH) }

func FieldDecl(ts [][]*Token) [][]*Token {
	a := fromState(IdentifierList(toState(ts)))
	if len(a) != 0 {
		a = Type(a)
	}
	ts = append(fromState(AnonymousField(toState(ts))), a...)
	return append(ts, tokenReader(ts, token.STRING)...)
}

// The problem is that the renderer slice gets copied
func FieldDeclState(ss []State) []State {
	typ := IdentifierList(ss)
	if len(typ) != 0 {
		typ = TypeState(typ)
	}
	typTag := tokenParserState(typ, token.STRING)

	af := AnonymousField(ss)
	afTag := tokenParserState(af, token.STRING)

	for i, t := range typ {
		fd := fieldDecl{idList: t.r[len(t.r)-2], typ: t.r[len(t.r)-1]}
		typ[i].r = rAppend(t.r, 2, fd)
	}
	for i, t := range typTag {
		fd := fieldDecl{idList: t.r[len(t.r)-3], typ: t.r[len(t.r)-2], tag: t.r[len(t.r)-1]}
		typTag[i].r = rAppend(t.r, 3, fd)
	}
	for i, a := range af {
		fd := fieldDecl{anonField: a.r[len(a.r)-1]}
		af[i].r = rAppend(a.r, 1, fd)
	}
	for i, a := range afTag {
		fd := fieldDecl{anonField: a.r[len(a.r)-2], tag: a.r[len(a.r)-1]}
		afTag[i].r = rAppend(a.r, 2, fd)
	}

	return append(append(typ, typTag...), append(af, afTag...)...)
}

type fieldDecl struct {
	idList, typ, anonField, tag Renderer
}

func (f fieldDecl) Render() []byte {
	var ret []byte
	if f.anonField != nil {
		ret = f.anonField.Render()
	} else {
		ret = append(append(f.idList.Render(), ` `...), f.typ.Render()...)
	}
	if f.tag != nil {
		ret = append(append(ret, ` `...), f.tag.Render()...)
	}
	return ret
}

func ForClause(ts [][]*Token) [][]*Token {
	ts = append(ts, SimpleStmt(ts)...)
	ts = tokenReader(ts, token.SEMICOLON)
	ts = append(ts, Expression(ts)...)
	ts = tokenReader(ts, token.SEMICOLON)
	return append(ts, SimpleStmt(ts)...)
}

func ForStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.FOR)
	ts = append(append(ts, Expression(ts)...), append(ForClause(ts), RangeClause(ts)...)...)
	return Block(ts)
}

func Function(ts [][]*Token) [][]*Token {
	ts = Signature(ts)
	return Block(ts)
}

// bad spec
// "func" FunctionName Function
func FunctionDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.FUNC)
	ts = tokenReader(ts, token.IDENT)
	return Function(ts)
}

func FunctionLit(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.FUNC)
	return Function(ts)
}

func FunctionType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.FUNC)
	return Signature(ts)
}

func FunctionTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.FUNC)
	ss = SignatureState(ss)

	for i, s := range ss {
		ft := functionType{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, ft)
	}
	return ss
}

type functionType struct{ r Renderer }

func (f functionType) Render() []byte { return append([]byte(`func`), f.r.Render()...) }

func GoStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.GO)
	return Expression(ts)
}

func GotoStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.GOTO)
	return tokenReader(ts, token.IDENT)
}

func IdentifierList(ss []State) []State {
	ss = tokenParserState(ss, token.IDENT)

	loop := ss
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = tokenParserState(loop, token.IDENT)
		ss = append(ss, loop...)
	}

	for i, s := range ss {
		idList := identifierList{s.r[len(s.r)-1]}
		s.r = s.r[:len(s.r)-1]

		for len(s.r) > 1 {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
				idList = append(idList, s.r[len(s.r)-2])
				s.r = s.r[:len(s.r)-2]
			} else {
				break
			}
		}

		// reverse
		for i := 0; i < len(idList)/2; i++ {
			j := len(idList) - 1 - i
			idList[i], idList[j] = idList[j], idList[i]
		}

		ss[i].r = rAppend(s.r, 0, idList)
	}
	return ss
}

type identifierList []Renderer

func (il identifierList) Render() []byte {
	var result [][]byte
	for i := 0; i < len(il); i++ {
		result = append(result, il[i].Render())
	}
	return bytes.Join(result, []byte(`,`))
}

func IfStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IF)
	simple := SimpleStmt(ts)
	ts = append(ts, tokenReader(simple, token.SEMICOLON)...)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	ts = Block(ts)
	els := tokenReader(ts, token.ELSE)
	els = append(IfStmt(els), Block(els)...)
	return append(ts, els...)
}

// bad spec
// "import" ( ImportSpec | "(" [ ImportSpec { ";" ImportSpec } [ ";" ]] ")" )
func ImportDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IMPORT)
	group := tokenReader(ts, token.LPAREN)
	elements := ImportSpec(group)
	next := elements
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = ImportSpec(current)
		elements = append(elements, current...)
		next = current
	}
	elements = append(elements, tokenReader(elements, token.SEMICOLON)...)
	group = append(group, elements...)
	group = tokenReader(group, token.RPAREN)
	return append(group, ImportSpec(ts)...)
}

// bad spec
// [ "." | "_" | PackageName ] ImportPath
func ImportSpec(ts [][]*Token) [][]*Token {
	names := append(tokenReader(ts, token.PERIOD), tokenReader(ts, token.IDENT)...)
	ts = append(ts, names...)
	return tokenReader(ts, token.STRING)
}

func IncDecStmt(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	return append(tokenReader(ts, token.INC), tokenReader(ts, token.DEC)...)
}

func Index(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACK)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	return tokenReader(ts, token.RBRACK)
}

// bad spec
// "interface" "{" [ MethodSpec { ";" MethodSpec } [ ";" ]] "}"
func InterfaceType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.INTERFACE)
	ts = tokenReader(ts, token.LBRACE)
	list := MethodSpec(ts)
	next := list
	for len(next) != 0 {
		ms := tokenReader(next, token.SEMICOLON)
		ms = MethodSpec(ms)
		list = append(list, ms...)
		next = ms
	}
	list = append(list, tokenReader(list, token.SEMICOLON)...)
	ts = append(ts, list...)
	return tokenReader(ts, token.RBRACE)
}

func InterfaceTypeState(ss []State) []State {
	ss = tokenParserState(ss, token.INTERFACE)
	ss = tokenReaderState(ss, token.LBRACE)
	methods := MethodSpecState(ss)
	loop := methods
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = MethodSpecState(loop)
		methods = append(methods, loop...)
	}
	methods = append(methods, tokenReaderState(methods, token.SEMICOLON)...)
	ss = append(ss, methods...)
	ss = tokenReaderState(ss, token.RBRACE)

	for i, s := range ss {
		it := interfaceType{}
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.INTERFACE {
				break
			}
			it = append(it, s.r[len(s.r)-1])
			s.r = s.r[:len(s.r)-1]
		}

		// reverse
		for j := 0; j < len(it); j++ {
			it[j], it[len(it)-1-i] = it[len(it)-1-i], it[j]
		}

		ss[i].r = rAppend(s.r, 1, it)
	}
	return ss
}

type interfaceType []Renderer

func (it interfaceType) Render() []byte {
	ret := []byte(`interface{`)
	for _, i := range it {
		ret = append(append(ret, i.Render()...), `;`...)
	}
	return append(ret, `}`...)
}

// TODO: missing LiteralValue
func Key(ts [][]*Token) [][]*Token {
	return append(tokenReader(ts, token.IDENT), Expression(ts)...)
}

func KeyState(ss []State) []State {
	fn := tokenParserState(ss, token.IDENT)
	// expr := Expression(ss)
	// lv := LiteralValue(ss)

	for i, f := range fn {
		k := key{f.r[len(f.r)-1]}
		fn[i].r = rAppend(f.r, 1, k)
	}

	ss = fn
	return ss
}

type key struct{ r Renderer }

func (k key) Render() []byte { return k.r.Render() }

func KeyedElement(ts [][]*Token) [][]*Token {
	with := Key(ts)
	with = tokenReader(with, token.COLON)
	return append(Element(ts), Element(with)...)
}

func LabeledStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IDENT)
	ts = tokenReader(ts, token.COLON)
	if len(ts) == 0 {
		return nil
	}
	return Statement(ts)
}

func Literal(ts [][]*Token) [][]*Token {
	basicLit := fromState(BasicLit(toState(ts)))
	return append(append(basicLit, CompositeLit(ts)...), FunctionLit(ts)...)
}

func LiteralState(ss []State) []State {
	return BasicLit(ss)
	// CompositeLit
	// FunctionLit
}

func LiteralType(ts [][]*Token) [][]*Token {
	tn := fromState(TypeName(toState(ts)))
	return append(
		append(
			append(StructType(ts), ArrayType(ts)...),
			append(EllipsisArrayType(ts), SliceType(ts)...)...),
		append(MapType(ts), tn...)...)
}

func LiteralTypeState(ss []State) []State {
	return append(
		append(append(StructTypeState(ss), EllipsisArrayTypeState(ss)...),
			append(SliceTypeState(ss), MapTypeState(ss)...)...),
		append(ArrayTypeState(ss), TypeName(ss)...)...)
}

func LiteralValue(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACE)
	if len(ts) == 0 {
		return nil
	}
	list := ElementList(ts)
	list = append(list, tokenReader(list, token.COMMA)...)
	ts = append(ts, list...)
	return tokenReader(ts, token.RBRACE)
}

func MapType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.MAP)
	ts = tokenReader(ts, token.LBRACK)
	if len(ts) == 0 {
		return nil
	}
	ts = Type(ts)
	ts = tokenReader(ts, token.RBRACK)
	return Type(ts)
}

func MapTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.MAP)
	ss = tokenReaderState(ss, token.LBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)
	ss = tokenReaderState(ss, token.RBRACK)
	ss = TypeState(ss)
	for i, s := range ss {
		mt := mapType{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, mt)
	}
	return ss
}

type mapType struct{ k, v Renderer }

func (m mapType) Render() []byte {
	return append(append(append([]byte(`map[`), m.k.Render()...), `]`...), m.v.Render()...)
}

// bad spec
// "func" Reciever MethodName Function
func MethodDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.FUNC)
	ts = Parameters(ts)
	ts = fromState(nonBlankIdent(toState(ts)))
	return Function(ts)
}

func MethodExpr(ts [][]*Token) [][]*Token {
	ts = ReceiverType(ts)
	ts = tokenReader(ts, token.PERIOD)
	return tokenReader(ts, token.IDENT)
}

func MethodSpec(ts [][]*Token) [][]*Token {
	sig := fromState(nonBlankIdent(toState(ts)))
	sig = Signature(sig)
	ts = fromState(TypeName(toState(ts)))
	return append(sig, ts...)
}

func MethodSpecState(ss []State) []State {
	sig := nonBlankIdent(ss)
	sig = SignatureState(sig)
	for i, s := range sig {
		ms := methodSpec{name: s.r[len(s.r)-2], signature: s.r[len(s.r)-1]}
		sig[i].r = rAppend(s.r, 2, ms)
	}

	tn := TypeName(ss)
	for i, t := range tn {
		ms := methodSpec{iTypeName: t.r[len(t.r)-1]}
		tn[i].r = rAppend(t.r, 1, ms)
	}
	return append(sig, tn...)
}

type methodSpec struct{ name, signature, iTypeName Renderer }

func (m methodSpec) Render() []byte {
	if m.iTypeName == nil {
		return append(m.name.Render(), m.signature.Render()...)
	}
	return m.iTypeName.Render()
}

func MulOp(ss []State) []State {
	var result []State
	for _, s := range ss {
		switch p := pop(&s.t); true {
		case p == nil:
		case p.tok == token.MUL, p.tok == token.QUO, p.tok == token.REM, p.tok == token.SHL,
			p.tok == token.SHR, p.tok == token.AND, p.tok == token.AND_NOT:
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func Operand(ts [][]*Token) [][]*Token {
	xp := tokenReader(ts, token.LPAREN)
	if len(xp) != 0 {
		xp = Expression(xp)
	}
	xp = tokenReader(xp, token.RPAREN)
	return append(append(Literal(ts), OperandName(ts)...), append(MethodExpr(ts), xp...)...)
}

func OperandState(ss []State) []State {
	ss = LiteralState(ss)
	// operand name
	// methodexpr
	// "(" expression ")"
	return ss
}

func OperandName(ts [][]*Token) [][]*Token {
	qi := fromState(QualifiedIdent(toState(ts)))
	return append(tokenReader(ts, token.IDENT), qi...)
}

func PackageClause(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.PACKAGE)
	return fromState(PackageName(toState(ts)))
}

func PackageName(ss []State) []State { return nonBlankIdent(ss) }

func ParameterDecl(ts [][]*Token) [][]*Token {
	idList := fromState(IdentifierList(toState(ts)))
	ts = append(ts, idList...)
	ts = append(ts, tokenReader(ts, token.ELLIPSIS)...)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func ParameterDeclIDListState(ss []State) []State {
	ss = IdentifierList(ss)
	ss = append(ss, tokenParserState(ss, token.ELLIPSIS)...)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)

	for i, s := range ss {
		pd := parameterDecl{typ: s.r[len(s.r)-1]}
		s.r = s.r[:len(s.r)-1]

		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.ELLIPSIS {
			pd.ellipsis = true
			s.r = s.r[:len(s.r)-1]
		}
		pd.idList = s.r[len(s.r)-1]
		ss[i].r = rAppend(s.r, 1, pd)
	}
	return ss
}

func ParameterDeclNoListState(ss []State) []State {
	ss = append(ss, tokenParserState(ss, token.ELLIPSIS)...)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)

	for i, s := range ss {
		pd := parameterDecl{typ: s.r[len(s.r)-1]}
		s.r = s.r[:len(s.r)-1]

		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.ELLIPSIS {
			pd.ellipsis = true
			s.r = s.r[:len(s.r)-1]
		}
		ss[i].r = rAppend(s.r, 0, pd)
	}
	return ss
}

type parameterDecl struct {
	idList   Renderer
	ellipsis bool
	typ      Renderer
}

func (pd parameterDecl) Render() []byte {
	var ret []byte
	if pd.idList != nil {
		ret = append(ret, pd.idList.Render()...)
		ret = append(ret, ` `...)
	}
	if pd.ellipsis {
		ret = append(ret, `... `...)
	}
	return append(ret, pd.typ.Render()...)
}

func ParameterList(ts [][]*Token) [][]*Token {
	ts = ParameterDecl(ts)
	next := ts
	for len(next) != 0 {
		decl := tokenReader(next, token.COMMA)
		decl = ParameterDecl(decl)
		ts = append(ts, decl...)
		next = decl
	}
	return ts
}

func ParameterListState(ss []State) []State {
	list := ParameterDeclIDListState(ss)
	loop := list
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = ParameterDeclIDListState(loop)
		list = append(list, loop...)
	}

	nolist := ParameterDeclNoListState(ss)
	loop = nolist
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = ParameterDeclNoListState(loop)
		nolist = append(nolist, loop...)
	}
	ss = append(nolist, list...)

	for i, s := range ss {
		pl := parameterList{s.r[len(s.r)-1]}
		s.r = s.r[:len(s.r)-1]

		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
				pl = append(pl, s.r[len(s.r)-2])
				s.r = s.r[:len(s.r)-2]
			} else {
				break
			}
		}

		// reverse
		for i := 0; i < len(pl)/2; i++ {
			pl[i], pl[len(pl)-1-i] = pl[len(pl)-1-i], pl[i]
		}

		ss[i].r = rAppend(s.r, 0, pl)
	}
	return ss
}

type parameterList []Renderer

func (pl parameterList) Render() []byte {
	var ret [][]byte
	for _, r := range pl {
		ret = append(ret, r.Render())
	}
	return bytes.Join(ret, []byte(`,`))
}

func Parameters(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LPAREN)
	params := ParameterList(ts)
	params = append(params, tokenReader(params, token.COMMA)...)
	ts = append(ts, params...)
	return tokenReader(ts, token.RPAREN)
}

func ParametersState(ss []State) []State {
	ss = tokenParserState(ss, token.LPAREN)
	params := ParameterListState(ss)
	params = append(params, tokenParserState(params, token.COMMA)...)
	ss = tokenReaderState(append(ss, params...), token.RPAREN)

	for i, s := range ss {
		p := parameters{}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
			p.comma = true
			s.r = s.r[:len(s.r)-1]
		}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.LPAREN {
			ss[i].r = rAppend(s.r, 1, p)
		} else {
			p.r = s.r[len(s.r)-1]
			ss[i].r = rAppend(s.r, 2, p)
		}
	}
	return ss
}

type parameters struct {
	r     Renderer
	comma bool
}

func (p parameters) Render() []byte {
	ret := []byte(`(`)
	if p.r != nil {
		ret = append(ret, p.r.Render()...)
		if p.comma {
			ret = append(ret, `,`...)
		}
	}
	return append(ret, `)`...)
}

func PointerType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.MUL)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

// TODO: figure out how to treat pointers.  Are they part of the thing they are
// attached to, or are they a universal wrapper?
func PointerTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.MUL)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)
	for _, s := range ss {
		pt := pointerType{s.r[len(s.r)-1]}
		s.r[len(s.r)-1] = pt
	}
	return ss
}

type pointerType struct{ r Renderer }

func (p pointerType) Render() []byte { return append([]byte(`*`), p.r.Render()...) }

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

func PrimaryExprState(ss []State) []State {
	ss = OperandState(ss)
	// conversion
	// selector
	// index
	// slice
	// type assertion
	// arguments
	return ss
}

func QualifiedIdent(ss []State) []State {
	ss = PackageName(ss)
	ss = tokenParserState(ss, token.PERIOD)
	ss = tokenParserState(ss, token.IDENT)
	for i, s := range ss {
		qi := qualifiedIdent{s.r[len(s.r)-3], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 3, qi)
	}
	return ss
}

type qualifiedIdent struct{ pkg, name Renderer }

func (q qualifiedIdent) Render() []byte {
	var ret []byte
	ret = append(ret, q.pkg.Render()...)
	ret = append(ret, `.`...)
	return append(ret, q.name.Render()...)
}

func RangeClause(ts [][]*Token) [][]*Token {
	exp := ExpressionList(ts)
	exp = tokenReader(exp, token.ASSIGN)
	id := fromState(IdentifierList(toState(ts)))
	id = tokenReader(id, token.DEFINE)
	ts = append(ts, append(exp, id...)...)
	ts = tokenReader(ts, token.RANGE)
	return Expression(ts)
}

func ReceiverType(ts [][]*Token) [][]*Token {
	ptr := tokenReader(ts, token.LPAREN)
	ptr = tokenReader(ptr, token.MUL)
	ptr = fromState(TypeName(toState(ptr)))
	ptr = tokenReader(ptr, token.RPAREN)

	par := tokenReader(ts, token.LPAREN)
	if len(par) != 0 {
		par = ReceiverType(par)
	}
	par = tokenReader(par, token.RPAREN)
	ts = fromState(TypeName(toState(ts)))

	return append(append(ptr, par...), ts...)
}

func RecvStmt(ts [][]*Token) [][]*Token {
	expr := ExpressionList(ts)
	expr = tokenReader(expr, token.ASSIGN)
	ident := fromState(IdentifierList(toState(ts)))
	ident = tokenReader(ident, token.DEFINE)
	return Expression(append(ts, append(expr, ident...)...))
}

func RelOp(ss []State) []State {
	var result []State
	for _, s := range ss {
		switch p := pop(&s.t); true {
		case p == nil:
		case p.tok == token.EQL, p.tok == token.NEQ, p.tok == token.LSS,
			p.tok == token.LEQ, p.tok == token.GTR, p.tok == token.GEQ:
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func Result(ts [][]*Token) [][]*Token {
	if len(ts) == 0 {
		return nil
	}
	return append(Parameters(ts), Type(ts)...)
}

func ResultState(ss []State) []State {
	if len(ss) == 0 {
		return nil
	}

	pp := ParametersState(ss)
	for i, p := range pp {
		r := result{parameters: p.r[len(p.r)-1]}
		pp[i].r = rAppend(p.r, 1, r)
	}

	tt := TypeState(ss)
	for i, t := range tt {
		r := result{typ: t.r[len(t.r)-1]}
		tt[i].r = rAppend(t.r, 1, r)
	}
	return append(pp, tt...)
}

type result struct{ parameters, typ Renderer }

func (r result) Render() []byte {
	if r.parameters == nil {
		return r.typ.Render()
	}
	return r.parameters.Render()
}

func ReturnStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.RETURN)
	return append(ts, ExpressionList(ts)...)
}

func Selector(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.PERIOD)
	return tokenReader(ts, token.IDENT)
}

// bad spec
// "select" "{" [ CommClause { ";" CommClause } [ ";" ]] "}"
func SelectStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.SELECT)
	ts = tokenReader(ts, token.LBRACE)
	clause := CommClause(ts)
	next := clause
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = CommClause(current)
		clause = append(clause, current...)
		next = current
	}
	clause = append(clause, tokenReader(clause, token.SEMICOLON)...)
	ts = append(ts, clause...)
	return tokenReader(ts, token.RBRACE)
}

func SendStmt(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	ts = tokenReader(ts, token.ARROW)
	return Expression(ts)
}

func ShortVarDecl(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	ts = tokenReader(ts, token.DEFINE)
	return ExpressionList(ts)
}

func Signature(ts [][]*Token) [][]*Token {
	ts = Parameters(ts)
	return append(ts, Result(ts)...)
}

func SignatureState(ss []State) []State {
	pp := ParametersState(ss)
	pr := ResultState(pp)

	for i, p := range pp {
		s := signature{parameters: p.r[len(p.r)-1]}
		pp[i].r = rAppend(p.r, 1, s)
	}
	for i, p := range pr {
		s := signature{parameters: p.r[len(p.r)-2], result: p.r[len(p.r)-1]}
		pr[i].r = rAppend(p.r, 2, s)
	}
	return append(pp, pr...)
}

type signature struct{ parameters, result Renderer }

func (s signature) Render() []byte {
	ret := s.parameters.Render()
	if s.result != nil {
		ret = append(ret, s.result.Render()...)
	}
	return ret
}

func SimpleStmt(ts [][]*Token) [][]*Token {
	return append(
		append(append(EmptyStmt(ts), ExpressionStmt(ts)...),
			append(SendStmt(ts), IncDecStmt(ts)...)...),
		append(Assignment(ts), ShortVarDecl(ts)...)...)
}

func Slice(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACK)
	ts = append(ts, Expression(ts)...)
	ts = tokenReader(ts, token.COLON)

	a := append(ts, Expression(ts)...)

	b := Expression(ts)
	b = tokenReader(b, token.COLON)
	b = Expression(b)

	return tokenReader(append(a, b...), token.RBRACK)
}

func SliceType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LBRACK)
	ts = tokenReader(ts, token.RBRACK)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
}

func SliceTypeState(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	ss = tokenReaderState(ss, token.RBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = TypeState(ss)
	for _, s := range ss {
		slice := sliceType{s.r[len(s.r)-1]}
		s.r[len(s.r)-1] = slice
	}
	return ss
}

type sliceType struct{ r Renderer }

func (s sliceType) Render() []byte { return append([]byte(`[]`), s.r.Render()...) }

func SourceFile(ts [][]*Token) [][]*Token {
	ts = PackageClause(ts)
	ts = tokenReader(ts, token.SEMICOLON)
	next := ts
	for len(next) != 0 {
		next = ImportDecl(next)
		next = tokenReader(next, token.SEMICOLON)
		ts = append(ts, next...)
	}
	next = ts
	for len(next) != 0 {
		next = TopLevelDecl(next)
		next = tokenReader(next, token.SEMICOLON)
		ts = append(ts, next...)
	}
	return ts
}

func Statement(ts [][]*Token) [][]*Token {
	fallthroughStmt := fromState(FallthroughStmt(toState(ts)))
	return append(
		append(append(append(Declaration(ts), LabeledStmt(ts)...), append(SimpleStmt(ts), GoStmt(ts)...)...),
			append(append(ReturnStmt(ts), BreakStmt(ts)...), append(ContinueStmt(ts), GotoStmt(ts)...)...)...),
		append(append(append(fallthroughStmt, Block(ts)...), append(IfStmt(ts), SwitchStmt(ts)...)...),
			append(append(SelectStmt(ts), ForStmt(ts)...), DeferStmt(ts)...)...)...)
}

// bad spec
// Statement { ";" Statement }
// force empty stmt
func StatementList(ts [][]*Token) [][]*Token {
	if len(ts) == 0 {
		return nil
	}
	ts = Statement(ts)
	next := ts
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = Statement(current)
		ts = append(ts, current...)
		next = current
	}
	return ts
}

// bad spec
// "struct" "{" [ FieldDecl { ";" FieldDecl } [ ";" ]] "}"
func StructType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.STRUCT)
	ts = tokenReader(ts, token.LBRACE)
	fields := FieldDecl(ts)
	next := fields
	for len(next) != 0 {
		field := tokenReader(next, token.SEMICOLON)
		field = FieldDecl(field)
		fields = append(fields, field...)
		next = field
	}
	fields = append(fields, tokenReader(fields, token.SEMICOLON)...)
	return tokenReader(append(ts, fields...), token.RBRACE)
}

func StructTypeState(ss []State) []State {
	ss = tokenParserState(ss, token.STRUCT)
	ss = tokenReaderState(ss, token.LBRACE)
	fields := FieldDeclState(ss)
	next := fields
	for len(next) != 0 {
		field := tokenReaderState(next, token.SEMICOLON)
		field = FieldDeclState(field)
		fields = append(fields, field...)
		next = field
	}
	fields = append(fields, tokenReaderState(fields, token.SEMICOLON)...)
	ss = tokenReaderState(append(ss, fields...), token.RBRACE)

	for i, s := range ss {
		st := structType{}
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); !ok || tok.tok != token.STRUCT {
				st = append(st, s.r[len(s.r)-1])
				s.r = s.r[:len(s.r)-1]
			} else {
				s.r = s.r[:len(s.r)-1]
				break
			}
		}

		// reverse
		for i := 0; i < len(st)/2; i++ {
			st[i], st[len(st)-1-i] = st[len(st)-1-i], st[i]
		}

		ss[i].r = rAppend(s.r, 0, st)
	}
	return ss
}

type structType []Renderer

func (st structType) Render() []byte {
	result := []byte(`struct{`)
	for i := 0; i < len(st); i++ {
		result = append(result, st[i].Render()...)
		result = append(result, `;`...)
	}
	return append(result, `}`...)
}

func SwitchStmt(ts [][]*Token) [][]*Token { return append(ExprSwitchStmt(ts), TypeSwitchStmt(ts)...) }

func TopLevelDecl(ts [][]*Token) [][]*Token {
	return append(append(Declaration(ts), FunctionDecl(ts)...), MethodDecl(ts)...)
}

func TypeState(ss []State) []State {
	paren := tokenParserState(ss, token.LPAREN) // TODO: wrap this in a struct to preserve the parens
	if len(paren) != 0 {
		paren = TypeState(paren)
	}
	paren = tokenParserState(paren, token.RPAREN)
	ss = append(append(TypeName(ss), TypeLitState(ss)...), paren...)

	for i, s := range ss {
		var p int
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.RPAREN {
				p += 1
				s.r = s.r[:len(s.r)-1]
			} else {
				break
			}
		}
		t := typ{s.r[len(s.r)-1], p}
		for pPrime := 0; pPrime < p; pPrime++ {
			if tok, ok := s.r[len(s.r)-1].(*Token); !ok || tok.tok != token.LPAREN {
				continue
			}
		}
		ss[i].r = rAppend(s.r, p+1, t)
	}
	return ss
}

func Type(ts [][]*Token) [][]*Token {
	paren := tokenReader(ts, token.LPAREN)
	if len(paren) != 0 {
		paren = Type(paren)
	}
	paren = tokenReader(paren, token.RPAREN)
	tn := fromState(TypeName(toState(ts)))
	return append(append(tn, TypeLit(ts)...), paren...)
}

type typ struct {
	r      Renderer
	parens int
}

func (t typ) Render() []byte {
	var ret []byte

	for i := 0; i < t.parens; i++ {
		ret = append(ret, `(`...)
	}
	ret = append(ret, t.r.Render()...)
	for i := 0; i < t.parens; i++ {
		ret = append(ret, `)`...)
	}
	return ret
}

func TypeAssertion(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.PERIOD)
	ts = tokenReader(ts, token.LPAREN)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	return tokenReader(ts, token.RPAREN)
}

func TypeCaseClause(ts [][]*Token) [][]*Token {
	ts = TypeSwitchCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func TypeDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.TYPE)
	multi := tokenReader(ts, token.LPAREN)
	multi = TypeSpec(multi)
	next := multi
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = TypeSpec(current)
		multi = append(multi, current...)
		next = current
	}
	multi = append(multi, tokenReader(multi, token.SEMICOLON)...)
	multi = tokenReader(multi, token.RPAREN)
	return append(TypeSpec(ts), multi...)
}

func TypeList(ts [][]*Token) [][]*Token {
	ts = Type(ts)
	next := ts
	for len(next) != 0 {
		current := tokenReader(next, token.COMMA)
		current = Type(current)
		ts = append(ts, current...)
		next = current
	}
	return ts
}

func TypeLitState(ss []State) []State {
	return append(
		append(append(ArrayTypeState(ss), StructTypeState(ss)...),
			append(PointerTypeState(ss), FunctionTypeState(ss)...)...),
		append(append(InterfaceTypeState(ss), SliceTypeState(ss)...),
			append(MapTypeState(ss), ChannelTypeState(ss)...)...)...)
}

func TypeLit(ts [][]*Token) [][]*Token {
	return append(
		append(append(ArrayType(ts), StructType(ts)...),
			append(PointerType(ts), FunctionType(ts)...)...),
		append(append(InterfaceType(ts), SliceType(ts)...),
			append(MapType(ts), ChannelType(ts)...)...)...)
}

func TypeName(ss []State) []State {
	return append(QualifiedIdent(ss), tokenParserState(ss, token.IDENT)...)
}

func TypeSpec(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IDENT)
	return Type(ts)
}

func TypeSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenReader(ts, token.CASE)
	return append(TypeList(cas), tokenReader(ts, token.DEFAULT)...)
}

func TypeSwitchGuard(ts [][]*Token) [][]*Token {
	id := tokenReader(ts, token.IDENT)
	ts = append(ts, tokenReader(id, token.DEFINE)...)
	ts = PrimaryExpr(ts)
	ts = tokenReader(ts, token.PERIOD)
	ts = tokenReader(ts, token.LPAREN)
	ts = tokenReader(ts, token.TYPE)
	return tokenReader(ts, token.RPAREN)
}

// bad spec
// "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" [ TypeCaseClause { ";" TypeCaseClause } [ ";" ]] "}"
func TypeSwitchStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.SWITCH)
	stmt := SimpleStmt(ts)
	ts = append(ts, tokenReader(stmt, token.SEMICOLON)...)
	ts = TypeSwitchGuard(ts)
	ts = tokenReader(ts, token.LBRACE)
	clauses := TypeCaseClause(ts)
	next := clauses
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = TypeCaseClause(current)
		clauses = append(clauses, current...)
		next = current
	}
	clauses = append(clauses, tokenReader(clauses, token.SEMICOLON)...)
	ts = append(ts, clauses...)
	return tokenReader(ts, token.RBRACE)
}

func UnaryExpr(ts [][]*Token) [][]*Token {
	uo := fromState(UnaryOp(toState(ts)))
	if len(uo) == 0 {
		return PrimaryExpr(ts)
	}
	return append(PrimaryExpr(ts), UnaryExpr(uo)...)
}

func UnaryExprState(ss []State) []State {
	ss = PrimaryExprState(ss)
	// UnaryOp UnaryExpr
	return ss
}

func UnaryOp(ss []State) []State {
	var result []State
	for _, s := range ss {
		switch p := pop(&s.t); true {
		case p == nil:
		case p.tok == token.ADD, p.tok == token.SUB, p.tok == token.NOT, p.tok == token.XOR,
			p.tok == token.MUL, p.tok == token.AND, p.tok == token.ARROW:
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func VarDecl(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.VAR)
	paren := tokenReader(ts, token.LPAREN)
	paren = VarSpec(paren)
	next := paren
	for len(next) != 0 {
		current := tokenReader(next, token.SEMICOLON)
		current = VarSpec(current)
		paren = append(paren, current...)
		next = current
	}
	paren = append(paren, tokenReader(paren, token.SEMICOLON)...)
	paren = tokenReader(paren, token.RPAREN)
	return append(VarSpec(ts), paren...)
}

func VarSpec(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	typ := Type(ts)
	extra := tokenReader(typ, token.ASSIGN)
	extra = ExpressionList(extra)
	typ = append(typ, extra...)
	assign := tokenReader(ts, token.ASSIGN)
	assign = ExpressionList(assign)
	return append(typ, assign...)
}

func nonBlankIdent(ss []State) []State {
	var result []State
	for _, s := range ss {
		if p := pop(&s.t); p != nil && p.tok == token.IDENT && p.lit != `_` {
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func tokenReader(ts [][]*Token, tok token.Token) [][]*Token {
	var result [][]*Token
	var p *Token
	for _, t := range ts {
		if p = pop(&t); p != nil && p.tok == tok {
			result = append(result, t)
		}
	}
	return result
}

func tokenParser(ts [][]*Token, tok token.Token) ([]Renderer, [][]*Token) {
	var result [][]*Token
	var tree []Renderer
	var p *Token
	for _, t := range ts {
		if p = pop(&t); p != nil && p.tok == tok {
			result = append(result, t)
			tree = append(tree, p)
		}
	}
	return tree, result
}

func tokenReaderState(ss []State, tok token.Token) []State {
	var result []State
	for _, s := range ss {
		if p := pop(&s.t); p != nil && p.tok == tok {
			result = append(result, s)
		}
	}
	return result
}

func tokenParserState(ss []State, tok token.Token) []State {
	var result []State
	for _, s := range ss {
		if p := pop(&s.t); p != nil && p.tok == tok {
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func print(ts [][]*Token) {
	for _, t := range ts {
		fmt.Println(t)
	}
	fmt.Println(`-----`)
}

func empties(n int) []Renderer {
	result := make([]Renderer, 0, n)
	for i := 0; i < n; i++ {
		result = append(result, e{})
	}
	return result
}

func toState(ts [][]*Token) []State {
	var s []State
	for _, t := range ts {
		s = append(s, State{empties(len(t)), t})
	}
	return s
}

func fromState(ss []State) [][]*Token {
	var t [][]*Token
	for _, s := range ss {
		t = append(t, s.t)
	}
	return t
}

func rAppend(from []Renderer, truncate int, end Renderer) []Renderer {
	to := make([]Renderer, len(from)-truncate, len(from)-truncate+1)
	copy(to, from[:len(from)-truncate])
	return append(to, end)
}
