package parser

import (
	"fmt"
	"go/token"
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

func AnonymousField(ts [][]*Token) [][]*Token {
	ts = append(ts, tokenReader(ts, token.MUL)...)
	return fromState(TypeName(toState(ts)))
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
	_, ts = IdentifierList(ts)
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
	_, a := IdentifierList(ts)
	if len(a) != 0 {
		a = Type(a)
	}
	ts = append(AnonymousField(ts), a...)
	return append(ts, tokenReader(ts, token.STRING)...)
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

func GoStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.GO)
	return Expression(ts)
}

func GotoStmt(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.GOTO)
	return tokenReader(ts, token.IDENT)
}

func IdentifierList(ts [][]*Token) ([]Renderer, [][]*Token) {
	var result [][]*Token
	var trees []Renderer
	var tree identifierList
	for _, t := range ts {
		outT, outS := tokenParser([][]*Token{t}, token.IDENT)
		if len(outT) == 0 {
			continue
		}
		tree = identifierList{outT[0]}
		trees = append(trees, tree)
		t = outS[0]
		result = append(result, t)

		nextS := [][]*Token{t}
		for {
			_, outS = tokenParser(nextS, token.COMMA)
			outT, outS = tokenParser(outS, token.IDENT)
			if len(outT) == 0 {
				break
			}
			tree = append(tree, outT[0])
			trees = append(trees, tree)
			result = append(result, outS...)
			nextS = outS
		}
	}
	return trees, result
}

type identifierList []Renderer

func (il identifierList) Render() []byte {
	var result []byte
	for i := 0; i < len(il); i++ {
		result = append(result, il[i].Render()...)
		if i != len(il)-1 {
			result = append(result, `,`...)
		}
	}
	return result
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

func Key(ts [][]*Token) [][]*Token {
	return append(tokenReader(ts, token.IDENT), Expression(ts)...)
}

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

func LiteralType(ts [][]*Token) [][]*Token {
	tn := fromState(TypeName(toState(ts)))
	return append(
		append(
			append(StructType(ts), ArrayType(ts)...),
			append(EllipsisArrayType(ts), SliceType(ts)...)...),
		append(MapType(ts), tn...)...)
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
	_, idList := IdentifierList(ts)
	ts = append(ts, idList...)
	ts = append(ts, tokenReader(ts, token.ELLIPSIS)...)
	if len(ts) == 0 {
		return nil
	}
	return Type(ts)
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

func Parameters(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.LPAREN)
	params := ParameterList(ts)
	params = append(params, tokenReader(params, token.COMMA)...)
	ts = append(ts, params...)
	return tokenReader(ts, token.RPAREN)
}

func PointerType(ts [][]*Token) [][]*Token {
	ts = tokenReader(ts, token.MUL)
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

func QualifiedIdent(ss []State) []State {
	ss = PackageName(ss)
	ss = tokenParserState(ss, token.PERIOD)
	ss = tokenParserState(ss, token.IDENT)
	for i, s := range ss {
		qi := qualifiedIdent{s.r[len(s.r)-3], s.r[len(s.r)-1]}
		ss[i].r = append(s.r[:len(s.r)-3], qi)
	}
	return ss
}

// func QualifiedIdent(ts [][]*Token) ([]Renderer, [][]*Token) {
// 	var result [][]*Token
// 	var qi []Renderer
// 	for _, t := range ts {

// 		pkg, outS := PackageName([][]*Token{t})
// 		_, outS = tokenParser(outS, token.PERIOD)
// 		name, outS := tokenParser(outS, token.IDENT)
// 		if len(name) == 0 {
// 			continue
// 		}
// 		qi = append(qi, &qualifiedIdent{pkg[0], name[0]})
// 		result = append(result, outS...)
// 	}
// 	return qi, result
// }

type qualifiedIdent struct{ pkg, name Renderer }

func (q qualifiedIdent) Render() []byte {
	var ret []byte
	ret = append(ret, q.pkg.Render()...)
	ret = append(ret, "."...)
	return append(ret, q.name.Render()...)
}

func RangeClause(ts [][]*Token) [][]*Token {
	exp := ExpressionList(ts)
	exp = tokenReader(exp, token.ASSIGN)
	_, id := IdentifierList(ts)
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
	_, ident := IdentifierList(ts)
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
	_, ts = IdentifierList(ts)
	ts = tokenReader(ts, token.DEFINE)
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

func SwitchStmt(ts [][]*Token) [][]*Token { return append(ExprSwitchStmt(ts), TypeSwitchStmt(ts)...) }

func TopLevelDecl(ts [][]*Token) [][]*Token {
	return append(append(Declaration(ts), FunctionDecl(ts)...), MethodDecl(ts)...)
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
	_, ts = IdentifierList(ts)
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
	fmt.Println("-----")
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
