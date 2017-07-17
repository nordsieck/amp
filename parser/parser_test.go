package parser

import (
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

type Tmap map[string][][]*Token
type Output struct {
	tree []string
	rest [][]*Token
}
type Omap map[string]Output

func GetStateOutput(ss []State) []StateOutput {
	var result []StateOutput
	for _, s := range ss {
		var r []string
		for _, rend := range s.r {
			r = append(r, string(rend.Render()))
		}
		result = append(result, StateOutput{r, s.t})
	}
	return result
}

type StateOutput struct {
	s []string
	t []*Token
}

var (
	empty  = [][]*Token(nil)
	ret    = &Token{token.SEMICOLON, "\n"}
	semi   = &Token{token.SEMICOLON, `;`}
	a      = &Token{token.IDENT, `a`}
	b      = &Token{token.IDENT, `b`}
	c      = &Token{token.IDENT, `c`}
	d      = &Token{token.IDENT, `d`}
	_int   = &Token{token.IDENT, `int`}
	zero   = &Token{token.INT, `0`}
	one    = &Token{token.INT, `1`}
	two    = &Token{token.INT, `2`}
	dot    = &Token{tok: token.PERIOD}
	comma  = &Token{tok: token.COMMA}
	lparen = &Token{tok: token.LPAREN}
	rparen = &Token{tok: token.RPAREN}
	_if    = &Token{token.IF, `if`}
	_else  = &Token{token.ELSE, `else`}
	rbrace = &Token{tok: token.RBRACE}
	lbrace = &Token{tok: token.LBRACE}
	inc    = &Token{tok: token.INC}
)

func TestAddOp(t *testing.T) {
	resultState(t, AddOp, map[string][]StateOutput{
		`+`: {{[]string{``, `+`}, []*Token{}}},
		`-`: {{[]string{``, `-`}, []*Token{}}},
		`|`: {{[]string{``, `|`}, []*Token{}}},
		`&`: {{[]string{``, `&`}, []*Token{}}},
		`1`: nil,
	})
}

func TestAnonymousField(t *testing.T) {
	resultState(t, AnonymousField, map[string][]StateOutput{
		`a.a`:  {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`*a.a`: {{[]string{``, `*a.a`}, []*Token{ret}}, {[]string{``, `*a`}, []*Token{ret, a, dot}}},
	})
}

func TestAnonymousField_Render(t *testing.T) {
	defect.Equal(t, string(anonymousField{true, a}.Render()), `*a`)
	defect.Equal(t, string(anonymousField{false, &qualifiedIdent{a, b}}.Render()), `a.b`)
}

func TestArguments(t *testing.T) {
	remaining(t, Arguments, Tmap{
		`()`:        {{ret}},
		`(a)`:       {{ret}},
		`(a,a)`:     {{ret}},
		`(a,a,)`:    {{ret}},
		`(a,a...,)`: {{ret}},
	})
}

func TestArrayType(t *testing.T) {
	remaining(t, ArrayType, Tmap{
		`[1]int`: {{ret}},
		`[a]int`: {{ret}},
	})
}

func TestAssignment(t *testing.T) {
	remaining(t, Assignment, Tmap{`a = 1`: {{ret}}})
}

func TestAssignOp(t *testing.T) {
	resultState(t, AssignOp, map[string][]StateOutput{
		`+=`:  {{[]string{``, `+=`}, []*Token{}}},
		`-=`:  {{[]string{``, `-=`}, []*Token{}}},
		`*=`:  {{[]string{``, `*=`}, []*Token{}}},
		`/=`:  {{[]string{``, `/=`}, []*Token{}}},
		`%=`:  {{[]string{``, `%=`}, []*Token{}}},
		`&=`:  {{[]string{``, `&=`}, []*Token{}}},
		`|=`:  {{[]string{``, `|=`}, []*Token{}}},
		`^=`:  {{[]string{``, `^=`}, []*Token{}}},
		`<<=`: {{[]string{``, `<<=`}, []*Token{}}},
		`>>=`: {{[]string{``, `>>=`}, []*Token{}}},
		`&^=`: {{[]string{``, `&^=`}, []*Token{}}},
		`=`:   {{[]string{``, `=`}, []*Token{}}},
	})
}

func TestBasicLit(t *testing.T) {
	resultState(t, BasicLit, map[string][]StateOutput{
		``:    nil,
		`1`:   {{[]string{``, `1`}, []*Token{ret}}},
		`1.1`: {{[]string{``, `1.1`}, []*Token{ret}}},
		`1i`:  {{[]string{``, `1i`}, []*Token{ret}}},
		`'a'`: {{[]string{``, `'a'`}, []*Token{ret}}},
		`"a"`: {{[]string{``, `"a"`}, []*Token{ret}}},
		`a`:   nil,
		`_`:   nil,
	})
}

func TestBinaryOp(t *testing.T) {
	resultState(t, BinaryOp, map[string][]StateOutput{
		`==`: {{[]string{``, `==`}, []*Token{}}},
		`+`:  {{[]string{``, `+`}, []*Token{}}},
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`||`: {{[]string{``, `||`}, []*Token{}}},
		`&&`: {{[]string{``, `&&`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestBlock(t *testing.T) {
	remaining(t, Block, Tmap{
		`{}`:        {{ret}},
		`{a()}`:     {{ret}},
		`{a();}`:    {{ret}},
		`{a();b()}`: {{ret}},
	})
}

func TestBreakStmt(t *testing.T) {
	remaining(t, BreakStmt, Tmap{
		`break a`: {{ret, a}, {ret}},
		`a`:       empty,
	})
}

func TestChannelType(t *testing.T) {
	remaining(t, ChannelType, Tmap{
		`chan int`:   {{ret}},
		`<-chan int`: {{ret}},
		`chan<- int`: {{ret}},
		`int`:        empty,
	})
}

func TestChannelTypeState(t *testing.T) {
	resultState(t, ChannelTypeState, map[string][]StateOutput{
		`chan int`:   {{[]string{``, `chan int`}, []*Token{ret}}},
		`<-chan int`: {{[]string{``, `<-chan int`}, []*Token{ret}}},
		`chan<- int`: {{[]string{``, `chan<- int`}, []*Token{ret}}},
		`int`:        nil,
	})
}

func TestChannelType_Render(t *testing.T) {
	defect.Equal(t, string(channelType{nil, &Token{token.IDENT, `int`}}.Render()), `chan int`)
	defect.Equal(t, string(channelType{&_true, &Token{token.IDENT, `int`}}.Render()), `<-chan int`)
	defect.Equal(t, string(channelType{&_false, &Token{token.IDENT, `int`}}.Render()), `chan<- int`)
}

func TestCompositeLit(t *testing.T) {
	remaining(t, CompositeLit, Tmap{
		`T{1}`: {{ret}},
		`T{foo: "bar", baz: "quux",}`: {{ret}, {ret}, {ret}, {ret}},
	})
}

func TestCommCase(t *testing.T) {
	remaining(t, CommCase, Tmap{
		`default`:   {{}},
		`case <-a`:  {{ret}},
		`case a<-5`: {{ret}, {ret, {token.INT, `5`}, {tok: token.ARROW}}},
	})
}

func TestCommClause(t *testing.T) {
	remaining(t, CommClause, Tmap{
		`default:`:     {{}},
		`case <-a:`:    {{}},
		`case a<-5:`:   {{}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case <-a: b(); c()`: {
			{ret, rparen, lparen, c, semi, rparen, lparen, b},
			{ret, rparen, lparen, c, semi, rparen, lparen},
			{ret, rparen, lparen, c, semi},
			{ret, rparen, lparen, c},
			{ret, rparen, lparen}, {ret}, {}},
	})
}

func TestConstDecl(t *testing.T) {
	remaining(t, ConstDecl, Tmap{
		`const a = 1`:                         {{ret}},
		`const a int = 1`:                     {{ret}},
		`const (a = 1)`:                       {{ret}},
		`const (a int = 1)`:                   {{ret}},
		`const (a, b int = 1, 2; c int = 3;)`: {{ret}},
	})
}

func TestConstSpec(t *testing.T) {
	remaining(t, ConstSpec, Tmap{
		`a = 1`:           {{ret}},
		`a int = 1`:       {{ret}},
		`a, b int = 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestContinueStmt(t *testing.T) {
	remaining(t, ContinueStmt, Tmap{
		`continue a`: {{ret, a}, {ret}},
		`a`:          empty,
	})
}

func TestConversion(t *testing.T) {
	remaining(t, Conversion, Tmap{
		`float(1)`:   {{ret}},
		`(int)(5,)`:  {{ret}},
		`a.a("foo")`: {{ret}},
	})
}

func TestDeclaration(t *testing.T) {
	remaining(t, Declaration, Tmap{
		`const a = 1`: {{ret}},
		`type a int`:  {{ret}},
		`var a int`:   {{ret}},
	})
}

func TestDeferStmt(t *testing.T) {
	remaining(t, DeferStmt, Tmap{
		`defer a()`: {{ret, rparen, lparen}, {ret}},
		`defer`:     empty,
		`a()`:       empty,
	})
}

func TestElement(t *testing.T) {
	remaining(t, Element, Tmap{
		`1+1`:   {{ret, one, {tok: token.ADD}}, {ret}},
		`{1+1}`: {{ret}},
	})
}

func TestElementList(t *testing.T) {
	remaining(t, ElementList, Tmap{
		`1:1`: {{ret, one, {tok: token.COLON}}, {ret}},
		`1:1, 1:1`: {
			{ret, one, {tok: token.COLON}, one, comma, one, {tok: token.COLON}},
			{ret, one, {tok: token.COLON}, one, comma},
			{ret, one, {tok: token.COLON}}, {ret}},
	})
}

func TestEllipsisArrayType(t *testing.T) {
	remaining(t, EllipsisArrayType, Tmap{`[...]int`: {{ret}}})
}

func TestEmptyStmt(t *testing.T) {
	remaining(t, EmptyStmt, Tmap{`1`: {{ret, one}}})
}

func TestExprCaseClause(t *testing.T) {
	remaining(t, ExprCaseClause, Tmap{
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
	})
}

func TestExpression(t *testing.T) {
	remaining(t, Expression, Tmap{
		`1`:    {{ret}},
		`1+1`:  {{ret, one, {tok: token.ADD}}, {ret}},
		`1+-1`: {{ret, one, {tok: token.SUB}, {tok: token.ADD}}, {ret}},
	})
}

func TestExpressionList(t *testing.T) {
	remaining(t, ExpressionList, Tmap{
		`1`:   {{ret}},
		`1,1`: {{ret, one, comma}, {ret}},
	})
}

func TestExpressionStmt(t *testing.T) {
	remaining(t, ExpressionStmt, Tmap{`1`: {{ret}}})
}

func TestExprSwitchCase(t *testing.T) {
	remaining(t, ExprSwitchCase, Tmap{
		`case 1`:    {{ret}},
		`case 1, 1`: {{ret, one, comma}, {ret}},
		`default`:   {{}},
	})
}

func TestExprSwitchStmt(t *testing.T) {
	remaining(t, ExprSwitchStmt, Tmap{
		`switch{}`:                                               {{ret}},
		`switch{default:}`:                                       {{ret}},
		`switch a := 1; {}`:                                      {{ret}},
		`switch a {}`:                                            {{ret}},
		`switch a := 1; a {}`:                                    {{ret}},
		`switch a := 1; a {default:}`:                            {{ret}},
		`switch {case true:}`:                                    {{ret}},
		`switch {case true: a()}`:                                {{ret}},
		`switch{ case true: a(); case false: b(); default: c()}`: {{ret}},
	})
}

func TestFallthroughStmt(t *testing.T) {
	resultState(t, FallthroughStmt, map[string][]StateOutput{
		`fallthrough`: {{[]string{``, `fallthrough`}, []*Token{ret}}},
		`a`:           nil,
	})
}

func TestFieldDecl(t *testing.T) {
	remaining(t, FieldDecl, Tmap{
		`a,a int`:    {{ret, _int, a, comma}, {ret}},
		`*int`:       {{ret}},
		`*int "foo"`: {{ret, {token.STRING, `"foo"`}}, {ret}},
	})
}

func TestFieldDeclState(t *testing.T) {
	resultState(t, FieldDeclState, map[string][]StateOutput{
		`a,a int`: {{[]string{``, `a,a int`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, _int, a, comma}}},
		`a,a int "foo"`: {{[]string{``, `a,a int`}, []*Token{ret, {token.STRING, `"foo"`}}}, {[]string{``, `a,a int "foo"`}, []*Token{ret}},
			{[]string{``, `a`}, []*Token{ret, {token.STRING, `"foo"`}, _int, a, comma}}},
		`*int`:       {{[]string{``, `*int`}, []*Token{ret}}},
		`*int "foo"`: {{[]string{``, `*int`}, []*Token{ret, {token.STRING, `"foo"`}}}, {[]string{``, `*int "foo"`}, []*Token{ret}}},
	})
}

func TestFieldDecl_Render(t *testing.T) {
	defect.Equal(t, string(fieldDecl{anonField: &Token{token.IDENT, `int`}, tag: &Token{token.STRING, `a`}}.Render()), `int a`)
	defect.Equal(t, string(fieldDecl{idList: identifierList{&Token{token.IDENT, `a`}, &Token{token.IDENT, `b`}}, typ: &Token{token.IDENT, `int`}}.Render()), `a,b int`)
}

func TestForClause(t *testing.T) {
	remaining(t, ForClause, Tmap{
		`;;`:       {{}, {}, {}, {}},
		`a := 1;;`: {{}, {}},
		`; a < 5;`: {{}, {}, {}, {}},
		`;; a++`: {{ret, inc, a}, {ret, inc, a}, {ret, inc, a},
			{ret, inc, a}, {ret, inc}, {ret, inc}, {ret}, {ret}},
		`a := 1; a < 5; a++`: {{ret, inc, a}, {ret, inc, a}, {ret, inc}, {ret}},
	})
}

func TestForStmt(t *testing.T) {
	remaining(t, ForStmt, Tmap{
		`for {}`:                               {{ret}},
		`for true {}`:                          {{ret}},
		`for a := 0; a < 5; a++ {}`:            {{ret}},
		`for range a {}`:                       {{ret}},
		`for { a() }`:                          {{ret}},
		`for { a(); }`:                         {{ret}},
		`for { a(); b() }`:                     {{ret}},
		`for a := 0; a < 5; a++ { a(); b(); }`: {{ret}},
	})
}

func TestFunction(t *testing.T) {
	remaining(t, Function, Tmap{
		`(){}`:                 {{ret}},
		`()(){}`:               {{ret}},
		`(int)int{ return 0 }`: {{ret}},
		`(){ a() }`:            {{ret}},
	})
}

func TestFunctionDecl(t *testing.T) {
	remaining(t, FunctionDecl, Tmap{
		`func f(){}`:                        {{ret}},
		`func f()(){}`:                      {{ret}},
		`func f(int)int{ a(); return b() }`: {{ret}},
	})
}

func TestFunctionLit(t *testing.T) {
	remaining(t, FunctionLit, Tmap{
		`func(){}`: {{ret}},
	})
}

func TestFunctionType(t *testing.T) {
	remaining(t, FunctionType, Tmap{
		`func()`:             {{ret}},
		`func()()`:           {{ret, rparen, lparen}, {ret}},
		`func(int, int) int`: {{ret, _int}, {ret}},
	})
}

func TestGoStmt(t *testing.T) {
	remaining(t, GoStmt, Tmap{`go a()`: {{ret, rparen, lparen}, {ret}}})
}

func TestGotoStmt(t *testing.T) {
	remaining(t, GotoStmt, Tmap{`goto a`: {{ret}}, `goto`: empty, `a`: empty})
}

func TestIdentifierList(t *testing.T) {
	resultState(t, IdentifierList, map[string][]StateOutput{
		`a`:   {{[]string{``, `a`}, []*Token{ret}}},
		`a,b`: {{[]string{``, `a`}, []*Token{ret, b, comma}}, {[]string{``, `a,b`}, []*Token{ret}}},
		`1`:   nil,
		`_`:   {{[]string{``, `_`}, []*Token{ret}}},
	})
}

func TestIdentifierList_Render(t *testing.T) {
	defect.Equal(t, string(identifierList{a, one}.Render()), `a,1`)
}

func TestIfStmt(t *testing.T) {
	remaining(t, IfStmt, Tmap{
		`if a {}`:              {{ret}},
		`if a := false; a {}`:  {{ret}},
		`if a {} else {}`:      {{ret, rbrace, lbrace, _else}, {ret}},
		`if a {} else if b {}`: {{ret, rbrace, lbrace, b, _if, _else}, {ret}},
		`if a := false; a {} else if b {} else {}`: {{ret, rbrace, lbrace, _else, rbrace, lbrace, b, _if, _else},
			{ret, rbrace, lbrace, _else}, {ret}},
		`if a { b() }`: {{ret}},
	})
}

func TestImportDecl(t *testing.T) {
	remaining(t, ImportDecl, Tmap{
		`import "a"`:                     {{ret}},
		`import ()`:                      {{ret}},
		`import ("a")`:                   {{ret}},
		`import ("a";)`:                  {{ret}},
		`import ("a"; "b")`:              {{ret}},
		`import ("a";. "b";_ "c";d "d")`: {{ret}},
	})
}

func TestImportSpec(t *testing.T) {
	remaining(t, ImportSpec, Tmap{
		`"a"`:   {{ret}},
		`. "a"`: {{ret}},
		`_ "a"`: {{ret}},
		`a "a"`: {{ret}},
	})
}

func TestIncDecStmt(t *testing.T) {
	remaining(t, IncDecStmt, Tmap{
		`a++`: {{ret}},
		`a--`: {{ret}},
	})
}

func TestIndex(t *testing.T) {
	remaining(t, Index, Tmap{
		`[a]`: {{ret}},
		`[1]`: {{ret}},
	})
}

func TestInterfaceType(t *testing.T) {
	remaining(t, InterfaceType, Tmap{
		`interface{}`:       {{ret}},
		`interface{a}`:      {{ret}},
		`interface{a()}`:    {{ret}},
		`interface{a();a;}`: {{ret}},
	})
}

func TestKey(t *testing.T) {
	remaining(t, Key, Tmap{
		`a`:     {{ret}, {ret}},
		`1+a`:   {{ret, a, {tok: token.ADD}}, {ret}},
		`"foo"`: {{ret}},
	})
}

func TestKeyedElement(t *testing.T) {
	remaining(t, KeyedElement, Tmap{
		`1`:   {{ret}},
		`1:1`: {{ret, one, {tok: token.COLON}}, {ret}},
	})
}

func TestLabeledStmt(t *testing.T) {
	remaining(t, LabeledStmt, Tmap{
		`a: var b int`: {{ret}, {ret, _int, b, {token.VAR, `var`}}},
	})
}

func TestLiteral(t *testing.T) {
	remaining(t, Literal, Tmap{
		`1`:        {{ret}},
		`T{1}`:     {{ret}},
		`a`:        empty,
		`_`:        empty,
		`func(){}`: {{ret}},
	})
}

func TestLiteralType(t *testing.T) {
	remaining(t, LiteralType, Tmap{
		`struct{}`:    {{ret}},
		`[1]int`:      {{ret}},
		`[...]int`:    {{ret}},
		`[]int`:       {{ret}},
		`map[int]int`: {{ret}},
		`a.a`:         {{ret}, {ret, a, dot}},
	})
}

func TestLiteralValue(t *testing.T) {
	remaining(t, LiteralValue, Tmap{
		`{1}`:           {{ret}},
		`{0: 1, 1: 2,}`: {{ret}},
	})
}

func TestMapType(t *testing.T) {
	remaining(t, MapType, Tmap{`map[int]int`: {{ret}}})
}

func TestMapTypeState(t *testing.T) {
	resultState(t, MapTypeState, map[string][]StateOutput{
		`map[int]int`: {{[]string{``, `map[int]int`}, []*Token{ret}}},
		`int`:         nil,
	})
}

func TestMapType_Render(t *testing.T) {
	defect.Equal(t, string(mapType{&Token{token.IDENT, `a`}, &Token{token.IDENT, `b`}}.Render()), `map[a]b`)
}

func TestMethodDecl(t *testing.T) {
	remaining(t, MethodDecl, Tmap{
		`func (m M) f(){}`:                 {{ret}},
		`func (m M) f()(){}`:               {{ret}},
		`func (m M) f(int)int{ return 0 }`: {{ret}},
		`func (m *M) f() { a(); b(); }`:    {{ret}},
	})
}

func TestMethodExpr(t *testing.T) {
	remaining(t, MethodExpr, Tmap{
		`1`:        empty,
		`a.a`:      {{ret}},
		`a.a.a`:    {{ret}, {ret, a, dot}},
		`(a).a`:    {{ret}},
		`(a.a).a`:  {{ret}},
		`(*a.a).a`: {{ret}},
	})
}

func TestMethodSpec(t *testing.T) {
	remaining(t, MethodSpec, Tmap{
		`a()`:   {{ret}, {ret, rparen, lparen}},
		`a`:     {{ret}},
		`a.a()`: {{ret, rparen, lparen}, {ret, rparen, lparen, a, dot}},
	})
}

func TestMulOp(t *testing.T) {
	resultState(t, MulOp, map[string][]StateOutput{
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`/`:  {{[]string{``, `/`}, []*Token{}}},
		`%`:  {{[]string{``, `%`}, []*Token{}}},
		`<<`: {{[]string{``, `<<`}, []*Token{}}},
		`>>`: {{[]string{``, `>>`}, []*Token{}}},
		`&`:  {{[]string{``, `&`}, []*Token{}}},
		`&^`: {{[]string{``, `&^`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestOperand(t *testing.T) {
	remaining(t, Operand, Tmap{
		`1`:     {{ret}},
		`a.a`:   {{ret, a, dot}, {ret}, {ret}},
		`(a.a)`: {{ret}, {ret}, {ret}},
	})
}

func TestOperandName(t *testing.T) {
	remaining(t, OperandName, Tmap{
		`1`:   empty,
		`_`:   {{ret}},
		`a`:   {{ret}},
		`a.a`: {{ret, a, dot}, {ret}},
	})
}

func TestPackageName(t *testing.T) {
	resultState(t, PackageName, map[string][]StateOutput{
		`a`: {{[]string{``, `a`}, []*Token{ret}}},
		`1`: nil,
		`_`: nil,
	})
}

func TestPackageClause(t *testing.T) {
	remaining(t, PackageClause, Tmap{
		`package a`: {{ret}},
		`package _`: empty,
		`package`:   empty,
		`a`:         empty,
	})
}

func TestParameterDecl(t *testing.T) {
	remaining(t, ParameterDecl, Tmap{
		`int`:      {{ret}},
		`...int`:   {{ret}},
		`a, b int`: {{ret, _int, b, comma}, {ret}},
		`b... int`: {{ret, _int, {tok: token.ELLIPSIS}}, {ret}},
	})
}

func TestParameterList(t *testing.T) {
	remaining(t, ParameterList, Tmap{
		`int`:      {{ret}},
		`int, int`: {{ret, _int, comma}, {ret}},
		`a, b int, c, d int`: {
			{ret, _int, d, comma, c, comma, _int, b, comma},
			{ret, _int, d, comma, c, comma},
			{ret, _int, d, comma, c, comma, _int},
			{ret, _int, d, comma},
			{ret, _int, d, comma, c, comma},
			{ret}, {ret, _int},
			{ret, _int, d, comma},
			{ret}, {ret}, {ret, _int}, {ret},
		},
	})
}

func TestParameters(t *testing.T) {
	remaining(t, Parameters, Tmap{
		`(int)`:                {{ret}},
		`(int, int)`:           {{ret}},
		`(int, int,)`:          {{ret}},
		`(a, b int)`:           {{ret}, {ret}},
		`(a, b int, c, d int)`: {{ret}, {ret}, {ret}, {ret}},
	})
}

func TestPointerType(t *testing.T) {
	remaining(t, PointerType, Tmap{
		`*int`: {{ret}},
		`int`:  empty,
	})
}

func TestPointerTypeState(t *testing.T) {
	resultState(t, PointerTypeState, map[string][]StateOutput{
		`*int`: {{[]string{``, `*int`}, []*Token{ret}}},
		`int`:  nil,
	})
}

func TestPointerType_Renderer(t *testing.T) {
	defect.Equal(t, string(pointerType{&Token{tok: token.IDENT, lit: `a`}}.Render()), `*a`)
}

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, Tmap{
		`1`: {{ret}},
		`(a.a)("foo")`: {
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret, rparen, {token.STRING, `"foo"`}, lparen},
			{ret}, {ret}, {ret}, {ret}},
		`a.a`:     {{ret, a, dot}, {ret}, {ret}, {ret}},
		`a[1]`:    {{ret, {tok: token.RBRACK}, one, {tok: token.LBRACK}}, {ret}},
		`a[:]`:    {{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}}, {ret}},
		`a.(int)`: {{ret, rparen, _int, lparen, dot}, {ret}},
		`a(b...)`: {{ret, rparen, {tok: token.ELLIPSIS}, b, lparen}, {ret}},
		`a(b...)[:]`: {
			{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK},
				rparen, {tok: token.ELLIPSIS}, b, lparen},
			{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}},
			{ret}},
	})
}

func TestQualifiedIdent(t *testing.T) {
	resultState(t, QualifiedIdent, map[string][]StateOutput{
		`1`:   nil,
		`a`:   nil,
		`a.a`: {{[]string{``, `a.a`}, []*Token{ret}}},
		`_.a`: nil,
		`a._`: {{[]string{``, `a._`}, []*Token{ret}}},
		`_._`: nil,
	})
}

func TestRangeClause(t *testing.T) {
	remaining(t, RangeClause, Tmap{
		`range a`:              {{ret}},
		`a[0], a[1] = range b`: {{ret}},
		`k, v := range a`:      {{ret}},
	})
}

func TestRAppend(t *testing.T) {
	first := []Renderer{one, two, a, b, c, d}
	second := rAppend(first, 2, _int)
	defect.DeepEqual(t, second, []Renderer{one, two, a, b, _int})
	first[0] = zero
	defect.DeepEqual(t, second, []Renderer{one, two, a, b, _int})
}

func TestReceiverType(t *testing.T) {
	remaining(t, ReceiverType, Tmap{
		`1`:      empty,
		`a.a`:    {{ret}, {ret, a, dot}},
		`(a.a)`:  {{ret}},
		`(*a.a)`: {{ret}},
	})
}

func TestRecvStmt(t *testing.T) {
	remaining(t, RecvStmt, Tmap{
		`<-a`: {{ret}},
		`a[0], b = <-c`: {{ret, c, {tok: token.ARROW}, {tok: token.ASSIGN}, b, comma,
			{tok: token.RBRACK}, zero, {tok: token.LBRACK}},
			{ret, c, {tok: token.ARROW}, {tok: token.ASSIGN}, b, comma}, {ret}},
		`a, b := <-c`: {{ret, c, {tok: token.ARROW}, {tok: token.DEFINE}, b, comma}, {ret}},
		`a[0], b := <-c`: {{ret, c, {tok: token.ARROW}, {tok: token.DEFINE}, b, comma,
			{tok: token.RBRACK}, zero, {tok: token.LBRACK}},
			{ret, c, {tok: token.ARROW}, {tok: token.DEFINE}, b, comma}},
	})
}

func TestRelOp(t *testing.T) {
	resultState(t, RelOp, map[string][]StateOutput{
		`==`: {{[]string{``, `==`}, []*Token{}}},
		`!=`: {{[]string{``, `!=`}, []*Token{}}},
		`>`:  {{[]string{``, `>`}, []*Token{}}},
		`>=`: {{[]string{``, `>=`}, []*Token{}}},
		`<`:  {{[]string{``, `<`}, []*Token{}}},
		`<=`: {{[]string{``, `<=`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestResult(t *testing.T) {
	remaining(t, Result, Tmap{
		`int`:        {{ret}},
		`(int, int)`: {{ret}},
	})
}

func TestReturnStmt(t *testing.T) {
	remaining(t, ReturnStmt, Tmap{
		`return 1`:    {{ret, one}, {ret}},
		`return 1, 2`: {{ret, two, comma, one}, {ret, two, comma}, {ret}},
	})
}

func TestSelector(t *testing.T) {
	remaining(t, Selector, Tmap{
		`1`:  empty,
		`a`:  empty,
		`.a`: {{ret}},
	})
}

func TestSelectStmt(t *testing.T) {
	remaining(t, SelectStmt, Tmap{
		`select {}`:                                            {{ret}},
		`select {default:}`:                                    {{ret}},
		`select {default: a()}`:                                {{ret}},
		`select {default: a();}`:                               {{ret}, {ret}},
		`select {case <-a: ;default:}`:                         {{ret}},
		`select {case <-a: b(); case c<-d: e(); default: f()}`: {{ret}},
	})
}

func TestSendStmt(t *testing.T) {
	remaining(t, SendStmt, Tmap{`a <- 1`: {{ret}}})
}

func TestShortVarDecl(t *testing.T) {
	remaining(t, ShortVarDecl, Tmap{
		`a := 1`:       {{ret}},
		`a, b := 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestSignature(t *testing.T) {
	remaining(t, Signature, Tmap{
		`()`:             {{ret}},
		`()()`:           {{ret, rparen, lparen}, {ret}},
		`(int, int) int`: {{ret, _int}, {ret}},
	})
}

func TestSimpleStmt(t *testing.T) {
	remaining(t, SimpleStmt, Tmap{
		`1`:      {{ret, one}, {ret}},
		`a <- 1`: {{ret, one, {tok: token.ARROW}, a}, {ret, one, {tok: token.ARROW}}, {ret}},
		`a++`:    {{ret, inc, a}, {ret, inc}, {ret}},
		`a = 1`:  {{ret, one, {tok: token.ASSIGN}, a}, {ret, one, {tok: token.ASSIGN}}, {ret}},
		`a := 1`: {{ret, one, {tok: token.DEFINE}, a}, {ret, one, {tok: token.DEFINE}}, {ret}},
	})
}

func TestSlice(t *testing.T) {
	remaining(t, Slice, Tmap{
		`[:]`:     {{ret}},
		`[1:]`:    {{ret}},
		`[:1]`:    {{ret}},
		`[1:1]`:   {{ret}},
		`[:1:1]`:  {{ret}},
		`[1:1:1]`: {{ret}},
	})
}

func TestSliceType(t *testing.T) {
	remaining(t, SliceType, Tmap{`[]int`: {{ret}}})
}

func TestSliceTypeState(t *testing.T) {
	resultState(t, SliceTypeState, map[string][]StateOutput{
		`[]int`: {{[]string{``, `[]int`}, []*Token{ret}}},
		`int`:   nil,
	})
}

func TestSliceType_Render(t *testing.T) {
	defect.Equal(t, string(sliceType{&Token{tok: token.IDENT, lit: `a`}}.Render()), `[]a`)
}

func TestSourceFile(t *testing.T) {
	remaining(t, SourceFile, Tmap{
		`package p`:             {{}},
		`package p; import "a"`: {{ret, {token.STRING, `"a"`}, {token.IMPORT, `import`}}, {}},
		`package p; var a int`:  {{ret, _int, a, {token.VAR, `var`}}, {}},
		`package p; import "a"; import "b"; var c int; var d int`: {
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}, semi, {token.STRING, `"b"`}, {token.IMPORT, `import`},
				semi, {token.STRING, `"a"`}, {token.IMPORT, `import`}},
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}, semi, {token.STRING, `"b"`}, {token.IMPORT, `import`}},
			{ret, _int, d, {token.VAR, `var`}, semi, _int, c, {token.VAR, `var`}},
			{ret, _int, d, {token.VAR, `var`}}, {}},
	})
}

func TestStatement(t *testing.T) {
	remaining(t, Statement, Tmap{
		`var a int`: {{ret}, {ret, _int, a, {token.VAR, `var`}}},
		`a: var b int`: {{ret},
			{ret, _int, b, {token.VAR, `var`}},
			{ret, _int, b, {token.VAR, `var`}, {tok: token.COLON}, a},
			{ret, _int, b, {token.VAR, `var`}, {tok: token.COLON}}},
		`a := 1`:      {{ret, one, {tok: token.DEFINE}, a}, {ret, one, {tok: token.DEFINE}}, {ret}},
		`go a()`:      {{ret, rparen, lparen, a, {token.GO, `go`}}, {ret, rparen, lparen}, {ret}},
		`return 1`:    {{ret, one, {token.RETURN, `return`}}, {ret, one}, {ret}},
		`break a`:     {{ret, a, {token.BREAK, `break`}}, {ret, a}, {ret}},
		`continue a`:  {{ret, a, {token.CONTINUE, `continue`}}, {ret, a}, {ret}},
		`goto a`:      {{ret, a, {token.GOTO, `goto`}}, {ret}},
		`fallthrough`: {{ret, {token.FALLTHROUGH, `fallthrough`}}, {ret}},
		`{a()}`:       {{ret, rbrace, rparen, lparen, a, lbrace}, {ret}},
		`if a {}`:     {{ret, rbrace, lbrace, a, _if}, {ret}},
		`switch {}`:   {{ret, rbrace, lbrace, {token.SWITCH, `switch`}}, {ret}},
		`select {}`:   {{ret, rbrace, lbrace, {token.SELECT, `select`}}, {ret}},
		`for {}`:      {{ret, rbrace, lbrace, {token.FOR, `for`}}, {ret}},
		`defer a()`:   {{ret, rparen, lparen, a, {token.DEFER, `defer`}}, {ret, rparen, lparen}, {ret}},
	})
}

func TestStatementList(t *testing.T) {
	remaining(t, StatementList, Tmap{
		`fallthrough`:  {{ret, {token.FALLTHROUGH, `fallthrough`}}, {ret}, {}},
		`fallthrough;`: {{semi, {token.FALLTHROUGH, `fallthrough`}}, {semi}, {}},
		`fallthrough; fallthrough`: {{ret, {token.FALLTHROUGH, `fallthrough`}, semi, {token.FALLTHROUGH, `fallthrough`}},
			{ret, {token.FALLTHROUGH, `fallthrough`}, semi},
			{ret, {token.FALLTHROUGH, `fallthrough`}}, {ret}, {}},
	})
}

func TestStructType(t *testing.T) {
	remaining(t, StructType, Tmap{
		`struct{}`:                      {{ret}},
		`struct{int}`:                   {{ret}},
		`struct{int;}`:                  {{ret}},
		`struct{int;float64;}`:          {{ret}},
		`struct{a int}`:                 {{ret}},
		`struct{a, b int}`:              {{ret}},
		`struct{a, b int;}`:             {{ret}},
		`struct{a, b int; c, d string}`: {{ret}},
	})
}

func TestStructTypeState(t *testing.T) {
	resultState(t, StructTypeState, map[string][]StateOutput{
		`struct{}`:                      {{[]string{``, `struct{}`}, []*Token{ret}}},
		`struct{int}`:                   {{[]string{``, `struct{int;}`}, []*Token{ret}}},
		`struct{int;}`:                  {{[]string{``, `struct{int;}`}, []*Token{ret}}},
		`struct{int;float64;}`:          {{[]string{``, `struct{int;float64;}`}, []*Token{ret}}},
		`struct{a int}`:                 {{[]string{``, `struct{a int;}`}, []*Token{ret}}},
		`struct{a, b int}`:              {{[]string{``, `struct{a,b int;}`}, []*Token{ret}}},
		`struct{a, b int;}`:             {{[]string{``, `struct{a,b int;}`}, []*Token{ret}}},
		`struct{a, b int; c, d string}`: {{[]string{``, `struct{a,b int;c,d string;}`}, []*Token{ret}}},
	})
}

func TestStructType_Render(t *testing.T) {
	defect.Equal(t, string(structType{&Token{tok: token.FALLTHROUGH}}.Render()), `struct{fallthrough;}`)
	defect.Equal(t, string(structType{&Token{tok: token.FALLTHROUGH}, &Token{tok: token.FALLTHROUGH}}.Render()), `struct{fallthrough;fallthrough;}`)
}

func TestSwitchStmt(t *testing.T) {
	remaining(t, SwitchStmt, Tmap{
		`switch {}`:          {{ret}},
		`switch a.(type) {}`: {{ret}},
	})
}

func TestTopLevelDecl(t *testing.T) {
	remaining(t, TopLevelDecl, Tmap{
		`var a int`:        {{ret}},
		`func f(){}`:       {{ret}},
		`func (m M) f(){}`: {{ret}},
	})
}

func TestType(t *testing.T) {
	remaining(t, Type, Tmap{
		`a`:        {{ret}},
		`a.a`:      {{ret}, {ret, a, dot}},
		`1`:        empty,
		`_`:        {{ret}},
		`(a.a)`:    {{ret}},
		`(((_)))`:  {{ret}},
		`chan int`: {{ret}},
	})
}

func TestTypeState(t *testing.T) {
	resultState(t, TypeState, map[string][]StateOutput{
		`a`:   {{[]string{``, `a`}, []*Token{ret}}},
		`a.a`: {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`1`:   nil,
		`_`:   {{[]string{``, `_`}, []*Token{ret}}},
		//`(a.a)`:    {{[]string{``, `(a.a)`}, []*Token{ret}}},
		//`(((_)))`:  {{[]string{``, `(((_)))`}, []*Token{ret}}},
		`chan int`: {{[]string{``, `chan int`}, []*Token{ret}}},
	})
}

func TestTypeAssertion(t *testing.T) {
	remaining(t, TypeAssertion, Tmap{
		`.(int)`: {{ret}},
		`1`:      empty,
	})
}

func TestTypeCaseClause(t *testing.T) {
	remaining(t, TypeCaseClause, Tmap{
		`case a:`:      {{}},
		`case a: b()`:  {{ret, rparen, lparen, b}, {ret, rparen, lparen}, {ret}, {}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case a, b: c(); d()`: {
			{ret, rparen, lparen, d, semi, rparen, lparen, c},
			{ret, rparen, lparen, d, semi, rparen, lparen},
			{ret, rparen, lparen, d, semi},
			{ret, rparen, lparen, d},
			{ret, rparen, lparen}, {ret}, {}},
	})
}

func TestTypeDecl(t *testing.T) {
	remaining(t, TypeDecl, Tmap{
		`type a int`:           {{ret}},
		`type (a int)`:         {{ret}},
		`type (a int;)`:        {{ret}},
		`type (a int; b int)`:  {{ret}},
		`type (a int; b int;)`: {{ret}},
	})
}

func TestTypeList(t *testing.T) {
	remaining(t, TypeList, Tmap{
		`a`:     {{ret}},
		`a, b`:  {{ret, b, comma}, {ret}},
		`a, b,`: {{comma, b, comma}, {comma}},
	})
}

func TestTypeLit(t *testing.T) {
	remaining(t, TypeLit, Tmap{
		`[1]int`:      {{ret}},
		`struct{}`:    {{ret}},
		`*int`:        {{ret}},
		`func()`:      {{ret}},
		`interface{}`: {{ret}},
		`[]int`:       {{ret}},
		`map[int]int`: {{ret}},
		`chan int`:    {{ret}},
	})
}

func TestTypeLitState(t *testing.T) {
	resultState(t, TypeLitState, map[string][]StateOutput{
		//`[1]int`:      {{[]string{``, `[1]int`}, []*Token{ret}}},
		`struct{}`: {{[]string{``, `struct{}`}, []*Token{ret}}},
		`*int`:     {{[]string{``, `*int`}, []*Token{ret}}},
		//`func()`:      {{[]string{``, `func()`}, []*Token{ret}}},
		//`interface{}`: {{[]string{``, `interface{}`}, []*Token{ret}}},
		`[]int`:       {{[]string{``, `[]int`}, []*Token{ret}}},
		`map[int]int`: {{[]string{``, `map[int]int`}, []*Token{ret}}},
		`chan int`:    {{[]string{``, `chan int`}, []*Token{ret}}},
	})
}

func TestTypeName(t *testing.T) {
	resultState(t, TypeName, map[string][]StateOutput{
		`a`:   {{[]string{``, `a`}, []*Token{ret}}},
		`a.a`: {{[]string{``, `a.a`}, []*Token{ret}}, {[]string{``, `a`}, []*Token{ret, a, dot}}},
		`1`:   nil,
		`_`:   {{[]string{``, `_`}, []*Token{ret}}},
	})
}

func TestTypeSpec(t *testing.T) {
	remaining(t, TypeSpec, Tmap{
		`a int`: {{ret}},
		`a`:     empty,
	})
}

func TestTypeSwitchCase(t *testing.T) {
	remaining(t, TypeSwitchCase, Tmap{
		`default`:   {{}},
		`case a`:    {{ret}},
		`case a, b`: {{ret, b, comma}, {ret}},
	})
}

func TestTypeSwitchGuard(t *testing.T) {
	remaining(t, TypeSwitchGuard, Tmap{
		`a.(type)`:      {{ret}},
		`a := a.(type)`: {{ret}},
	})
}

func TestTypeSwitchStmt(t *testing.T) {
	remaining(t, TypeSwitchStmt, Tmap{
		`switch a.(type) {}`:                              {{ret}},
		`switch a := 5; a := a.(type) {}`:                 {{ret}},
		`switch a.(type) { case int: }`:                   {{ret}},
		`switch a.(type) { case int:; }`:                  {{ret}, {ret}},
		`switch a.(type) { case int: b() }`:               {{ret}},
		`switch a.(type) { case int: b(); }`:              {{ret}, {ret}},
		`switch a.(type) { case int: b(); default: c() }`: {{ret}},
	})
}

func TestUnaryExpr(t *testing.T) {
	remaining(t, UnaryExpr, Tmap{
		`1`:  {{ret}},
		`-1`: {{ret}},
		`!a`: {{ret}},
	})
}

func TestUnaryOp(t *testing.T) {
	resultState(t, UnaryOp, map[string][]StateOutput{
		`+`:  {{[]string{``, `+`}, []*Token{}}},
		`-`:  {{[]string{``, `-`}, []*Token{}}},
		`!`:  {{[]string{``, `!`}, []*Token{}}},
		`^`:  {{[]string{``, `^`}, []*Token{}}},
		`*`:  {{[]string{``, `*`}, []*Token{}}},
		`&`:  {{[]string{``, `&`}, []*Token{}}},
		`<-`: {{[]string{``, `<-`}, []*Token{}}},
		`1`:  nil,
	})
}

func TestVarDecl(t *testing.T) {
	remaining(t, VarDecl, Tmap{
		`var a int`:                 {{ret}},
		`var (a int)`:               {{ret}},
		`var (a int;)`:              {{ret}},
		`var (a, b = 1, 2)`:         {{ret}},
		`var (a, b = 1, 2; c int;)`: {{ret}},
	})
}

func TestVarSpec(t *testing.T) {
	remaining(t, VarSpec, Tmap{
		`a int`:       {{ret}},
		`a int = 1`:   {{ret, one, {tok: token.ASSIGN}}, {ret}},
		`a = 1`:       {{ret}},
		`a, b = 1, 2`: {{ret, two, comma}, {ret}},
	})
}

func TestTokenReader(t *testing.T) {
	toks := [][]*Token{
		{ret, rparen},
		{ret, rparen, a, dot},
	}
	defect.DeepEqual(t, tokenReader(toks, token.RPAREN), [][]*Token{{ret}})
}

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{ret, rparen},
		{ret, rparen, a, dot},
	}
	tree, state := tokenParser(toks, token.RPAREN)
	defect.DeepEqual(t, state, [][]*Token{{ret}})
	defect.DeepEqual(t, tree, []Renderer{rparen})
}

func TestTokenReaderState(t *testing.T) {
	state := []State{
		{[]Renderer{e{}}, []*Token{ret, rparen}},
		{[]Renderer{e{}}, []*Token{ret, rparen, a, dot}},
	}
	newState := tokenReaderState(state, token.RPAREN)
	defect.DeepEqual(t, newState, []State{{[]Renderer{e{}}, []*Token{ret}}})
}

func TestTokenParserState(t *testing.T) {
	state := []State{
		{[]Renderer{e{}}, []*Token{ret, rparen}},
		{[]Renderer{e{}}, []*Token{ret, rparen, a, dot}},
	}
	newState := tokenParserState(state, token.RPAREN)
	defect.DeepEqual(t, newState, []State{{[]Renderer{e{}, rparen}, []*Token{ret}}})
}
