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
	empty = [][]*Token(nil)
	semi  = &Token{tok: token.SEMICOLON, lit: "\n"}
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
		`a.a`:  {{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`*a.a`: {{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
	})
}

func TestArguments(t *testing.T) {
	remaining(t, Arguments, Tmap{
		`()`:        {{semi}},
		`(a)`:       {{semi}},
		`(a,a)`:     {{semi}},
		`(a,a,)`:    {{semi}},
		`(a,a...,)`: {{semi}},
	})
}

func TestArrayType(t *testing.T) {
	remaining(t, ArrayType, Tmap{
		`[1]int`: {{semi}},
		`[a]int`: {{semi}},
	})
}

func TestAssignment(t *testing.T) {
	remaining(t, Assignment, Tmap{`a = 1`: {{semi}}})
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
		`1`:   {[]interface{}{&Token{tok: token.INT, lit: `1`}}, [][]*Token{{semi}}},
		`1.1`: {[]interface{}{&Token{tok: token.FLOAT, lit: `1.1`}}, [][]*Token{{semi}}},
		`1i`:  {[]interface{}{&Token{tok: token.IMAG, lit: `1i`}}, [][]*Token{{semi}}},
		`'a'`: {[]interface{}{&Token{tok: token.CHAR, lit: `'a'`}}, [][]*Token{{semi}}},
		`"a"`: {[]interface{}{&Token{tok: token.STRING, lit: `"a"`}}, [][]*Token{{semi}}},
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
		`{}`:        {{semi}},
		`{a()}`:     {{semi}},
		`{a();}`:    {{semi}},
		`{a();b()}`: {{semi}},
	})
}

func TestBreakStmt(t *testing.T) {
	remaining(t, BreakStmt, Tmap{
		`break a`: {{semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`a`:       empty,
	})
}

func TestChannelType(t *testing.T) {
	remaining(t, ChannelType, Tmap{
		`chan int`:   {{semi}},
		`<-chan int`: {{semi}},
		`chan<- int`: {{semi}},
		`int`:        empty,
	})
}

func TestCompositeLit(t *testing.T) {
	remaining(t, CompositeLit, Tmap{
		`T{1}`: {{semi}},
		`T{foo: "bar", baz: "quux",}`: {{semi}, {semi}, {semi}, {semi}},
	})
}

func TestCommCase(t *testing.T) {
	remaining(t, CommCase, Tmap{
		`default`:   {{}},
		`case <-a`:  {{semi}},
		`case a<-5`: {{semi}, {semi, {tok: token.INT, lit: `5`}, {tok: token.ARROW}}},
	})
}

func TestCommClause(t *testing.T) {
	remaining(t, CommClause, Tmap{
		`default:`:   {{}},
		`case <-a:`:  {{}},
		`case a<-5:`: {{}},
		`default: a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}},
		`case <-a: b(); c()`: {
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `c`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `b`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `c`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.RPAREN}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `c`}, {tok: token.SEMICOLON, lit: `;`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `c`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {},
		},
	})
}

func TestConstDecl(t *testing.T) {
	remaining(t, ConstDecl, Tmap{
		`const a = 1`:                         {{semi}},
		`const a int = 1`:                     {{semi}},
		`const (a = 1)`:                       {{semi}},
		`const (a int = 1)`:                   {{semi}},
		`const (a, b int = 1, 2; c int = 3;)`: {{semi}},
	})
}

func TestConstSpec(t *testing.T) {
	remaining(t, ConstSpec, Tmap{
		`a = 1`:     {{semi}},
		`a int = 1`: {{semi}},
		`a, b int = 1, 2`: {
			{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi},
		},
	})
}

func TestContinueStmt(t *testing.T) {
	remaining(t, ContinueStmt, Tmap{
		`continue a`: {{semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`a`:          empty,
	})
}

func TestConversion(t *testing.T) {
	remaining(t, Conversion, Tmap{
		`float(1)`:   {{semi}},
		`(int)(5,)`:  {{semi}},
		`a.a("foo")`: {{semi}},
	})
}

func TestDeclaration(t *testing.T) {
	remaining(t, Declaration, Tmap{
		`const a = 1`: {{semi}},
		`type a int`:  {{semi}},
		`var a int`:   {{semi}},
	})
}

func TestDeferStmt(t *testing.T) {
	remaining(t, DeferStmt, Tmap{
		`defer a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`defer`:     empty,
		`a()`:       empty,
	})
}

func TestElement(t *testing.T) {
	remaining(t, Element, Tmap{
		`1+1`:   {{semi, {tok: token.INT, lit: `1`}, {tok: token.ADD}}, {semi}},
		`{1+1}`: {{semi}},
	})
}

func TestElementList(t *testing.T) {
	remaining(t, ElementList, Tmap{
		`1:1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}}, {semi}},
		`1:1, 1:1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}, {tok: token.INT, lit: `1`},
			{tok: token.COMMA}, {tok: token.INT, lit: `1`}, {tok: token.COLON}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}, {tok: token.INT, lit: `1`}, {tok: token.COMMA}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}},
			{semi}},
	})
}

func TestEllipsisArrayType(t *testing.T) {
	remaining(t, EllipsisArrayType, Tmap{`[...]int`: {{semi}}})
}

func TestEmptyStmt(t *testing.T) {
	remaining(t, EmptyStmt, Tmap{`1`: {{semi, {tok: token.INT, lit: `1`}}}})
}

func TestExprCaseClause(t *testing.T) {
	remaining(t, ExprCaseClause, Tmap{
		`default: a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}},
	})
}

func TestExpression(t *testing.T) {
	remaining(t, Expression, Tmap{
		`1`:    {{semi}},
		`1+1`:  {{semi, {tok: token.INT, lit: `1`}, {tok: token.ADD}}, {semi}},
		`1+-1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.SUB}, {tok: token.ADD}}, {semi}},
	})
}

func TestExpressionList(t *testing.T) {
	remaining(t, ExpressionList, Tmap{
		`1`:   {{semi}},
		`1,1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestExpressionStmt(t *testing.T) {
	remaining(t, ExpressionStmt, Tmap{`1`: {{semi}}})
}

func TestExprSwitchCase(t *testing.T) {
	remaining(t, ExprSwitchCase, Tmap{
		`case 1`:    {{semi}},
		`case 1, 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COMMA}}, {semi}},
		`default`:   {{}},
	})
}

func TestExprSwitchStmt(t *testing.T) {
	remaining(t, ExprSwitchStmt, Tmap{
		`switch{}`:                                               {{semi}},
		`switch{default:}`:                                       {{semi}},
		`switch a := 1; {}`:                                      {{semi}},
		`switch a {}`:                                            {{semi}},
		`switch a := 1; a {}`:                                    {{semi}},
		`switch a := 1; a {default:}`:                            {{semi}},
		`switch {case true:}`:                                    {{semi}},
		`switch {case true: a()}`:                                {{semi}},
		`switch{ case true: a(); case false: b(); default: c()}`: {{semi}},
	})
}

func TestFallthroughStmt(t *testing.T) {
	result(t, FallthroughStmt, Omap{
		`fallthrough`: {[]interface{}{&Token{tok: token.FALLTHROUGH, lit: `fallthrough`}}, [][]*Token{{semi}}},
		`a`:           {},
	})
}

func TestFieldDecl(t *testing.T) {
	remaining(t, FieldDecl, Tmap{
		`a,a int`:    {{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `a`}, {tok: token.COMMA}}, {semi}},
		`*int`:       {{semi}},
		`*int "foo"`: {{semi, {tok: token.STRING, lit: `"foo"`}}, {semi}},
	})
}

func TestForClause(t *testing.T) {
	remaining(t, ForClause, Tmap{
		`;;`:       {{}, {}, {}, {}},
		`a := 1;;`: {{}, {}},
		`; a < 5;`: {{}, {}, {}, {}},
		`;; a++`: {{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}}, {semi, {tok: token.INC}}, {semi}, {semi}},
		`a := 1; a < 5; a++`: {{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INC}}, {semi}},
	})
}

func TestForStmt(t *testing.T) {
	remaining(t, ForStmt, Tmap{
		`for {}`:                               {{semi}},
		`for true {}`:                          {{semi}},
		`for a := 0; a < 5; a++ {}`:            {{semi}},
		`for range a {}`:                       {{semi}},
		`for { a() }`:                          {{semi}},
		`for { a(); }`:                         {{semi}},
		`for { a(); b() }`:                     {{semi}},
		`for a := 0; a < 5; a++ { a(); b(); }`: {{semi}},
	})
}

func TestFunction(t *testing.T) {
	remaining(t, Function, Tmap{
		`(){}`:                 {{semi}},
		`()(){}`:               {{semi}},
		`(int)int{ return 0 }`: {{semi}},
		`(){ a() }`:            {{semi}},
	})
}

func TestFunctionDecl(t *testing.T) {
	remaining(t, FunctionDecl, Tmap{
		`func f(){}`:                        {{semi}},
		`func f()(){}`:                      {{semi}},
		`func f(int)int{ a(); return b() }`: {{semi}},
	})
}

func TestFunctionLit(t *testing.T) {
	remaining(t, FunctionLit, Tmap{
		`func(){}`: {{semi}},
	})
}

func TestFunctionType(t *testing.T) {
	remaining(t, FunctionType, Tmap{
		`func()`:             {{semi}},
		`func()()`:           {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`func(int, int) int`: {{semi, {tok: token.IDENT, lit: `int`}}, {semi}},
	})
}

func TestGoStmt(t *testing.T) {
	remaining(t, GoStmt, Tmap{`go a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}}})
}

func TestGotoStmt(t *testing.T) {
	remaining(t, GotoStmt, Tmap{`goto a`: {{semi}}, `goto`: empty, `a`: empty})
}

func TestIdentifierList(t *testing.T) {
	remaining(t, IdentifierList, Tmap{
		`a`:   {{semi}},
		`a,a`: {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.COMMA}}, {semi}},
		`1`:   empty,
		`_`:   {{semi}},
	})
}

func TestIfStmt(t *testing.T) {
	remaining(t, IfStmt, Tmap{
		`if a {}`:             {{semi}},
		`if a := false; a {}`: {{semi}},
		`if a {} else {}`:     {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}}, {semi}},
		`if a {} else if b {}`: {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`},
			{tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}}, {semi}},
		`if a := false; a {} else if b {} else {}`: {{semi, {tok: token.RBRACE}, {tok: token.LBRACE},
			{tok: token.ELSE, lit: `else`}, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`},
			{tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}}, {semi}},
		`if a { b() }`: {{semi}},
	})
}

func TestImportDecl(t *testing.T) {
	remaining(t, ImportDecl, Tmap{
		`import "a"`:                     {{semi}},
		`import ()`:                      {{semi}},
		`import ("a")`:                   {{semi}},
		`import ("a";)`:                  {{semi}},
		`import ("a"; "b")`:              {{semi}},
		`import ("a";. "b";_ "c";d "d")`: {{semi}},
	})
}

func TestImportSpec(t *testing.T) {
	remaining(t, ImportSpec, Tmap{
		`"a"`:   {{semi}},
		`. "a"`: {{semi}},
		`_ "a"`: {{semi}},
		`a "a"`: {{semi}},
	})
}

func TestIncDecStmt(t *testing.T) {
	remaining(t, IncDecStmt, Tmap{
		`a++`: {{semi}},
		`a--`: {{semi}},
	})
}

func TestIndex(t *testing.T) {
	remaining(t, Index, Tmap{
		`[a]`: {{semi}},
		`[1]`: {{semi}},
	})
}

func TestInterfaceType(t *testing.T) {
	remaining(t, InterfaceType, Tmap{
		`interface{}`:       {{semi}},
		`interface{a}`:      {{semi}},
		`interface{a()}`:    {{semi}},
		`interface{a();a;}`: {{semi}},
	})
}

func TestKey(t *testing.T) {
	remaining(t, Key, Tmap{
		`a`:     {{semi}, {semi}},
		`1+a`:   {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.ADD}}, {semi}},
		`"foo"`: {{semi}},
	})
}

func TestKeyedElement(t *testing.T) {
	remaining(t, KeyedElement, Tmap{
		`1`:   {{semi}},
		`1:1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}}, {semi}},
	})
}

func TestLabeledStmt(t *testing.T) {
	remaining(t, LabeledStmt, Tmap{
		`a: var b int`: {{semi}, {semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}}},
	})
}

func TestLiteral(t *testing.T) {
	remaining(t, Literal, Tmap{
		`1`:        {{semi}},
		`T{1}`:     {{semi}},
		`a`:        empty,
		`_`:        empty,
		`func(){}`: {{semi}},
	})
}

func TestLiteralType(t *testing.T) {
	remaining(t, LiteralType, Tmap{
		`struct{}`:    {{semi}},
		`[1]int`:      {{semi}},
		`[...]int`:    {{semi}},
		`[]int`:       {{semi}},
		`map[int]int`: {{semi}},
		`a.a`:         {{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
	})
}

func TestLiteralValue(t *testing.T) {
	remaining(t, LiteralValue, Tmap{
		`{1}`:           {{semi}},
		`{0: 1, 1: 2,}`: {{semi}},
	})
}

func TestMapType(t *testing.T) {
	remaining(t, MapType, Tmap{`map[int]int`: {{semi}}})
}

func TestMethodDecl(t *testing.T) {
	remaining(t, MethodDecl, Tmap{
		`func (m M) f(){}`:                 {{semi}},
		`func (m M) f()(){}`:               {{semi}},
		`func (m M) f(int)int{ return 0 }`: {{semi}},
		`func (m *M) f() { a(); b(); }`:    {{semi}},
	})
}

func TestMethodExpr(t *testing.T) {
	remaining(t, MethodExpr, Tmap{
		`1`:        empty,
		`a.a`:      {{semi}},
		`a.a.a`:    {{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`(a).a`:    {{semi}},
		`(a.a).a`:  {{semi}},
		`(*a.a).a`: {{semi}},
	})
}

func TestMethodSpec(t *testing.T) {
	remaining(t, MethodSpec, Tmap{
		`a()`: {{semi}, {semi, {tok: token.RPAREN}, {tok: token.LPAREN}}},
		`a`:   {{semi}},
		`a.a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
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
		`1`:     {{semi}},
		`a.a`:   {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}, {semi}},
		`(a.a)`: {{semi}, {semi}, {semi}},
	})
}

func TestOperandName(t *testing.T) {
	remaining(t, OperandName, Tmap{
		`1`:   empty,
		`_`:   {{semi}},
		`a`:   {{semi}},
		`a.a`: {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}},
	})
}

func TestPackageName(t *testing.T) {
	remaining(t, PackageName, Tmap{
		`a`: {{semi}},
		`1`: empty,
		`_`: empty,
	})
}

func TestPackageClause(t *testing.T) {
	remaining(t, PackageClause, Tmap{
		`package a`: {{semi}},
		`package _`: empty,
		`package`:   empty,
		`a`:         empty,
	})
}

func TestParameterDecl(t *testing.T) {
	remaining(t, ParameterDecl, Tmap{
		`int`:      {{semi}},
		`...int`:   {{semi}},
		`a, b int`: {{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`b... int`: {{semi, {tok: token.IDENT, lit: `int`}, {tok: token.ELLIPSIS}}, {semi}},
	})
}

func TestParameterList(t *testing.T) {
	remaining(t, ParameterList, Tmap{
		`int`:      {{semi}},
		`int, int`: {{semi, {tok: token.IDENT, lit: `int`}, {tok: token.COMMA}}, {semi}},
		`a, b int, c, d int`: {
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}, {tok: token.IDENT, lit: `c`}, {tok: token.COMMA},
				{tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}, {tok: token.IDENT, lit: `c`}, {tok: token.COMMA}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}, {tok: token.IDENT, lit: `c`}, {tok: token.COMMA},
				{tok: token.IDENT, lit: `int`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}, {tok: token.IDENT, lit: `c`}, {tok: token.COMMA}},
			{semi}, {semi, {tok: token.IDENT, lit: `int`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.COMMA}},
			{semi}, {semi}, {semi, {tok: token.IDENT, lit: `int`}}, {semi},
		},
	})
}

func TestParameters(t *testing.T) {
	remaining(t, Parameters, Tmap{
		`(int)`:                {{semi}},
		`(int, int)`:           {{semi}},
		`(int, int,)`:          {{semi}},
		`(a, b int)`:           {{semi}, {semi}},
		`(a, b int, c, d int)`: {{semi}, {semi}, {semi}, {semi}},
	})
}

func TestPointerType(t *testing.T) {
	remaining(t, PointerType, Tmap{
		`*int`: {{semi}},
		`int`:  empty,
	})
}

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, Tmap{
		`1`: {{semi}},
		`(a.a)("foo")`: {
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi}, {semi}, {semi}, {semi},
		},
		`a.a`: {
			{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}, {semi}, {semi},
		},
		`a[1]`: {
			{semi, {tok: token.RBRACK}, {tok: token.INT, lit: `1`}, {tok: token.LBRACK}}, {semi},
		},
		`a[:]`: {
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}}, {semi},
		},
		`a.(int)`: {
			{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `int`}, {tok: token.LPAREN}, {tok: token.PERIOD}}, {semi},
		},
		`a(b...)`: {
			{semi, {tok: token.RPAREN}, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, {tok: token.LPAREN}}, {semi},
		},
		`a(b...)[:]`: {
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK},
				{tok: token.RPAREN}, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, {tok: token.LPAREN}},
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}},
			{semi},
		},
	})
}

func TestQualifiedIdent(t *testing.T) {
	remaining(t, QualifiedIdent, Tmap{
		`1`:   empty,
		`a`:   empty,
		`a.a`: {{semi}},
		`_.a`: empty,
		`a._`: {{semi}},
		`_._`: empty,
	})
}

func TestRangeClause(t *testing.T) {
	remaining(t, RangeClause, Tmap{
		`range a`:              {{semi}},
		`a[0], a[1] = range b`: {{semi}},
		`k, v := range a`:      {{semi}},
	})
}

func TestReceiverType(t *testing.T) {
	remaining(t, ReceiverType, Tmap{
		`1`:      empty,
		`a.a`:    {{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`(a.a)`:  {{semi}},
		`(*a.a)`: {{semi}},
	})
}

func TestRecvStmt(t *testing.T) {
	remaining(t, RecvStmt, Tmap{
		`<-a`: {{semi}},
		`a[0], b = <-c`: {{semi, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.ASSIGN}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA},
			{tok: token.RBRACK}, {tok: token.INT, lit: `0`}, {tok: token.LBRACK}},
			{semi, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.ASSIGN}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`a, b := <-c`: {{semi, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`a[0], b := <-c`: {{semi, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA},
			{tok: token.RBRACK}, {tok: token.INT, lit: `0`}, {tok: token.LBRACK}},
			{semi, {tok: token.IDENT, lit: `c`}, {tok: token.ARROW}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}},
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
		`int`:        {{semi}},
		`(int, int)`: {{semi}},
	})
}

func TestReturnStmt(t *testing.T) {
	remaining(t, ReturnStmt, Tmap{
		`return 1`: {{semi, {tok: token.INT, lit: `1`}}, {semi}},
		`return 1, 2`: {{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}, {tok: token.INT, lit: `1`}},
			{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestSelector(t *testing.T) {
	remaining(t, Selector, Tmap{
		`1`:  empty,
		`a`:  empty,
		`.a`: {{semi}},
	})
}

func TestSelectStmt(t *testing.T) {
	remaining(t, SelectStmt, Tmap{
		`select {}`:                                            {{semi}},
		`select {default:}`:                                    {{semi}},
		`select {default: a()}`:                                {{semi}},
		`select {default: a();}`:                               {{semi}, {semi}},
		`select {case <-a: ;default:}`:                         {{semi}},
		`select {case <-a: b(); case c<-d: e(); default: f()}`: {{semi}},
	})
}

func TestSendStmt(t *testing.T) {
	remaining(t, SendStmt, Tmap{`a <- 1`: {{semi}}})
}

func TestShortVarDecl(t *testing.T) {
	remaining(t, ShortVarDecl, Tmap{
		`a := 1`:       {{semi}},
		`a, b := 1, 2`: {{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestSignature(t *testing.T) {
	remaining(t, Signature, Tmap{
		`()`:             {{semi}},
		`()()`:           {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`(int, int) int`: {{semi, {tok: token.IDENT, lit: `int`}}, {semi}},
	})
}

func TestSimpleStmt(t *testing.T) {
	remaining(t, SimpleStmt, Tmap{
		`1`: {{semi, {tok: token.INT, lit: `1`}}, {semi}},
		`a <- 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.ARROW}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.ARROW}}, {semi}},
		`a++`: {{semi, {tok: token.INC}, {tok: token.IDENT, lit: `a`}}, {semi, {tok: token.INC}}, {semi}},
		`a = 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.ASSIGN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.ASSIGN}}, {semi}},
		`a := 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.DEFINE}}, {semi}},
	})
}

func TestSlice(t *testing.T) {
	remaining(t, Slice, Tmap{
		`[:]`:     {{semi}},
		`[1:]`:    {{semi}},
		`[:1]`:    {{semi}},
		`[1:1]`:   {{semi}},
		`[:1:1]`:  {{semi}},
		`[1:1:1]`: {{semi}},
	})
}

func TestSliceType(t *testing.T) {
	remaining(t, SliceType, Tmap{`[]int`: {{semi}}})
}

func TestSourceFile(t *testing.T) {
	remaining(t, SourceFile, Tmap{
		`package p`:             {{}},
		`package p; import "a"`: {{semi, {tok: token.STRING, lit: `"a"`}, {tok: token.IMPORT, lit: `import`}}, {}},
		`package p; var a int`:  {{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `a`}, {tok: token.VAR, lit: `var`}}, {}},
		`package p; import "a"; import "b"; var c int; var d int`: {
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `c`}, {tok: token.VAR, lit: `var`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.STRING, lit: `"b"`}, {tok: token.IMPORT, lit: `import`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.STRING, lit: `"a"`}, {tok: token.IMPORT, lit: `import`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `c`}, {tok: token.VAR, lit: `var`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.STRING, lit: `"b"`}, {tok: token.IMPORT, lit: `import`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `c`}, {tok: token.VAR, lit: `var`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `d`}, {tok: token.VAR, lit: `var`}}, {}},
	})
}

func TestStatement(t *testing.T) {
	remaining(t, Statement, Tmap{
		`var a int`: {{semi}, {semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `a`}, {tok: token.VAR, lit: `var`}}},
		`a: var b int`: {{semi},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}, {tok: token.COLON}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}, {tok: token.COLON}}},
		`a := 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.DEFINE}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.DEFINE}}, {semi}},
		`go a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.GO, lit: `go`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`return 1`:    {{semi, {tok: token.INT, lit: `1`}, {tok: token.RETURN, lit: `return`}}, {semi, {tok: token.INT, lit: `1`}}, {semi}},
		`break a`:     {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.BREAK, lit: `break`}}, {semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`continue a`:  {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.CONTINUE, lit: `continue`}}, {semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`goto a`:      {{semi, {tok: token.IDENT, lit: `a`}, {tok: token.GOTO, lit: `goto`}}, {semi}},
		`fallthrough`: {{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi}},
		`{a()}`:       {{semi, {tok: token.RBRACE}, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.LBRACE}}, {semi}},
		`if a {}`:     {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `a`}, {tok: token.IF, lit: `if`}}, {semi}},
		`switch {}`:   {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.SWITCH, lit: `switch`}}, {semi}},
		`select {}`:   {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.SELECT, lit: `select`}}, {semi}},
		`for {}`:      {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.FOR, lit: `for`}}, {semi}},
		`defer a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.DEFER, lit: `defer`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
	})
}

func TestStatementList(t *testing.T) {
	remaining(t, StatementList, Tmap{
		`fallthrough`: {{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi}, {}},
		`fallthrough;`: {{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{{tok: token.SEMICOLON, lit: `;`}}, {}},
		`fallthrough; fallthrough`: {{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}, {tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}, {tok: token.SEMICOLON, lit: `;`}},
			{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi}, {}},
	})
}

func TestStructType(t *testing.T) {
	remaining(t, StructType, Tmap{
		`struct{}`:                      {{semi}},
		`struct{int}`:                   {{semi}},
		`struct{int;}`:                  {{semi}},
		`struct{int;float64;}`:          {{semi}},
		`struct{a int}`:                 {{semi}},
		`struct{a, b int}`:              {{semi}},
		`struct{a, b int;}`:             {{semi}},
		`struct{a, b int; c, d string}`: {{semi}},
	})
}

func TestSwitchStmt(t *testing.T) {
	remaining(t, SwitchStmt, Tmap{
		`switch {}`:          {{semi}},
		`switch a.(type) {}`: {{semi}},
	})
}

func TestTopLevelDecl(t *testing.T) {
	remaining(t, TopLevelDecl, Tmap{
		`var a int`:        {{semi}},
		`func f(){}`:       {{semi}},
		`func (m M) f(){}`: {{semi}},
	})
}

func TestType(t *testing.T) {
	remaining(t, Type, Tmap{
		`a`:        {{semi}},
		`a.a`:      {{semi}, {semi, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:        empty,
		`_`:        {{semi}},
		`(a.a)`:    {{semi}},
		`(((_)))`:  {{semi}},
		`chan int`: {{semi}},
	})
}

func TestTypeAssertion(t *testing.T) {
	remaining(t, TypeAssertion, Tmap{
		`.(int)`: {{semi}},
		`1`:      empty,
	})
}

func TestTypeCaseClause(t *testing.T) {
	remaining(t, TypeCaseClause, Tmap{
		`case a:`: {{}},
		`case a: b()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `b`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}},
		`default: a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}},
		`case a, b: c(); d()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `d`}, {tok: token.SEMICOLON, lit: `;`},
			{tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `c`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `d`}, {tok: token.SEMICOLON, lit: `;`},
				{tok: token.RPAREN}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `d`}, {tok: token.SEMICOLON, lit: `;`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `d`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}},
	})
}

func TestTypeDecl(t *testing.T) {
	remaining(t, TypeDecl, Tmap{
		`type a int`:           {{semi}},
		`type (a int)`:         {{semi}},
		`type (a int;)`:        {{semi}},
		`type (a int; b int)`:  {{semi}},
		`type (a int; b int;)`: {{semi}},
	})
}

func TestTypeList(t *testing.T) {
	remaining(t, TypeList, Tmap{
		`a`:     {{semi}},
		`a, b`:  {{semi, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`a, b,`: {{{tok: token.COMMA}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {{tok: token.COMMA}}},
	})
}

func TestTypeLit(t *testing.T) {
	remaining(t, TypeLit, Tmap{
		`[1]int`:      {{semi}},
		`struct{}`:    {{semi}},
		`*int`:        {{semi}},
		`func()`:      {{semi}},
		`interface{}`: {{semi}},
		`[]int`:       {{semi}},
		`map[int]int`: {{semi}},
		`chan int`:    {{semi}},
	})
}

func TestTypeName(t *testing.T) {
	remaining(t, TypeName, Tmap{
		`a`:   {{semi}},
		`a.a`: {{semi}, {semi, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:   empty,
		`_`:   {{semi}},
	})
}

func TestTypeSpec(t *testing.T) {
	remaining(t, TypeSpec, Tmap{
		`a int`: {{semi}},
		`a`:     empty,
	})
}

func TestTypeSwitchCase(t *testing.T) {
	remaining(t, TypeSwitchCase, Tmap{
		`default`:   {{}},
		`case a`:    {{semi}},
		`case a, b`: {{semi, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestTypeSwitchGuard(t *testing.T) {
	remaining(t, TypeSwitchGuard, Tmap{
		`a.(type)`:      {{semi}},
		`a := a.(type)`: {{semi}},
	})
}

func TestTypeSwitchStmt(t *testing.T) {
	remaining(t, TypeSwitchStmt, Tmap{
		`switch a.(type) {}`:                              {{semi}},
		`switch a := 5; a := a.(type) {}`:                 {{semi}},
		`switch a.(type) { case int: }`:                   {{semi}},
		`switch a.(type) { case int:; }`:                  {{semi}, {semi}},
		`switch a.(type) { case int: b() }`:               {{semi}},
		`switch a.(type) { case int: b(); }`:              {{semi}, {semi}},
		`switch a.(type) { case int: b(); default: c() }`: {{semi}},
	})
}

func TestUnaryExpr(t *testing.T) {
	remaining(t, UnaryExpr, Tmap{
		`1`:  {{semi}},
		`-1`: {{semi}},
		`!a`: {{semi}},
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
		`var a int`:                 {{semi}},
		`var (a int)`:               {{semi}},
		`var (a int;)`:              {{semi}},
		`var (a, b = 1, 2)`:         {{semi}},
		`var (a, b = 1, 2; c int;)`: {{semi}},
	})
}

func TestVarSpec(t *testing.T) {
	remaining(t, VarSpec, Tmap{
		`a int`:       {{semi}},
		`a int = 1`:   {{semi, {tok: token.INT, lit: `1`}, {tok: token.ASSIGN}}, {semi}},
		`a = 1`:       {{semi}},
		`a, b = 1, 2`: {{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestTokenReader(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	defect.DeepEqual(t, tokenReader(toks, token.RPAREN), [][]*Token{{semi}})
}

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	tree, state := tokenParser(toks, token.RPAREN)
	defect.DeepEqual(t, state, [][]*Token{{semi}})
	defect.DeepEqual(t, tree, []interface{}{&Token{tok: token.RPAREN}})
}
