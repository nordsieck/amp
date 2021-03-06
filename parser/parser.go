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

func AliasDeclState(ss []State) []State {
	ss = tokenParserState(ss, token.IDENT)
	ss = tokenReaderState(ss, token.ASSIGN)
	ss = Type(ss)
	for i, s := range ss {
		ad := aliasDecl{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, ad)
	}
	return ss
}

type aliasDecl struct{ id, typ Renderer }

func (a aliasDecl) Render() []byte { return append(append(a.id.Render(), `=`...), a.typ.Render()...) }

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
func Arguments(ss []State) []State {
	ss = tokenParserState(ss, token.LPAREN)
	expr := ExpressionListState(ss)
	expr = append(expr, tokenParserState(expr, token.ELLIPSIS)...)
	expr = append(expr, tokenParserState(expr, token.COMMA)...)
	ss = append(ss, expr...)
	ss = tokenReaderState(ss, token.RPAREN)

	for i, s := range ss {
		var arg arguments
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
			arg.comma = true
			s.r = s.r[:len(s.r)-1]
		}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.ELLIPSIS {
			arg.ellipsis = true
			s.r = s.r[:len(s.r)-1]
		}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.LPAREN {
			ss[i].r = rAppend(s.r, 1, arg)
		} else {
			arg.expressionList = s.r[len(s.r)-1]
			ss[i].r = rAppend(s.r, 2, arg)
		}
	}
	return ss
}

type arguments struct {
	expressionList  Renderer
	ellipsis, comma bool
}

func (a arguments) Render() []byte {
	ret := []byte(`(`)
	if a.expressionList != nil {
		ret = append(ret, a.expressionList.Render()...)
	}
	if a.ellipsis {
		ret = append(ret, `...`...)
	}
	if a.comma {
		ret = append(ret, `,`...)
	}
	return append(ret, `)`...)
}

func ArrayType(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = ExpressionState(ss)
	ss = tokenReaderState(ss, token.RBRACK)
	ss = Type(ss)

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

func AssignmentState(ss []State) []State {
	ss = ExpressionListState(ss)
	ss = AssignOp(ss)
	ss = ExpressionListState(ss)
	for i, s := range ss {
		a := assignment{s.r[len(s.r)-3], s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 3, a)
	}
	return ss
}

type assignment struct {
	first Renderer
	op    Renderer
	last  Renderer
}

func (a assignment) Render() []byte {
	return append(append(a.first.Render(), a.op.Render()...), a.last.Render()...)
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

// bad spec
// "{" StatementList [ ";" ] "}"
// omit optional semi because of empty stmt
func BlockState(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACE)
	ss = StatementListState(ss)
	ss = tokenReaderState(ss, token.RBRACE)
	for i, s := range ss {
		b := block{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, b)
	}
	return ss
}

type block struct{ r Renderer }

func (b block) Render() []byte {
	ret := []byte(`{`)
	if b.r != nil {
		ret = append(ret, b.r.Render()...)
	}
	return append(ret, `}`...)
}

func BreakStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.BREAK)
	return append(ts, tokenReader(ts, token.IDENT)...)
}

func BreakStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.BREAK)
	lbl := tokenParserState(ss, token.IDENT)

	for i, l := range lbl {
		bs := breakStmt{l.r[len(l.r)-1]}
		lbl[i].r = rAppend(l.r, 1, bs)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, breakStmt{})
	}
	return append(ss, lbl...)
}

type breakStmt struct{ r Renderer }

func (b breakStmt) Render() []byte {
	if b.r == nil {
		return []byte(`break`)
	}
	return append([]byte(`break `), b.r.Render()...)
}

func ChannelType(ss []State) []State {
	plain := tokenReaderState(ss, token.CHAN)
	after := tokenReaderState(plain, token.ARROW)
	before := tokenReaderState(ss, token.ARROW)
	before = tokenReaderState(before, token.CHAN)

	if len(plain) != 0 {
		plain = Type(plain)
		for _, s := range plain {
			ct := channelType{nil, s.r[len(s.r)-1]}
			s.r[len(s.r)-1] = ct
		}
	}
	if len(after) != 0 {
		after = Type(after)
		for _, s := range after {
			ct := channelType{&_false, s.r[len(s.r)-1]}
			s.r[len(s.r)-1] = ct
		}
	}
	if len(before) != 0 {
		before = Type(before)
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

func CommCaseState(ss []State) []State {
	def := tokenReaderState(ss, token.DEFAULT)
	for i, d := range def {
		def[i].r = rAppend(d.r, 0, commCase{})
	}
	ss = tokenReaderState(ss, token.CASE)
	send := SendStmtState(ss)
	for i, s := range send {
		cc := commCase{s.r[len(s.r)-1], true}
		send[i].r = rAppend(s.r, 1, cc)
	}

	recv := RecvStmtState(ss)
	for i, r := range recv {
		cc := commCase{r.r[len(r.r)-1], false}
		recv[i].r = rAppend(r.r, 1, cc)
	}
	return append(append(send, recv...), def...)
}

type commCase struct {
	stmt Renderer
	send bool // here so we can tell what type stmt is
}

func (c commCase) Render() []byte {
	if c.stmt == nil {
		return []byte(`default`)
	}
	return append([]byte(`case `), c.stmt.Render()...)
}

func CommClause(ts [][]*Token) [][]*Token {
	ts = CommCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func CommClauseState(ss []State) []State {
	ss = CommCaseState(ss)
	ss = tokenReaderState(ss, token.COLON)
	ss = StatementListState(ss)
	for i, s := range ss {
		cc := commClause{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, cc)
	}
	return ss
}

type commClause struct{ commCase, stmtList Renderer }

func (c commClause) Render() []byte {
	return append(append(c.commCase.Render(), `:`...), c.stmtList.Render()...)
}

func CompositeLit(ts [][]*Token) [][]*Token {
	ts = fromState(LiteralType(toState(ts)))
	return LiteralValue(ts)
}

func CompositeLitState(ss []State) []State {
	ss = LiteralType(ss)
	ss = LiteralValueState(ss)

	for i, s := range ss {
		cl := compositeLit{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, cl)
	}
	return ss
}

type compositeLit struct{ typ, val Renderer }

func (cl compositeLit) Render() []byte { return append(cl.typ.Render(), cl.val.Render()...) }

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
// "const" ( ConstSpec | "(" [ ConstSpec { ";" ConstSpec } [ ";" ]] ")" )
func ConstDeclState(ss []State) []State {
	ss = tokenParserState(ss, token.CONST)

	bare := ConstSpecState(ss)
	for i, b := range bare {
		cd := constDecl{r: []Renderer{b.r[len(b.r)-1]}}
		bare[i].r = rAppend(b.r, 2, cd)
	}

	ss = tokenReaderState(ss, token.LPAREN)
	loop := ConstSpecState(ss)
	ss = append(ss, loop...)

	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = ConstSpecState(loop)
		ss = append(ss, loop...)
	}
	ss = append(ss, tokenParserState(ss, token.SEMICOLON)...)
	ss = tokenReaderState(ss, token.RPAREN)

	for i, s := range ss {
		cd := constDecl{parens: true}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.SEMICOLON {
			cd.trailingSemi = true
			s.r = s.r[:len(s.r)-1]
		}
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.CONST {
				break
			}
			cd.r = append(cd.r, s.r[len(s.r)-1])
			s.r = s.r[:len(s.r)-1]
		}

		for j := 0; j < len(cd.r)/2; j++ {
			cd.r[j], cd.r[len(cd.r)-j-1] = cd.r[len(cd.r)-j-1], cd.r[j]
		}
		ss[i].r = rAppend(s.r, 1, cd)
	}
	return append(bare, ss...)
}

type constDecl struct {
	r                    []Renderer
	parens, trailingSemi bool
}

func (c constDecl) Render() []byte {
	ret := []byte(`const `)
	if !c.parens {
		return append(ret, c.r[0].Render()...)
	}

	ret = append(ret, `(`...)
	if len(c.r) == 0 {
		return append(ret, `)`...)
	}
	ret = append(ret, c.r[0].Render()...)
	for i := 1; i < len(c.r); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, c.r[i].Render()...)
	}
	if c.trailingSemi {
		ret = append(ret, `;`...)
	}
	return append(ret, `)`...)
}

// bad spec
// IdentifierLit [ Type ] "=" ExpressionList
func ConstSpec(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	ts = append(ts, fromState(Type(toState(ts)))...)
	ts = tokenReader(ts, token.ASSIGN)
	return ExpressionList(ts)
}

// bad spec
// IdentifierLit [ Type ] "=" ExpressionList
func ConstSpecState(ss []State) []State {
	ss = IdentifierList(ss)

	typ := Type(ss)
	typ = tokenReaderState(typ, token.ASSIGN)
	typ = ExpressionListState(typ)
	for i, t := range typ {
		cs := constSpec{t.r[len(t.r)-3], t.r[len(t.r)-2], t.r[len(t.r)-1]}
		typ[i].r = rAppend(t.r, 3, cs)
	}

	ss = tokenReaderState(ss, token.ASSIGN)
	ss = ExpressionListState(ss)
	for i, s := range ss {
		cs := constSpec{idList: s.r[len(s.r)-2], exprList: s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, cs)
	}

	return append(typ, ss...)
}

type constSpec struct{ idList, typ, exprList Renderer }

func (c constSpec) Render() []byte {
	ret := c.idList.Render()
	if c.typ != nil {
		ret = append(ret, ` `...)
		ret = append(ret, c.typ.Render()...)
	}
	ret = append(ret, `=`...)
	return append(ret, c.exprList.Render()...)
}

func ContinueStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.CONTINUE)
	return append(ts, tokenReader(ts, token.IDENT)...)
}

func ContinueStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.CONTINUE)
	label := tokenParserState(ss, token.IDENT)

	for i, l := range label {
		cs := continueStmt{l.r[len(l.r)-1]}
		label[i].r = rAppend(l.r, 1, cs)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, continueStmt{})
	}
	return append(ss, label...)
}

type continueStmt struct{ r Renderer }

func (c continueStmt) Render() []byte {
	if c.r == nil {
		return []byte(`continue`)
	}
	return append([]byte(`continue `), c.r.Render()...)
}

func Conversion(ts [][]*Token) [][]*Token {
	ts = fromState(Type(toState(ts)))
	ts = tokenReader(ts, token.LPAREN)
	if len(ts) == 0 {
		return nil
	}
	ts = Expression(ts)
	ts = append(ts, tokenReader(ts, token.COMMA)...)
	return tokenReader(ts, token.RPAREN)
}

func ConversionState(ss []State) []State {
	ss = Type(ss)
	ss = tokenReaderState(ss, token.LPAREN)
	if len(ss) == 0 {
		return nil
	}
	ss = ExpressionState(ss)
	ss = append(ss, tokenParserState(ss, token.COMMA)...)
	ss = tokenReaderState(ss, token.RPAREN)

	for i, s := range ss {
		var c conversion
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
			c.comma = true
			s.r = s.r[:len(s.r)-1]
		}
		c.expr = s.r[len(s.r)-1]
		c.typ = s.r[len(s.r)-2]

		typ := c.typ.(typ)
		_, tok := typ.r.(*Token)
		_, qi := typ.r.(qualifiedIdent)

		if tok || qi {
			ss[i].r = rAppend(s.r, 2, conversionOrArguments{c})
		} else {
			ss[i].r = rAppend(s.r, 2, c)
		}
	}
	return ss
}

type conversionOrArguments struct{ r Renderer }

func (c conversionOrArguments) Render() []byte { return c.r.Render() }

type conversion struct {
	typ, expr Renderer
	comma     bool
}

func (c conversion) Render() []byte {
	ret := append(append(c.typ.Render(), `(`...), c.expr.Render()...)
	if c.comma {
		ret = append(ret, `,`...)
	}
	return append(ret, `)`...)
}

func Declaration(ts [][]*Token) [][]*Token {
	return append(append(ConstDecl(ts), TypeDecl(ts)...), VarDecl(ts)...)
}

func DeclarationState(ss []State) []State {
	return append(append(ConstDeclState(ss), TypeDeclState(ss)...), VarDeclState(ss)...)
}

func DeferStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.DEFER)
	return Expression(ts)
}

func Element(ts [][]*Token) [][]*Token {
	return append(Expression(ts), LiteralValue(ts)...)
}

func ElementState(ss []State) []State {
	ss = append(ExpressionState(ss), LiteralValueState(ss)...)
	for i, s := range ss {
		e := element{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, e)
	}
	return ss
}

type element struct{ r Renderer }

func (e element) Render() []byte { return e.r.Render() }

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

func ElementListState(ss []State) []State {
	ss = KeyedElementState(ss)
	for i, s := range ss {
		el := elementList{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, el)
	}

	loop := ss
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.COMMA)
		loop = KeyedElementState(loop)
		for i, l := range loop {
			el := l.r[len(l.r)-2].(elementList)
			el = append(el, l.r[len(l.r)-1])
			loop[i].r = rAppend(l.r, 2, el)
		}

		ss = append(ss, loop...)
	}
	return ss
}

type elementList []Renderer

func (el elementList) Render() []byte {
	var ret []byte
	if len(el) == 0 {
		return nil
	}
	ret = el[0].Render()
	for i := 1; i < len(el); i++ {
		ret = append(ret, `,`...)
		ret = append(ret, el[i].Render()...)
	}
	return ret
}

func EllipsisArrayType(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	ss = tokenReaderState(ss, token.ELLIPSIS)
	ss = tokenReaderState(ss, token.RBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)

	for i, s := range ss {
		eat := ellipsisArrayType{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, eat)
	}
	return ss
}

type ellipsisArrayType struct{ r Renderer }

func (e ellipsisArrayType) Render() []byte { return append([]byte(`[...]`), e.r.Render()...) }

func EmptyStmt(ts [][]*Token) [][]*Token { return ts }

// TODO: figure out a way to minimize the use of EmptyStmt
func EmptyStmtState(ss []State) []State { return ss }

func ExprCaseClause(ts [][]*Token) [][]*Token {
	ts = ExprSwitchCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func ExprCaseClauseState(ss []State) []State {
	ss = ExprSwitchCaseState(ss)
	ss = tokenParserState(ss, token.COLON)
	ss = StatementListState(ss)
	for i, s := range ss {
		ecc := exprCaseClause{s.r[len(s.r)-3], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 3, ecc)
	}
	return ss
}

type exprCaseClause struct{ switchCase, stmtList Renderer }

// expects a stmtList, even if it's empty
func (e exprCaseClause) Render() []byte {
	return append(append(e.switchCase.Render(), `:`...), e.stmtList.Render()...)
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
	if len(ss) == 0 {
		return nil
	}

	base := UnaryExprState(ss)
	comp := BinaryOp(base)
	comp = ExpressionState(comp)

	for i, b := range base {
		expr := expression{b.r[len(b.r)-1]}
		base[i].r = rAppend(b.r, 1, expr)

	}

	for i, c := range comp {
		expr := expression{c.r[len(c.r)-3], c.r[len(c.r)-2], c.r[len(c.r)-1]}
		comp[i].r = rAppend(c.r, 3, expr)
	}

	return append(base, comp...)
}

// either 1 or 3 long
type expression []Renderer

func (e expression) Render() []byte {
	var ret []byte
	for _, elem := range e {
		ret = append(ret, elem.Render()...)
	}
	return ret
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

func ExpressionListState(ss []State) []State {
	ss = ExpressionState(ss)
	loop := ss
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = ExpressionState(loop)
		ss = append(ss, loop...)
	}
	for i, s := range ss {
		el := expressionList{s.r[len(s.r)-1]}
		s.r = s.r[:len(s.r)-1]
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
				el = append(el, s.r[len(s.r)-2])
				s.r = s.r[:len(s.r)-2]
			} else {
				break
			}
		}

		// reverse
		for i := 0; i < len(el)/2; i++ {
			el[i], el[len(el)-1-i] = el[len(el)-1-i], el[i]
		}

		ss[i].r = rAppend(s.r, 0, el)
	}
	return ss
}

type expressionList []Renderer

func (e expressionList) Render() []byte {
	if len(e) == 0 {
		return nil
	}
	var ret []byte
	for i := 0; i < len(e)-1; i++ {
		ret = append(ret, e[i].Render()...)
		ret = append(ret, `,`...)
	}
	return append(ret, e[len(e)-1].Render()...)
}

func ExpressionStmt(ts [][]*Token) [][]*Token { return Expression(ts) }

func ExpressionStmtState(ss []State) []State { return ExpressionState(ss) }

func ExprSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenReader(ts, token.CASE)
	cas = ExpressionList(cas)
	return append(tokenReader(ts, token.DEFAULT), cas...)
}

func ExprSwitchCaseState(ss []State) []State {
	cas := tokenReaderState(ss, token.CASE)
	cas = ExpressionListState(cas)
	for i, c := range cas {
		esc := exprSwitchCase{c.r[len(c.r)-1]}
		cas[i].r = rAppend(c.r, 1, esc)
	}

	def := tokenReaderState(ss, token.DEFAULT)
	for i, d := range def {
		def[i].r = rAppend(d.r, 0, exprSwitchCase{})
	}
	return append(def, cas...)
}

type exprSwitchCase struct{ r Renderer }

func (e exprSwitchCase) Render() []byte {
	if e.r == nil {
		return []byte(`default`)
	}
	return append([]byte(`case `), e.r.Render()...)
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

// bad spec
// "switch" [ SimpleStmt ";" ] [ Expression ] "{" [ ExprCaseClause { ";" ExprCaseClause } [ ";" ]] "}"
func ExprSwitchStmtState(ss []State) []State {
	ss = tokenParserState(ss, token.SWITCH)

	simple := SimpleStmtState(ss)
	simple = tokenReaderState(simple, token.SEMICOLON)
	for i, s := range simple {
		var ess exprSwitchStmt
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.SWITCH {
			ess.simple = e{}
		} else {
			ess.simple = s.r[len(s.r)-1]
			s.r = s.r[:len(s.r)-1]
		}
		simple[i].r = rAppend(s.r, 1, ess)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 1, exprSwitchStmt{})
	}
	ss = append(ss, simple...)

	expr := ExpressionState(ss)
	for i, e := range expr {
		ess := e.r[len(e.r)-2].(exprSwitchStmt)
		ess.expr = e.r[len(e.r)-1]
		expr[i].r = rAppend(e.r, 2, ess)
	}
	ss = append(ss, expr...)
	ss = tokenReaderState(ss, token.LBRACE)
	loop := ExprCaseClauseState(ss)
	for i, l := range loop {
		ess := l.r[len(l.r)-2].(exprSwitchStmt)
		ess.clauses = append(ess.clauses, l.r[len(l.r)-1])
		loop[i].r = rAppend(l.r, 2, ess)
	}
	ss = append(ss, loop...)
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = ExprCaseClauseState(loop)
		for i, l := range loop {
			ess := l.r[len(l.r)-2].(exprSwitchStmt)
			ess.clauses = append(ess.clauses, l.r[len(l.r)-1])
			loop[i].r = rAppend(l.r, 2, ess)
		}
		ss = append(ss, loop...)
	}
	return tokenReaderState(ss, token.RBRACE)
}

type exprSwitchStmt struct {
	simple, expr Renderer
	clauses      []Renderer
}

func (e exprSwitchStmt) Render() []byte {
	ret := []byte(`switch `)
	if e.simple != nil {
		ret = append(ret, e.simple.Render()...)
		ret = append(ret, `;`...)
	}
	if e.expr != nil {
		ret = append(ret, e.expr.Render()...)
	}
	ret = append(ret, `{`...)
	if len(e.clauses) == 0 {
		return append(ret, `}`...)
	}
	ret = append(ret, e.clauses[0].Render()...)
	for i := 1; i < len(e.clauses); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, e.clauses[i].Render()...)
	}
	return append(ret, `}`...)
}

func FallthroughStmt(ss []State) []State { return tokenParserState(ss, token.FALLTHROUGH) }

func FieldDecl(ss []State) []State {
	typ := IdentifierList(ss)
	if len(typ) != 0 {
		typ = Type(typ)
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
	ts = fromState(Signature(toState(ts)))
	return Block(ts)
}

func FunctionState(ss []State) []State {
	ss = Signature(ss)
	ss = BlockState(ss)
	for i, s := range ss {
		f := function{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, f)
	}
	return ss
}

type function struct{ sig, block Renderer }

func (f function) Render() []byte { return append(f.sig.Render(), f.block.Render()...) }

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

func FunctionLitState(ss []State) []State {
	ss = tokenReaderState(ss, token.FUNC)
	ss = FunctionState(ss)
	for i, s := range ss {
		fl := functionLit{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, fl)
	}
	return ss
}

type functionLit struct{ r Renderer }

func (f functionLit) Render() []byte { return append([]byte(`func`), f.r.Render()...) }

func FunctionType(ss []State) []State {
	ss = tokenReaderState(ss, token.FUNC)
	ss = Signature(ss)

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

func GoStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.GO)
	ss = ExpressionState(ss)
	for i, s := range ss {
		gs := goStmt{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, gs)
	}
	return ss
}

type goStmt struct{ r Renderer }

func (g goStmt) Render() []byte { return append([]byte(`go `), g.r.Render()...) }

func GotoStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.GOTO)
	return tokenReader(ts, token.IDENT)
}

func GotoStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.GOTO)
	ss = tokenParserState(ss, token.IDENT)
	for i, s := range ss {
		gs := gotoStmt{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, gs)
	}
	return ss
}

type gotoStmt struct{ r Renderer }

func (g gotoStmt) Render() []byte { return append([]byte(`goto `), g.r.Render()...) }

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

func IfStmtState(ss []State) []State {
	if len(ss) == 0 {
		return nil
	}
	ss = tokenParserState(ss, token.IF)
	simple := SimpleStmtState(ss)
	simple = tokenReaderState(simple, token.SEMICOLON)
	simple = ExpressionState(simple)
	simple = BlockState(simple)
	simpleElse := tokenReaderState(simple, token.ELSE)
	simpleElse = append(IfStmtState(simpleElse), BlockState(simpleElse)...)

	for i, s := range simpleElse {
		is := ifStmt{nil, s.r[len(s.r)-3], s.r[len(s.r)-2], s.r[len(s.r)-1]}
		if tok, ok := s.r[len(s.r)-4].(*Token); ok && tok.tok == token.IF {
			is.simple = e{}
			simpleElse[i].r = rAppend(s.r, 4, is)
		} else {
			is.simple = s.r[len(s.r)-4]
			simpleElse[i].r = rAppend(s.r, 5, is)
		}
	}
	for i, s := range simple {
		is := ifStmt{nil, s.r[len(s.r)-2], s.r[len(s.r)-1], nil}
		if tok, ok := s.r[len(s.r)-3].(*Token); ok && tok.tok == token.IF {
			is.simple = e{}
			simple[i].r = rAppend(s.r, 3, is)
		} else {
			is.simple = s.r[len(s.r)-3]
			simple[i].r = rAppend(s.r, 4, is)
		}
	}

	ss = ExpressionState(ss)
	ss = BlockState(ss)
	ssElse := tokenReaderState(ss, token.ELSE)
	ssElse = append(IfStmtState(ssElse), BlockState(ssElse)...)

	for i, s := range ssElse {
		is := ifStmt{nil, s.r[len(s.r)-3], s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ssElse[i].r = rAppend(s.r, 4, is)
	}
	for i, s := range ss {
		is := ifStmt{nil, s.r[len(s.r)-2], s.r[len(s.r)-1], nil}
		ss[i].r = rAppend(s.r, 3, is)
	}

	return append(append(ss, ssElse...), append(simple, simpleElse...)...)
}

type ifStmt struct{ simple, expr, block, tail Renderer }

func (i ifStmt) Render() []byte {
	ret := []byte(`if `)
	if i.simple != nil {
		ret = append(ret, i.simple.Render()...)
		ret = append(ret, `;`...)
	}
	ret = append(ret, i.expr.Render()...)
	ret = append(ret, i.block.Render()...)
	if i.tail != nil {
		ret = append(ret, ` else `...)
		ret = append(ret, i.tail.Render()...)
	}
	return ret
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

func IncDecStmtState(ss []State) []State {
	ss = ExpressionState(ss)
	inc := tokenReaderState(ss, token.INC)
	for i, n := range inc {
		ids := incDecStmt{n.r[len(n.r)-1], true}
		inc[i].r = rAppend(n.r, 1, ids)
	}

	dec := tokenReaderState(ss, token.DEC)
	for i, d := range dec {
		ids := incDecStmt{d.r[len(d.r)-1], false}
		dec[i].r = rAppend(d.r, 1, ids)
	}
	return append(inc, dec...)
}

type incDecStmt struct {
	r   Renderer
	inc bool
}

func (ids incDecStmt) Render() []byte {
	ret := ids.r.Render()
	if ids.inc {
		return append(ret, `++`...)
	}
	return append(ret, `--`...)
}

func Index(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = ExpressionState(ss)
	ss = tokenReaderState(ss, token.RBRACK)
	for i, s := range ss {
		idx := index{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, idx)
	}
	return ss
}

type index struct{ r Renderer }

func (i index) Render() []byte { return append(append([]byte(`[`), i.r.Render()...), `]`...) }

// bad spec
// "interface" "{" [ MethodSpec { ";" MethodSpec } [ ";" ]] "}"
func InterfaceType(ss []State) []State {
	ss = tokenParserState(ss, token.INTERFACE)
	ss = tokenReaderState(ss, token.LBRACE)
	methods := MethodSpec(ss)
	loop := methods
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = MethodSpec(loop)
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

func Key(ts [][]*Token) [][]*Token {
	return append(append(tokenReader(ts, token.IDENT), Expression(ts)...),
		LiteralValue(ts)...)
}

func KeyState(ss []State) []State {
	// TODO: do we need field name at all?  Currently we use OperandName in Expression
	expr := ExpressionState(ss)
	lv := LiteralValueState(ss)
	ss = append(expr, lv...)

	for i, s := range ss {
		k := key{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, k)
	}

	return ss
}

type key struct{ r Renderer }

func (k key) Render() []byte { return k.r.Render() }

func KeyedElement(ts [][]*Token) [][]*Token {
	with := Key(ts)
	with = tokenReader(with, token.COLON)
	return append(Element(ts), Element(with)...)
}

func KeyedElementState(ss []State) []State {
	with := KeyState(ss)
	with = tokenReaderState(with, token.COLON)
	with = ElementState(with)
	for i, w := range with {
		ke := keyedElement{key: w.r[len(w.r)-2], element: w.r[len(w.r)-1]}
		with[i].r = rAppend(w.r, 2, ke)
	}

	solo := ElementState(ss)
	for i, s := range solo {
		ke := keyedElement{element: s.r[len(s.r)-1]}
		solo[i].r = rAppend(s.r, 1, ke)
	}

	return append(solo, with...)
}

type keyedElement struct{ key, element Renderer }

func (ke keyedElement) Render() []byte {
	var ret []byte
	if ke.key != nil {
		ret = append(ret, ke.key.Render()...)
		ret = append(ret, `:`...)
	}
	return append(ret, ke.element.Render()...)
}

func LabeledStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IDENT)
	ts = tokenReader(ts, token.COLON)
	if len(ts) == 0 {
		return nil
	}
	return Statement(ts)
}

func LabeledStmtState(ss []State) []State {
	ss = tokenParserState(ss, token.IDENT)
	ss = tokenParserState(ss, token.COLON)
	ss = StatementState(ss)
	for i, s := range ss {
		var ls labeledStmt
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COLON {
			ls.label = s.r[len(s.r)-2]
			s.r = s.r[:len(s.r)-2]
		} else {
			ls.label = s.r[len(s.r)-3]
			ls.stmt = s.r[len(s.r)-1]
			s.r = s.r[:len(s.r)-3]
		}
		ss[i].r = rAppend(s.r, 0, ls)
	}
	return ss
}

type labeledStmt struct{ label, stmt Renderer }

func (l labeledStmt) Render() []byte {
	ret := append(l.label.Render(), `:`...)
	if l.stmt != nil {
		return append(ret, l.stmt.Render()...)
	}
	return ret
}

func Literal(ts [][]*Token) [][]*Token {
	basicLit := fromState(BasicLit(toState(ts)))
	return append(append(basicLit, CompositeLit(ts)...), FunctionLit(ts)...)
}

func LiteralState(ss []State) []State {
	return append(append(BasicLit(ss), CompositeLitState(ss)...), FunctionLitState(ss)...)
}

func LiteralType(ss []State) []State {
	return append(
		append(append(StructType(ss), EllipsisArrayType(ss)...),
			append(SliceType(ss), MapType(ss)...)...),
		append(ArrayType(ss), TypeName(ss)...)...)
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

func LiteralValueState(ss []State) []State {
	ss = tokenParserState(ss, token.LBRACE)

	if len(ss) == 0 {
		return nil
	}

	list := ElementListState(ss)
	list = append(list, tokenParserState(list, token.COMMA)...)
	ss = append(ss, list...)
	ss = tokenReaderState(ss, token.RBRACE)

	for i, s := range ss {
		var lv literalValue

		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.COMMA {
			lv.comma = true
			s.r = s.r[:len(s.r)-1]
		}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.LBRACE {
			ss[i].r = rAppend(s.r, 1, lv)
		} else {
			lv.elementList = s.r[len(s.r)-1]
			ss[i].r = rAppend(s.r, 2, lv)
		}
	}

	return ss
}

type literalValue struct {
	elementList Renderer
	comma       bool
}

func (l literalValue) Render() []byte {
	ret := []byte(`{`)
	if l.elementList != nil {
		ret = append(ret, l.elementList.Render()...)
		if l.comma {
			ret = append(ret, `,`...)
		}
	}
	return append(ret, `}`...)
}

func MapType(ss []State) []State {
	ss = tokenReaderState(ss, token.MAP)
	ss = tokenReaderState(ss, token.LBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)
	ss = tokenReaderState(ss, token.RBRACK)
	ss = Type(ss)
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
	ts = fromState(Parameters(toState(ts)))
	ts = fromState(nonBlankIdent(toState(ts)))
	return Function(ts)
}

// TODO: double check that the OR case is actually doing what it's supposed to do
func MethodExpr(ts [][]*Token) [][]*Token {
	ts = ReceiverType(ts)
	ts = tokenReader(ts, token.PERIOD)
	return tokenReader(ts, token.IDENT)
}

func MethodExprState(ss []State) []State {
	ss = ReceiverTypeState(ss)
	ss = tokenReaderState(ss, token.PERIOD)
	ss = tokenParserState(ss, token.IDENT)

	for i, s := range ss {
		me := methodExpr{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		rt := me.receiverType.(receiverType)
		if tok, ok := rt.r.(*Token); ok && tok.tok == token.IDENT && rt.parens == 0 && !rt.pointer {
			or := operandNameOrMethodExpr{qualifiedIdent{me.receiverType, me.methodName}}
			ss[i].r = rAppend(s.r, 2, or)
		} else {
			ss[i].r = rAppend(s.r, 2, me)
		}
	}
	return ss
}

type methodExpr struct{ receiverType, methodName Renderer }

func (m methodExpr) Render() []byte {
	return append(append(m.receiverType.Render(), `.`...), m.methodName.Render()...)
}

func MethodSpec(ss []State) []State {
	sig := nonBlankIdent(ss)
	sig = Signature(sig)
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

func nonBlankIdent(ss []State) []State {
	var result []State
	for _, s := range ss {
		if p := pop(&s.t); p != nil && p.tok == token.IDENT && p.lit != `_` {
			result = append(result, State{append(s.r, p), s.t})
		}
	}
	return result
}

func NonEmptyStatementState(ss []State) []State {
	if len(ss) == 0 {
		return nil
	}
	return append(nonEmptySimpleStmtState(ss), nonSimpleStatementState(ss)...)
}

func nonEmptySimpleStmtState(ss []State) []State {
	return append(
		append(append(ExpressionStmtState(ss), SendStmtState(ss)...),
			append(IncDecStmtState(ss), AssignmentState(ss)...)...),
		ShortVarDeclState(ss)...)
}

func nonSimpleStatementState(ss []State) []State {
	return append(append(
		append(append(DeclarationState(ss), LabeledStmtState(ss)...), append(GoStmtState(ss), ReturnStmtState(ss)...)...),
		append(append(BreakStmtState(ss), ContinueStmtState(ss)...), append(GotoStmtState(ss), FallthroughStmt(ss)...)...)...),
		append(append(BlockState(ss), IfStmtState(ss)...), append(SwitchStmtState(ss), SelectStmtState(ss)...)...)...)
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
	parens := tokenReaderState(ss, token.LPAREN)
	if len(parens) != 0 {
		parens = ExpressionState(parens)
	}
	parens = tokenReaderState(parens, token.RPAREN)
	for i, p := range parens {
		o := operand{p.r[len(p.r)-1], true}
		parens[i].r = rAppend(p.r, 1, o)
	}

	noParens := append(append(LiteralState(ss), OperandNameState(ss)...), MethodExprState(ss)...)
	for i, n := range noParens {
		o := operand{n.r[len(n.r)-1], false}
		noParens[i].r = rAppend(n.r, 1, o)
	}
	return append(noParens, parens...)
}

type operand struct {
	r      Renderer
	parens bool
}

func (o operand) Render() []byte {
	ret := o.r.Render()
	if o.parens {
		ret = append(append([]byte(`(`), ret...), `)`...)
	}
	return ret
}

func OperandName(ts [][]*Token) [][]*Token {
	qi := fromState(QualifiedIdent(toState(ts)))
	return append(tokenReader(ts, token.IDENT), qi...)
}

func OperandNameState(ss []State) []State { return tokenParserState(ss, token.IDENT) }

type operandNameOrMethodExpr struct{ r Renderer }

func (o operandNameOrMethodExpr) Render() []byte { return o.r.Render() }

func PackageClause(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.PACKAGE)
	return fromState(PackageName(toState(ts)))
}

func PackageName(ss []State) []State { return nonBlankIdent(ss) }

func ParameterDeclIDList(ss []State) []State {
	ss = IdentifierList(ss)
	ss = append(ss, tokenParserState(ss, token.ELLIPSIS)...)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)

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

func ParameterDeclNoList(ss []State) []State {
	ss = append(ss, tokenParserState(ss, token.ELLIPSIS)...)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)

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

func ParameterList(ss []State) []State {
	list := ParameterDeclIDList(ss)
	loop := list
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = ParameterDeclIDList(loop)
		list = append(list, loop...)
	}

	nolist := ParameterDeclNoList(ss)
	loop = nolist
	for len(loop) != 0 {
		loop = tokenParserState(loop, token.COMMA)
		loop = ParameterDeclNoList(loop)
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

func Parameters(ss []State) []State {
	ss = tokenParserState(ss, token.LPAREN)
	params := ParameterList(ss)
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

func PointerType(ss []State) []State {
	ss = tokenReaderState(ss, token.MUL)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)
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
		additions := fromState(Selector(toState(newBase)))
		additions = append(additions, fromState(Index(toState(newBase)))...)
		additions = append(additions, fromState(Slice(toState(newBase)))...)
		additions = append(additions, fromState(TypeAssertion(toState(newBase)))...)
		additions = append(additions, fromState(Arguments(toState(newBase)))...)
		base = append(base, additions...)
		newBase = additions
	}
	return base
}

func PrimaryExprState(ss []State) []State {
	ss = append(OperandState(ss), ConversionState(ss)...)
	for i, s := range ss {
		pe := primaryExpr{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, pe)
	}

	loop := ss
	for len(loop) != 0 {
		sel := Selector(loop)
		idx := Index(loop)
		slc := Slice(loop)
		typ := TypeAssertion(loop)
		arg := Arguments(loop)
		loop = append(append(append(sel, idx...), append(slc, typ...)...), arg...)
		for i, s := range loop {
			pe := s.r[len(s.r)-2].(primaryExpr)
			pe = append(pe, s.r[len(s.r)-1])
			loop[i].r = rAppend(s.r, 2, pe)
		}

		ss = append(ss, loop...)
	}

	// filter for uniqueness of parse
	var newss []State
	for _, s := range ss {
		pe := s.r[len(s.r)-1].(primaryExpr)
		if len(pe) == 1 {
			newss = append(newss, s)
			continue
		}

		if op, ok := pe[0].(operand); ok {
			if _, sel := pe[1].(selector); sel {
				// a.a or (a).a
				_, tok := op.r.(*Token)
				_, or := op.r.(operandNameOrMethodExpr)

				if tok || or || op.parens {
					continue
				}
			}
			if arg, ok := pe[1].(arguments); ok && !arg.ellipsis {
				var el expressionList
				if arg.expressionList != nil {
					el = arg.expressionList.(expressionList)
				}
				if len(el) != 1 {
					newss = append(newss, s)
					continue
				}

				// get the root level operand
				add := false
				for op.parens {
					expr, ok := op.r.(expression)
					if !ok || len(expr) != 1 {
						add = true
						break
					}

					ue, ok := expr[0].(unaryExpr)
					if !ok || len(ue) != 1 {
						add = true
						break
					}

					pe, ok := ue[0].(primaryExpr)
					if !ok || len(pe) != 1 {
						add = true
						break
					}

					if op, ok = pe[0].(operand); !ok {
						add = true
						break
					}
				}

				if add {
					newss = append(newss, s)
					continue
				}

				// (a)(a)
				if tok, ok := op.r.(*Token); ok && tok.tok == token.IDENT {
					continue
				}

				// (a.a)(a)
				if _, ok = op.r.(operandNameOrMethodExpr); ok {
					continue
				}
			}
		}
		newss = append(newss, s)
	}

	return newss
}

type primaryExpr []Renderer

func (pe primaryExpr) Render() []byte {
	var out []byte
	for _, elem := range pe {
		out = append(out, elem.Render()...)
	}
	return out
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

func RangeClauseState(ss []State) []State {
	expr := ExpressionListState(ss)
	expr = tokenReaderState(expr, token.ASSIGN)
	for i, e := range expr {
		rc := rangeClause{e.r[len(e.r)-1], nil, true}
		expr[i].r = rAppend(e.r, 1, rc)
	}

	ids := IdentifierList(ss)
	ids = tokenReaderState(ids, token.DEFINE)
	for i, id := range ids {
		rc := rangeClause{id.r[len(id.r)-1], nil, false}
		ids[i].r = rAppend(id.r, 1, rc)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, rangeClause{})
	}
	ss = append(append(ss, expr...), ids...)
	ss = tokenReaderState(ss, token.RANGE)
	ss = ExpressionState(ss)
	for i, s := range ss {
		rc := s.r[len(s.r)-2].(rangeClause)
		rc.expr = s.r[len(s.r)-1]
		ss[i].r = rAppend(s.r, 2, rc)
	}
	return ss
}

type rangeClause struct {
	list, expr Renderer
	assign     bool
}

func (r rangeClause) Render() []byte {
	var ret []byte
	if r.list != nil {
		if r.assign {
			ret = append(r.list.Render(), `=`...)
		} else {
			ret = append(r.list.Render(), `:=`...)
		}
	}
	ret = append(ret, `range `...)
	return append(ret, r.expr.Render()...)
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

func ReceiverTypeState(ss []State) []State {
	paren := tokenReaderState(ss, token.LPAREN)
	if len(paren) != 0 {
		paren = ReceiverTypeState(paren)
	}
	paren = tokenReaderState(paren, token.RPAREN)
	for i, p := range paren {
		rt := p.r[len(p.r)-1].(receiverType)
		rt.parens += 1
		paren[i].r = rAppend(p.r, 1, rt)
	}

	tn := TypeName(ss)
	for i, t := range tn {
		rt := receiverType{r: t.r[len(t.r)-1]}
		tn[i].r = rAppend(t.r, 1, rt)
	}

	ptr := tokenParserState(ss, token.LPAREN)
	ptr = tokenParserState(ptr, token.MUL)
	ptr = TypeName(ptr)
	ptr = tokenParserState(ptr, token.RPAREN)
	for i, p := range ptr {
		rt := receiverType{p.r[len(p.r)-2], 1, true}
		ptr[i].r = rAppend(p.r, 4, rt)
	}
	return append(append(tn, ptr...), paren...)
}

type receiverType struct {
	r       Renderer
	parens  int
	pointer bool
}

func (r receiverType) Render() []byte {
	var ret []byte
	for i := 0; i < r.parens; i++ {
		ret = append(ret, `(`...)
	}
	if r.pointer {
		ret = append(ret, `*`...)
	}
	ret = append(ret, r.r.Render()...)
	for i := 0; i < r.parens; i++ {
		ret = append(ret, `)`...)
	}
	return ret
}

func RecvStmt(ts [][]*Token) [][]*Token {
	expr := ExpressionList(ts)
	expr = tokenReader(expr, token.ASSIGN)
	ident := fromState(IdentifierList(toState(ts)))
	ident = tokenReader(ident, token.DEFINE)
	return Expression(append(ts, append(expr, ident...)...))
}

func RecvStmtState(ss []State) []State {
	expr := ExpressionListState(ss)
	expr = tokenReaderState(expr, token.ASSIGN)
	expr = ExpressionState(expr)
	for i, e := range expr {
		rs := recvStmt{e.r[len(e.r)-2], e.r[len(e.r)-1], false}
		expr[i].r = rAppend(e.r, 2, rs)
	}

	ids := IdentifierList(ss)
	ids = tokenReaderState(ids, token.DEFINE)
	ids = ExpressionState(ids)
	for i, id := range ids {
		rs := recvStmt{id.r[len(id.r)-2], id.r[len(id.r)-1], true}
		ids[i].r = rAppend(id.r, 2, rs)
	}

	ss = ExpressionState(ss)
	for i, s := range ss {
		rs := recvStmt{expr: s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, rs)
	}
	return append(append(expr, ids...), ss...)
}

type recvStmt struct {
	list, expr Renderer
	define     bool
}

func (r recvStmt) Render() []byte {
	var ret []byte
	if r.list != nil {
		if r.define {
			ret = append(r.list.Render(), `:=`...)
		} else {
			ret = append(r.list.Render(), `=`...)
		}
	}
	return append(ret, r.expr.Render()...)
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

func Result(ss []State) []State {
	if len(ss) == 0 {
		return nil
	}

	pp := Parameters(ss)
	for i, p := range pp {
		r := result{parameters: p.r[len(p.r)-1]}
		pp[i].r = rAppend(p.r, 1, r)
	}

	tt := Type(ss)
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

func ReturnStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.RETURN)
	list := ExpressionListState(ss)
	for i, l := range list {
		rs := returnStmt{l.r[len(l.r)-1]}
		list[i].r = rAppend(l.r, 1, rs)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, returnStmt{})
	}
	return append(ss, list...)
}

type returnStmt struct{ r Renderer }

func (r returnStmt) Render() []byte {
	if r.r == nil {
		return []byte(`return`)
	}
	return append([]byte(`return `), r.r.Render()...)
}

func Selector(ss []State) []State {
	ss = tokenReaderState(ss, token.PERIOD)
	ss = tokenParserState(ss, token.IDENT)
	for i, s := range ss {
		sel := selector{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, sel)
	}
	return ss
}

type selector struct{ r Renderer }

func (s selector) Render() []byte { return append([]byte(`.`), s.r.Render()...) }

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

// bad spec
// "select" "{" [ CommClause { ";" CommClause } [ ";" ]] "}"
func SelectStmtState(ss []State) []State {
	ss = tokenReaderState(ss, token.SELECT)
	ss = tokenReaderState(ss, token.LBRACE)

	loop := CommClauseState(ss)
	for i, l := range loop {
		sel := selectStmt{l.r[len(l.r)-1]}
		loop[i].r = rAppend(l.r, 1, sel)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, selectStmt{})
	}
	ss = append(ss, loop...)
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = CommClauseState(loop)
		for i, l := range loop {
			sel := l.r[len(l.r)-2].(selectStmt)
			sel = append(sel, l.r[len(l.r)-1])
			loop[i].r = rAppend(l.r, 2, sel)
		}
		ss = append(ss, loop...)
	}

	return tokenReaderState(ss, token.RBRACE)
}

type selectStmt []Renderer

func (s selectStmt) Render() []byte {
	ret := []byte(`select {`)
	if len(s) == 0 {
		return append(ret, `}`...)
	}
	ret = append(ret, s[0].Render()...)
	for i := 1; i < len(s); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, s[i].Render()...)
	}
	return append(ret, `}`...)
}

func SendStmt(ts [][]*Token) [][]*Token {
	ts = Expression(ts)
	ts = tokenReader(ts, token.ARROW)
	return Expression(ts)
}

func SendStmtState(ss []State) []State {
	ss = ExpressionState(ss)
	ss = tokenReaderState(ss, token.ARROW)
	ss = ExpressionState(ss)
	for i, s := range ss {
		snd := sendStmt{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, snd)
	}
	return ss
}

type sendStmt struct{ channel, expression Renderer }

func (ss sendStmt) Render() []byte {
	return append(append(ss.channel.Render(), `<-`...), ss.expression.Render()...)
}

func ShortVarDecl(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	ts = tokenReader(ts, token.DEFINE)
	return ExpressionList(ts)
}

func ShortVarDeclState(ss []State) []State {
	ss = IdentifierList(ss)
	ss = tokenReaderState(ss, token.DEFINE)
	ss = ExpressionListState(ss)
	for i, s := range ss {
		svd := shortVarDecl{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, svd)
	}
	return ss
}

type shortVarDecl struct{ id, expr Renderer }

func (svd shortVarDecl) Render() []byte {
	return append(append(svd.id.Render(), `:=`...), svd.expr.Render()...)
}

func Signature(ss []State) []State {
	pp := Parameters(ss)
	pr := Result(pp)

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

// TODO wrap the return of this function so I don't have to worry about returning nothing.
func SimpleStmtState(ss []State) []State {
	return append(EmptyStmtState(ss), nonEmptySimpleStmtState(ss)...)
}

func Slice(ss []State) []State {
	ss = tokenParserState(ss, token.LBRACK)
	ss = append(ss, ExpressionState(ss)...)
	ss = tokenParserState(ss, token.COLON)

	two := append(ss, ExpressionState(ss)...)
	two = tokenReaderState(two, token.RBRACK)
	for i, t := range two {
		s := slice{}
		if tok, ok := t.r[len(t.r)-1].(*Token); ok && tok.tok == token.COLON {
			t.r = t.r[:len(t.r)-1]
		} else {
			s[1] = t.r[len(t.r)-1]
			t.r = t.r[:len(t.r)-2]
		}
		if tok, ok := t.r[len(t.r)-1].(*Token); ok && tok.tok == token.LBRACK {
			two[i].r = rAppend(t.r, 1, s)
		} else {
			s[0] = t.r[len(t.r)-1]
			two[i].r = rAppend(t.r, 2, s)
		}
	}

	three := ExpressionState(ss)
	three = tokenReaderState(three, token.COLON)
	three = ExpressionState(three)
	three = tokenReaderState(three, token.RBRACK)
	for i, t := range three {
		s := slice{1: t.r[len(t.r)-2], 2: t.r[len(t.r)-1]}
		if tok, ok := t.r[len(t.r)-4].(*Token); ok && tok.tok == token.LBRACK {
			three[i].r = rAppend(t.r, 4, s)
		} else {
			s[0] = t.r[len(t.r)-4]
			three[i].r = rAppend(t.r, 5, s)
		}
	}

	return append(two, three...)
}

type slice [3]Renderer

func (s slice) Render() []byte {
	ret := []byte(`[`)
	if s[0] != nil {
		ret = append(ret, s[0].Render()...)
	}
	ret = append(ret, `:`...)
	if s[1] != nil {
		ret = append(ret, s[1].Render()...)
	}
	if s[2] != nil {
		ret = append(ret, `:`...)
		ret = append(ret, s[2].Render()...)
	}
	return append(ret, `]`...)
}

func SliceType(ss []State) []State {
	ss = tokenReaderState(ss, token.LBRACK)
	ss = tokenReaderState(ss, token.RBRACK)
	if len(ss) == 0 {
		return nil
	}
	ss = Type(ss)
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

// TODO: currently the medium-term goal
func Statement(ts [][]*Token) [][]*Token {
	fallthroughStmt := fromState(FallthroughStmt(toState(ts)))
	return append(
		append(append(append(Declaration(ts), LabeledStmt(ts)...), append(SimpleStmt(ts), GoStmt(ts)...)...),
			append(append(ReturnStmt(ts), BreakStmt(ts)...), append(ContinueStmt(ts), GotoStmt(ts)...)...)...),
		append(append(append(fallthroughStmt, Block(ts)...), append(IfStmt(ts), SwitchStmt(ts)...)...),
			append(append(SelectStmt(ts), ForStmt(ts)...), DeferStmt(ts)...)...)...)
}

func StatementState(ss []State) []State {
	if len(ss) == 0 {
		return nil
	}
	return append(nonSimpleStatementState(ss), SimpleStmtState(ss)...)
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
// Statement { ";" Statement }
// force empty stmt
func StatementListState(ss []State) []State {

	// very important that full is processed before empty.
	// empty stmt interacts strangely with append
	full := NonEmptyStatementState(ss)
	for i, f := range full {
		sl := statementList{f.r[len(f.r)-1]}
		full[i].r = rAppend(f.r, 1, sl)
	}
	empty := EmptyStmtState(ss)
	for i, em := range empty {
		sl := statementList{e{}}
		empty[i].r = rAppend(em.r, 0, sl)
	}

	loop := append(empty, full...)
	ss = loop
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		full = NonEmptyStatementState(loop)
		for i, f := range full {
			sl := f.r[len(f.r)-2].(statementList)
			sl = append(sl, f.r[len(f.r)-1])
			full[i].r = rAppend(f.r, 2, sl)
		}

		empty = EmptyStmtState(loop)
		for i, em := range empty {
			sl := em.r[len(em.r)-1].(statementList)
			sl = append(sl, e{})
			empty[i].r = rAppend(em.r, 1, sl)
		}

		loop = append(empty, full...)
		ss = append(ss, loop...)
	}
	return ss
}

type statementList []Renderer

func (s statementList) Render() []byte {
	if len(s) == 0 {
		return nil
	}
	ret := s[0].Render()
	for i := 1; i < len(s); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, s[i].Render()...)
	}
	return ret
}

// bad spec
// "struct" "{" [ FieldDecl { ";" FieldDecl } [ ";" ]] "}"
func StructType(ss []State) []State {
	ss = tokenParserState(ss, token.STRUCT)
	ss = tokenReaderState(ss, token.LBRACE)
	fields := FieldDecl(ss)
	next := fields
	for len(next) != 0 {
		field := tokenReaderState(next, token.SEMICOLON)
		field = FieldDecl(field)
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

func SwitchStmtState(ss []State) []State {
	return append(ExprSwitchStmtState(ss), TypeSwitchStmtState(ss)...)
}

func TopLevelDecl(ts [][]*Token) [][]*Token {
	return append(append(Declaration(ts), FunctionDecl(ts)...), MethodDecl(ts)...)
}

func Type(ss []State) []State {
	paren := tokenReaderState(ss, token.LPAREN)
	if len(paren) != 0 {
		paren = Type(paren)
	}
	paren = tokenReaderState(paren, token.RPAREN)
	for i, p := range paren {
		thisParen := p.r[len(p.r)-1].(typ)
		thisParen.parens += 1
		paren[i].r = rAppend(p.r, 1, thisParen)
	}

	ss = append(TypeName(ss), TypeLit(ss)...)

	for i, s := range ss {
		t := typ{s.r[len(s.r)-1], 0}
		ss[i].r = rAppend(s.r, 1, t)
	}
	return append(ss, paren...)
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

func TypeAssertion(ss []State) []State {
	ss = tokenReaderState(ss, token.PERIOD)
	ss = tokenReaderState(ss, token.LPAREN)
	ss = Type(ss)
	ss = tokenReaderState(ss, token.RPAREN)
	for i, s := range ss {
		t := typeAssertion{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, t)
	}
	return ss
}

type typeAssertion struct{ r Renderer }

func (t typeAssertion) Render() []byte { return append(append([]byte(`.(`), t.r.Render()...), `)`...) }

func TypeCaseClause(ts [][]*Token) [][]*Token {
	ts = TypeSwitchCase(ts)
	ts = tokenReader(ts, token.COLON)
	return StatementList(ts)
}

func TypeCaseClauseState(ss []State) []State {
	ss = TypeSwitchCaseState(ss)
	ss = tokenReaderState(ss, token.COLON)
	ss = StatementListState(ss)
	for i, s := range ss {
		tcc := typeCaseClause{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, tcc)
	}
	return ss
}

type typeCaseClause struct{ switchCase, stmtList Renderer }

func (t typeCaseClause) Render() []byte {
	return append(append(t.switchCase.Render(), `:`...), t.stmtList.Render()...)
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

// bad spec
// "type" ( TypeSpec | "(" [ TypeSpec { ";" TypeSpec } [ ";" ]] ")" )
func TypeDeclState(ss []State) []State {
	ss = tokenParserState(ss, token.TYPE)

	bare := TypeSpecState(ss)
	for i, b := range bare {
		td := typeDecl{r: []Renderer{b.r[len(b.r)-1]}}
		bare[i].r = rAppend(b.r, 2, td)
	}

	ss = tokenReaderState(ss, token.LPAREN)

	loop := TypeSpecState(ss)
	ss = append(ss, loop...)

	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = TypeSpecState(loop)
		ss = append(ss, loop...)
	}
	ss = append(ss, tokenParserState(ss, token.SEMICOLON)...)
	ss = tokenReaderState(ss, token.RPAREN)

	for i, s := range ss {
		td := typeDecl{paren: true}
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.SEMICOLON {
			td.trailingSemi = true
			s.r = s.r[:len(s.r)-1]
		}

		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.TYPE {
				break
			}
			td.r = append(td.r, s.r[len(s.r)-1])
			s.r = s.r[:len(s.r)-1]
		}

		for i := 0; i < len(td.r)/2; i++ {
			td.r[i], td.r[len(td.r)-i-1] = td.r[len(td.r)-i-1], td.r[i]
		}

		ss[i].r = rAppend(s.r, 1, td)
	}
	return append(bare, ss...)
}

type typeDecl struct {
	r                   []Renderer
	paren, trailingSemi bool
}

func (t typeDecl) Render() []byte {
	ret := []byte(`type `)
	if !t.paren {
		return append(ret, t.r[0].Render()...)
	}

	ret = append(ret, `(`...)
	if len(t.r) == 0 {
		return append(ret, `)`...)
	}
	ret = append(ret, t.r[0].Render()...)
	for i := 1; i < len(t.r); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, t.r[i].Render()...)
	}
	if t.trailingSemi {
		ret = append(ret, `;`...)
	}
	return append(ret, `)`...)
}

func TypeDefState(ss []State) []State {
	ss = tokenParserState(ss, token.IDENT)
	ss = Type(ss)
	for i, s := range ss {
		td := typeDef{s.r[len(s.r)-2], s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, td)
	}
	return ss
}

type typeDef struct{ id, typ Renderer }

func (t typeDef) Render() []byte {
	return append(append(t.id.Render(), ` `...), t.typ.Render()...)
}

func TypeList(ts [][]*Token) [][]*Token {
	ts = fromState(Type(toState(ts)))
	next := ts
	for len(next) != 0 {
		current := tokenReader(next, token.COMMA)
		current = fromState(Type(toState(current)))
		ts = append(ts, current...)
		next = current
	}
	return ts
}

func TypeListState(ss []State) []State {
	ss = Type(ss)
	for i, s := range ss {
		tl := typeList{s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, tl)
	}
	loop := ss
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.COMMA)
		loop = Type(loop)
		for i, l := range loop {
			tl := append(l.r[len(l.r)-2].(typeList), l.r[len(l.r)-1])
			loop[i].r = rAppend(l.r, 2, tl)
		}
		ss = append(ss, loop...)
	}
	return ss
}

type typeList []Renderer

func (t typeList) Render() []byte {
	if len(t) == 0 {
		return nil
	}
	ret := t[0].Render()
	for i := 1; i < len(t); i++ {
		ret = append(ret, `,`...)
		ret = append(ret, t[i].Render()...)
	}
	return ret
}

func TypeLit(ss []State) []State {
	return append(
		append(append(ArrayType(ss), StructType(ss)...),
			append(PointerType(ss), FunctionType(ss)...)...),
		append(append(InterfaceType(ss), SliceType(ss)...),
			append(MapType(ss), ChannelType(ss)...)...)...)
}

func TypeName(ss []State) []State {
	return append(QualifiedIdent(ss), tokenParserState(ss, token.IDENT)...)
}

func TypeSpec(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.IDENT)
	return fromState(Type(toState(ts)))
}

func TypeSpecState(ss []State) []State { return append(AliasDeclState(ss), TypeDefState(ss)...) }

func TypeSwitchCase(ts [][]*Token) [][]*Token {
	cas := tokenReader(ts, token.CASE)
	return append(TypeList(cas), tokenReader(ts, token.DEFAULT)...)
}

func TypeSwitchCaseState(ss []State) []State {
	cas := tokenReaderState(ss, token.CASE)
	cas = TypeListState(cas)
	for i, c := range cas {
		tsc := typeSwitchCase{c.r[len(c.r)-1]}
		cas[i].r = rAppend(c.r, 1, tsc)
	}

	ss = tokenReaderState(ss, token.DEFAULT)
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 0, typeSwitchCase{})
	}
	return append(ss, cas...)
}

type typeSwitchCase struct{ r Renderer }

func (t typeSwitchCase) Render() []byte {
	if t.r == nil {
		return []byte(`default`)
	}
	return append([]byte(`case `), t.r.Render()...)
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

func TypeSwitchGuardState(ss []State) []State {
	id := tokenParserState(ss, token.IDENT)
	id = tokenReaderState(id, token.DEFINE)
	id = PrimaryExprState(id)
	id = tokenReaderState(id, token.PERIOD)
	id = tokenReaderState(id, token.LPAREN)
	id = tokenReaderState(id, token.TYPE)
	id = tokenReaderState(id, token.RPAREN)
	for i, elem := range id {
		tsg := typeSwitchGuard{elem.r[len(elem.r)-2], elem.r[len(elem.r)-1]}
		id[i].r = rAppend(elem.r, 2, tsg)
	}

	ss = PrimaryExprState(ss)
	ss = tokenReaderState(ss, token.PERIOD)
	ss = tokenReaderState(ss, token.LPAREN)
	ss = tokenReaderState(ss, token.TYPE)
	ss = tokenReaderState(ss, token.RPAREN)
	for i, s := range ss {
		tsg := typeSwitchGuard{expr: s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 1, tsg)
	}

	return append(ss, id...)
}

type typeSwitchGuard struct{ id, expr Renderer }

func (t typeSwitchGuard) Render() []byte {
	var ret []byte
	if t.id != nil {
		ret = append(t.id.Render(), `:=`...)
	}
	ret = append(ret, t.expr.Render()...)
	return append(ret, `.(type)`...)
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

// bad spec
// "switch" [ SimpleStmt ";" ] TypeSwitchGuard "{" [ TypeCaseClause { ";" TypeCaseClause } [ ";" ]] "}"
func TypeSwitchStmtState(ss []State) []State {
	ss = tokenParserState(ss, token.SWITCH)
	simple := SimpleStmtState(ss)
	simple = tokenReaderState(simple, token.SEMICOLON)
	for i, s := range simple {
		var tss typeSwitchStmt
		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.SWITCH {
			tss.stmt = e{}
		} else {
			tss.stmt = s.r[len(s.r)-1]
			s.r = s.r[:len(s.r)-1]
		}
		simple[i].r = rAppend(s.r, 1, tss)
	}
	for i, s := range ss {
		ss[i].r = rAppend(s.r, 1, typeSwitchStmt{})
	}
	ss = append(ss, simple...)
	ss = TypeSwitchGuardState(ss)
	for i, s := range ss {
		tss := s.r[len(s.r)-2].(typeSwitchStmt)
		tss.guard = s.r[len(s.r)-1]
		ss[i].r = rAppend(s.r, 2, tss)
	}
	ss = tokenReaderState(ss, token.LBRACE)
	loop := TypeCaseClauseState(ss)
	for i, l := range loop {
		tss := l.r[len(l.r)-2].(typeSwitchStmt)
		tss.clauses = append(tss.clauses, l.r[len(l.r)-1])
		loop[i].r = rAppend(l.r, 2, tss)
	}
	ss = append(ss, loop...)
	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = TypeCaseClauseState(loop)
		for i, l := range loop {
			tss := l.r[len(l.r)-2].(typeSwitchStmt)
			tss.clauses = append(tss.clauses, l.r[len(l.r)-1])
			loop[i].r = rAppend(l.r, 2, tss)
		}
		ss = append(ss, loop...)
	}
	return tokenReaderState(ss, token.RBRACE)
}

type typeSwitchStmt struct {
	stmt, guard Renderer
	clauses     []Renderer
}

func (t typeSwitchStmt) Render() []byte {
	ret := []byte(`switch `)
	if t.stmt != nil {
		ret = append(ret, t.stmt.Render()...)
		ret = append(ret, `;`...)
	}
	ret = append(ret, t.guard.Render()...)
	ret = append(ret, `{`...)
	if len(t.clauses) == 0 {
		return append(ret, `}`...)
	}
	ret = append(ret, t.clauses[0].Render()...)
	for i := 1; i < len(t.clauses); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, t.clauses[i].Render()...)
	}
	return append(ret, `}`...)
}

func UnaryExpr(ts [][]*Token) [][]*Token {
	uo := fromState(UnaryOp(toState(ts)))
	if len(uo) == 0 {
		return PrimaryExpr(ts)
	}
	return append(PrimaryExpr(ts), UnaryExpr(uo)...)
}

func UnaryExprState(ss []State) []State {
	uo := UnaryOp(ss)
	if len(uo) != 0 {
		uo = UnaryExprState(uo)
	}
	for i, u := range uo {
		ue := u.r[len(u.r)-1].(unaryExpr)
		ue = append(ue, u.r[len(u.r)-2])
		uo[i].r = rAppend(u.r, 2, ue)
	}

	pe := PrimaryExprState(ss)
	for i, p := range pe {
		ue := unaryExpr{p.r[len(p.r)-1]}
		pe[i].r = rAppend(p.r, 1, ue)
	}
	return append(pe, uo...)
}

// slice is in reverse order
type unaryExpr []Renderer

func (u unaryExpr) Render() []byte {
	var ret []byte
	for i := len(u) - 1; i >= 0; i-- {
		ret = append(ret, u[i].Render()...)
	}
	return ret
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

func VarDeclState(ss []State) []State {
	ss = tokenParserState(ss, token.VAR)
	bare := VarSpecState(ss)
	for i, b := range bare {
		vd := varDecl{[]Renderer{b.r[len(b.r)-1]}, false, false}
		bare[i].r = rAppend(b.r, 2, vd)
	}

	ss = tokenReaderState(ss, token.LPAREN)
	loop := VarSpecState(ss)
	ss = append(ss, loop...)

	for len(loop) != 0 {
		loop = tokenReaderState(loop, token.SEMICOLON)
		loop = VarSpecState(loop)
		ss = append(ss, loop...)
	}
	ss = append(ss, tokenParserState(ss, token.SEMICOLON)...)
	ss = tokenReaderState(ss, token.RPAREN)
	for i, s := range ss {
		vd := varDecl{parens: true}

		if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.SEMICOLON {
			vd.trailingSemi = true
			s.r = s.r[:len(s.r)-1]
		}
		for {
			if tok, ok := s.r[len(s.r)-1].(*Token); ok && tok.tok == token.VAR {
				break
			}
			vd.r = append(vd.r, s.r[len(s.r)-1])
			s.r = s.r[:len(s.r)-1]
		}

		for i := 0; i < len(vd.r)/2; i++ {
			vd.r[i], vd.r[len(vd.r)-i-1] = vd.r[len(vd.r)-i-1], vd.r[i]
		}

		ss[i].r = rAppend(s.r, 1, vd)
	}

	return append(bare, ss...)
}

type varDecl struct {
	r                    []Renderer
	parens, trailingSemi bool
}

func (v varDecl) Render() []byte {
	ret := []byte(`var `)
	if !v.parens {
		return append(ret, v.r[0].Render()...)
	}

	ret = append(ret, `(`...)
	if len(v.r) == 0 {
		return append(ret, `)`...)
	}
	ret = append(ret, v.r[0].Render()...)
	for i := 1; i < len(v.r); i++ {
		ret = append(ret, `;`...)
		ret = append(ret, v.r[i].Render()...)
	}
	if v.trailingSemi {
		ret = append(ret, `;`...)
	}
	return append(ret, `)`...)
}

func VarSpec(ts [][]*Token) [][]*Token {
	ts = fromState(IdentifierList(toState(ts)))
	typ := fromState(Type(toState(ts)))
	extra := tokenReader(typ, token.ASSIGN)
	extra = ExpressionList(extra)
	typ = append(typ, extra...)
	assign := tokenReader(ts, token.ASSIGN)
	assign = ExpressionList(assign)
	return append(typ, assign...)
}

func VarSpecState(ss []State) []State {
	ss = IdentifierList(ss)

	typ := Type(ss)
	typExpr := tokenReaderState(typ, token.ASSIGN)
	typExpr = ExpressionListState(typExpr)

	for i, t := range typ {
		vs := varSpec{idList: t.r[len(t.r)-2], typ: t.r[len(t.r)-1]}
		typ[i].r = rAppend(t.r, 2, vs)
	}
	for i, t := range typExpr {
		vs := varSpec{t.r[len(t.r)-3], t.r[len(t.r)-2], t.r[len(t.r)-1]}
		typExpr[i].r = rAppend(t.r, 3, vs)
	}

	ss = tokenReaderState(ss, token.ASSIGN)
	ss = ExpressionListState(ss)
	for i, s := range ss {
		vs := varSpec{idList: s.r[len(s.r)-2], exprList: s.r[len(s.r)-1]}
		ss[i].r = rAppend(s.r, 2, vs)
	}

	return append(append(typ, typExpr...), ss...)
}

type varSpec struct{ idList, typ, exprList Renderer }

func (v varSpec) Render() []byte {
	ret := v.idList.Render()
	if v.typ != nil {
		ret = append(append(ret, ` `...), v.typ.Render()...)
	}
	if v.exprList != nil {
		ret = append(append(ret, `=`...), v.exprList.Render()...)
	}
	return ret
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
