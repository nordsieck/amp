package parser

import (
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

type Tmap map[string][][]*Token
type Output struct {
	tree []interface{}
	rest [][]*Token
}
type Omap map[string]Output

var (
	empty  = [][]*Token(nil)
	ret    = &Token{tok: token.SEMICOLON, lit: "\n"}
	semi   = &Token{tok: token.SEMICOLON, lit: `;`}
	a      = &Token{tok: token.IDENT, lit: `a`}
	_int   = &Token{tok: token.IDENT, lit: `int`}
	one    = &Token{tok: token.INT, lit: `1`}
	dot    = &Token{tok: token.PERIOD}
	comma  = &Token{tok: token.COMMA}
	lparen = &Token{tok: token.LPAREN}
	rparen = &Token{tok: token.RPAREN}
	_if    = &Token{tok: token.IF, lit: `if`}
	_else  = &Token{tok: token.ELSE, lit: `else`}
	rbrace = &Token{tok: token.RBRACE}
	lbrace = &Token{tok: token.LBRACE}
)

func TestAddOp(t *testing.T) {
	result(t, AddOp, Omap{
		`+`: {[]interface{}{&Token{tok: token.ADD}}, [][]*Token{{}}},
		`-`: {[]interface{}{&Token{tok: token.SUB}}, [][]*Token{{}}},
		`|`: {[]interface{}{&Token{tok: token.OR}}, [][]*Token{{}}},
		`&`: {[]interface{}{&Token{tok: token.AND}}, [][]*Token{{}}},
		`1`: {},
	})
}

func TestAnonymousField(t *testing.T) {
	remaining(t, AnonymousField, Tmap{
		`a.a`:  {{ret}, {ret, a, dot}},
		`*a.a`: {{ret}, {ret, a, dot}},
	})
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
	result(t, AssignOp, Omap{
		`+=`:  {[]interface{}{&Token{tok: token.ADD_ASSIGN}}, [][]*Token{{}}},
		`-=`:  {[]interface{}{&Token{tok: token.SUB_ASSIGN}}, [][]*Token{{}}},
		`*=`:  {[]interface{}{&Token{tok: token.MUL_ASSIGN}}, [][]*Token{{}}},
		`/=`:  {[]interface{}{&Token{tok: token.QUO_ASSIGN}}, [][]*Token{{}}},
		`%=`:  {[]interface{}{&Token{tok: token.REM_ASSIGN}}, [][]*Token{{}}},
		`&=`:  {[]interface{}{&Token{tok: token.AND_ASSIGN}}, [][]*Token{{}}},
		`|=`:  {[]interface{}{&Token{tok: token.OR_ASSIGN}}, [][]*Token{{}}},
		`^=`:  {[]interface{}{&Token{tok: token.XOR_ASSIGN}}, [][]*Token{{}}},
		`<<=`: {[]interface{}{&Token{tok: token.SHL_ASSIGN}}, [][]*Token{{}}},
		`>>=`: {[]interface{}{&Token{tok: token.SHR_ASSIGN}}, [][]*Token{{}}},
		`&^=`: {[]interface{}{&Token{tok: token.AND_NOT_ASSIGN}}, [][]*Token{{}}},
		`=`:   {[]interface{}{&Token{tok: token.ASSIGN}}, [][]*Token{{}}},
	})
}

func TestBasicLit(t *testing.T) {
	result(t, BasicLit, Omap{
		``:    {},
		`1`:   {[]interface{}{one}, [][]*Token{{ret}}},
		`1.1`: {[]interface{}{&Token{tok: token.FLOAT, lit: `1.1`}}, [][]*Token{{ret}}},
		`1i`:  {[]interface{}{&Token{tok: token.IMAG, lit: `1i`}}, [][]*Token{{ret}}},
		`'a'`: {[]interface{}{&Token{tok: token.CHAR, lit: `'a'`}}, [][]*Token{{ret}}},
		`"a"`: {[]interface{}{&Token{tok: token.STRING, lit: `"a"`}}, [][]*Token{{ret}}},
		`a`:   {},
		`_`:   {},
	})
}

func TestBinaryOp(t *testing.T) {
	result(t, BinaryOp, Omap{
		`==`: {[]interface{}{&Token{tok: token.EQL}}, [][]*Token{{}}},
		`+`:  {[]interface{}{&Token{tok: token.ADD}}, [][]*Token{{}}},
		`*`:  {[]interface{}{&Token{tok: token.MUL}}, [][]*Token{{}}},
		`||`: {[]interface{}{&Token{tok: token.LOR}}, [][]*Token{{}}},
		`&&`: {[]interface{}{&Token{tok: token.LAND}}, [][]*Token{{}}},
		`1`:  {},
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
		`case a<-5`: {{ret}, {ret, {tok: token.INT, lit: `5`}, {tok: token.ARROW}}},
	})
}

func TestCommClause(t *testing.T) {
	remaining(t, CommClause, Tmap{
		`default:`:     {{}},
		`case <-a:`:    {{}},
		`case a<-5:`:   {{}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case <-a: b(); c()`: {
			{ret, rparen, lparen, {tok: token.IDENT, lit: `c`}, semi, rparen, lparen, {tok: token.IDENT, lit: `b`}},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `c`}, semi, rparen, lparen},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `c`}, semi},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `c`}},
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
		`a, b int = 1, 2`: {{ret, {tok: token.INT, lit: `2`}, comma}, {ret}},
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
	result(t, FallthroughStmt, Omap{
		`fallthrough`: {[]interface{}{&Token{tok: token.FALLTHROUGH, lit: `fallthrough`}}, [][]*Token{{ret}}},
		`a`:           {},
	})
}

func TestFieldDecl(t *testing.T) {
	remaining(t, FieldDecl, Tmap{
		`a,a int`:    {{ret, _int, a, comma}, {ret}},
		`*int`:       {{ret}},
		`*int "foo"`: {{ret, {tok: token.STRING, lit: `"foo"`}}, {ret}},
	})
}

func TestForClause(t *testing.T) {
	remaining(t, ForClause, Tmap{
		`;;`:       {{}, {}, {}, {}},
		`a := 1;;`: {{}, {}},
		`; a < 5;`: {{}, {}, {}, {}},
		`;; a++`: {{ret, {tok: token.INC}, a}, {ret, {tok: token.INC}, a}, {ret, {tok: token.INC}, a},
			{ret, {tok: token.INC}, a}, {ret, {tok: token.INC}}, {ret, {tok: token.INC}}, {ret}, {ret}},
		`a := 1; a < 5; a++`: {{ret, {tok: token.INC}, a}, {ret, {tok: token.INC}, a}, {ret, {tok: token.INC}}, {ret}},
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
	result(t, IdentifierList, Omap{
		`a`:   {[]interface{}{[]*Token{a}}, [][]*Token{{ret}}},
		`a,a`: {[]interface{}{[]*Token{a}, []*Token{a, a}}, [][]*Token{{ret, a, comma}, {ret}}},
		`1`:   {},
		`_`:   {[]interface{}{[]*Token{{tok: token.IDENT, lit: `_`}}}, [][]*Token{{ret}}},
	})
}

func TestIdentifierList_Render(t *testing.T) {
	defect.DeepEqual(t, identifierList{a, one}.Render(), []byte(`a,1`))
}

func TestIfStmt(t *testing.T) {
	remaining(t, IfStmt, Tmap{
		`if a {}`:              {{ret}},
		`if a := false; a {}`:  {{ret}},
		`if a {} else {}`:      {{ret, rbrace, lbrace, _else}, {ret}},
		`if a {} else if b {}`: {{ret, rbrace, lbrace, {tok: token.IDENT, lit: `b`}, _if, _else}, {ret}},
		`if a := false; a {} else if b {} else {}`: {{ret, rbrace, lbrace, _else, rbrace, lbrace, {tok: token.IDENT, lit: `b`}, _if, _else},
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
		`a: var b int`: {{ret}, {ret, _int, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}}},
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
		`a()`: {{ret}, {ret, rparen, lparen}},
		`a`:   {{ret}},
		`a.a()`: {{ret, rparen, lparen},
			{ret, rparen, lparen, a, dot}},
	})
}

func TestMulOp(t *testing.T) {
	result(t, MulOp, Omap{
		`*`:  {[]interface{}{&Token{tok: token.MUL}}, [][]*Token{{}}},
		`/`:  {[]interface{}{&Token{tok: token.QUO}}, [][]*Token{{}}},
		`%`:  {[]interface{}{&Token{tok: token.REM}}, [][]*Token{{}}},
		`<<`: {[]interface{}{&Token{tok: token.SHL}}, [][]*Token{{}}},
		`>>`: {[]interface{}{&Token{tok: token.SHR}}, [][]*Token{{}}},
		`&`:  {[]interface{}{&Token{tok: token.AND}}, [][]*Token{{}}},
		`&^`: {[]interface{}{&Token{tok: token.AND_NOT}}, [][]*Token{{}}},
		`1`:  {},
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
	remaining(t, PackageName, Tmap{
		`a`: {{ret}},
		`1`: empty,
		`_`: empty,
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
		`a, b int`: {{ret, _int, {tok: token.IDENT, lit: `b`}, comma}, {ret}},
		`b... int`: {{ret, _int, {tok: token.ELLIPSIS}}, {ret}},
	})
}

func TestParameterList(t *testing.T) {
	remaining(t, ParameterList, Tmap{
		`int`:      {{ret}},
		`int, int`: {{ret, _int, comma}, {ret}},
		`a, b int, c, d int`: {
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma, {tok: token.IDENT, lit: `c`}, comma,
				_int, {tok: token.IDENT, lit: `b`}, comma},
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma, {tok: token.IDENT, lit: `c`}, comma},
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma, {tok: token.IDENT, lit: `c`}, comma,
				_int},
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma},
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma, {tok: token.IDENT, lit: `c`}, comma},
			{ret}, {ret, _int},
			{ret, _int, {tok: token.IDENT, lit: `d`}, comma},
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

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, Tmap{
		`1`: {{ret}},
		`(a.a)("foo")`: {
			{ret, rparen, {tok: token.STRING, lit: `"foo"`}, lparen},
			{ret, rparen, {tok: token.STRING, lit: `"foo"`}, lparen},
			{ret, rparen, {tok: token.STRING, lit: `"foo"`}, lparen},
			{ret}, {ret}, {ret}, {ret}},
		`a.a`:     {{ret, a, dot}, {ret}, {ret}, {ret}},
		`a[1]`:    {{ret, {tok: token.RBRACK}, one, {tok: token.LBRACK}}, {ret}},
		`a[:]`:    {{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}}, {ret}},
		`a.(int)`: {{ret, rparen, _int, lparen, dot}, {ret}},
		`a(b...)`: {{ret, rparen, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, lparen}, {ret}},
		`a(b...)[:]`: {
			{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK},
				rparen, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, lparen},
			{ret, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}},
			{ret}},
	})
}

func TestQualifiedIdent(t *testing.T) {
	remaining(t, QualifiedIdent, Tmap{
		`1`:   empty,
		`a`:   empty,
		`a.a`: {{ret}},
		`_.a`: empty,
		`a._`: {{ret}},
		`_._`: empty,
	})
}

func TestRangeClause(t *testing.T) {
	remaining(t, RangeClause, Tmap{
		`range a`:              {{ret}},
		`a[0], a[1] = range b`: {{ret}},
		`k, v := range a`:      {{ret}},
	})
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
		`a[0], b = <-c`: {{ret, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.ASSIGN}, {tok: token.IDENT, lit: `b`}, comma,
			{tok: token.RBRACK}, {tok: token.INT, lit: `0`}, {tok: token.LBRACK}},
			{ret, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.ASSIGN}, {tok: token.IDENT, lit: `b`}, comma}, {ret}},
		`a, b := <-c`: {{ret, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, comma}, {ret}},
		`a[0], b := <-c`: {{ret, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, comma,
			{tok: token.RBRACK}, {tok: token.INT, lit: `0`}, {tok: token.LBRACK}},
			{ret, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, comma}},
	})
}

func TestRelOp(t *testing.T) {
	result(t, RelOp, Omap{
		`==`: {[]interface{}{&Token{tok: token.EQL}}, [][]*Token{{}}},
		`!=`: {[]interface{}{&Token{tok: token.NEQ}}, [][]*Token{{}}},
		`>`:  {[]interface{}{&Token{tok: token.GTR}}, [][]*Token{{}}},
		`>=`: {[]interface{}{&Token{tok: token.GEQ}}, [][]*Token{{}}},
		`<`:  {[]interface{}{&Token{tok: token.LSS}}, [][]*Token{{}}},
		`<=`: {[]interface{}{&Token{tok: token.LEQ}}, [][]*Token{{}}},
		`1`:  {},
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
		`return 1, 2`: {{ret, {tok: token.INT, lit: `2`}, comma, one}, {ret, {tok: token.INT, lit: `2`}, comma}, {ret}},
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
		`a, b := 1, 2`: {{ret, {tok: token.INT, lit: `2`}, comma}, {ret}},
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
		`a++`:    {{ret, {tok: token.INC}, a}, {ret, {tok: token.INC}}, {ret}},
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

func TestSourceFile(t *testing.T) {
	remaining(t, SourceFile, Tmap{
		`package p`:             {{}},
		`package p; import "a"`: {{ret, {tok: token.STRING, lit: `"a"`}, {tok: token.IMPORT, lit: `import`}}, {}},
		`package p; var a int`:  {{ret, _int, a, {tok: token.VAR, lit: `var`}}, {}},
		`package p; import "a"; import "b"; var c int; var d int`: {
			{ret, _int, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, semi, _int, {tok: token.IDENT, lit: `c`},
				{tok: token.VAR, lit: `var`}, semi, {tok: token.STRING, lit: `"b"`}, {tok: token.IMPORT, lit: `import`}, semi,
				{tok: token.STRING, lit: `"a"`}, {tok: token.IMPORT, lit: `import`}},
			{ret, _int, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, semi, _int, {tok: token.IDENT, lit: `c`},
				{tok: token.VAR, lit: `var`}, semi, {tok: token.STRING, lit: `"b"`}, {tok: token.IMPORT, lit: `import`}},
			{ret, _int, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, semi, _int, {tok: token.IDENT, lit: `c`},
				{tok: token.VAR, lit: `var`}},
			{ret, _int, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}}, {}},
	})
}

func TestStatement(t *testing.T) {
	remaining(t, Statement, Tmap{
		`var a int`: {{ret}, {ret, _int, a, {tok: token.VAR, lit: `var`}}},
		`a: var b int`: {{ret},
			{ret, _int, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}},
			{ret, _int, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}, {tok: token.COLON}, a},
			{ret, _int, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}, {tok: token.COLON}}},
		`a := 1`:      {{ret, one, {tok: token.DEFINE}, a}, {ret, one, {tok: token.DEFINE}}, {ret}},
		`go a()`:      {{ret, rparen, lparen, a, {tok: token.GO, lit: `go`}}, {ret, rparen, lparen}, {ret}},
		`return 1`:    {{ret, one, {tok: token.RETURN, lit: `return`}}, {ret, one}, {ret}},
		`break a`:     {{ret, a, {tok: token.BREAK, lit: `break`}}, {ret, a}, {ret}},
		`continue a`:  {{ret, a, {tok: token.CONTINUE, lit: `continue`}}, {ret, a}, {ret}},
		`goto a`:      {{ret, a, {tok: token.GOTO, lit: `goto`}}, {ret}},
		`fallthrough`: {{ret, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {ret}},
		`{a()}`:       {{ret, rbrace, rparen, lparen, a, lbrace}, {ret}},
		`if a {}`:     {{ret, rbrace, lbrace, a, _if}, {ret}},
		`switch {}`:   {{ret, rbrace, lbrace, {tok: token.SWITCH, lit: `switch`}}, {ret}},
		`select {}`:   {{ret, rbrace, lbrace, {tok: token.SELECT, lit: `select`}}, {ret}},
		`for {}`:      {{ret, rbrace, lbrace, {tok: token.FOR, lit: `for`}}, {ret}},
		`defer a()`:   {{ret, rparen, lparen, a, {tok: token.DEFER, lit: `defer`}}, {ret, rparen, lparen}, {ret}},
	})
}

func TestStatementList(t *testing.T) {
	remaining(t, StatementList, Tmap{
		`fallthrough`:  {{ret, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {ret}, {}},
		`fallthrough;`: {{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi}, {}},
		`fallthrough; fallthrough`: {{ret, {tok: token.FALLTHROUGH, lit: `fallthrough`}, semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{ret, {tok: token.FALLTHROUGH, lit: `fallthrough`}, semi},
			{ret, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {ret}, {}},
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

func TestTypeAssertion(t *testing.T) {
	remaining(t, TypeAssertion, Tmap{
		`.(int)`: {{ret}},
		`1`:      empty,
	})
}

func TestTypeCaseClause(t *testing.T) {
	remaining(t, TypeCaseClause, Tmap{
		`case a:`:      {{}},
		`case a: b()`:  {{ret, rparen, lparen, {tok: token.IDENT, lit: `b`}}, {ret, rparen, lparen}, {ret}, {}},
		`default: a()`: {{ret, rparen, lparen, a}, {ret, rparen, lparen}, {ret}, {}},
		`case a, b: c(); d()`: {
			{ret, rparen, lparen, {tok: token.IDENT, lit: `d`}, semi, rparen, lparen, {tok: token.IDENT, lit: `c`}},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `d`}, semi, rparen, lparen},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `d`}, semi},
			{ret, rparen, lparen, {tok: token.IDENT, lit: `d`}},
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
		`a, b`:  {{ret, {tok: token.IDENT, lit: `b`}, comma}, {ret}},
		`a, b,`: {{comma, {tok: token.IDENT, lit: `b`}, comma}, {comma}},
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

func TestTypeName(t *testing.T) {
	remaining(t, TypeName, Tmap{
		`a`:   {{ret}},
		`a.a`: {{ret}, {ret, a, dot}},
		`1`:   empty,
		`_`:   {{ret}},
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
		`case a, b`: {{ret, {tok: token.IDENT, lit: `b`}, comma}, {ret}},
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
	result(t, UnaryOp, Omap{
		`+`:  {[]interface{}{token.ADD}, [][]*Token{{}}},
		`-`:  {[]interface{}{token.SUB}, [][]*Token{{}}},
		`!`:  {[]interface{}{token.NOT}, [][]*Token{{}}},
		`^`:  {[]interface{}{token.XOR}, [][]*Token{{}}},
		`*`:  {[]interface{}{token.MUL}, [][]*Token{{}}},
		`&`:  {[]interface{}{token.AND}, [][]*Token{{}}},
		`<-`: {[]interface{}{token.ARROW}, [][]*Token{{}}},
		`1`:  {},
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
		`a, b = 1, 2`: {{ret, {tok: token.INT, lit: `2`}, comma}, {ret}},
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
	defect.DeepEqual(t, tree, []interface{}{rparen})
}
