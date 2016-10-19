package parser

import (
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

var (
	empty     = [][]*Token(nil)
	semi      = &Token{tok: token.SEMICOLON, lit: "\n"}
	semiSlice = [][]*Token{{semi}}
)

func TestAddOp(t *testing.T) {
	remaining(t, AddOp, map[string][][]*Token{
		`+`: [][]*Token{{}},
		`-`: [][]*Token{{}},
		`|`: [][]*Token{{}},
		`&`: [][]*Token{{}},
		`1`: empty,
	})
}

func TestAnonymousField(t *testing.T) {
	remaining(t, AnonymousField, map[string][][]*Token{
		`a.a`:  [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`*a.a`: [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
	})
}

func TestArguments(t *testing.T) {
	remaining(t, Arguments, map[string][][]*Token{
		`()`:        semiSlice,
		`(a)`:       semiSlice,
		`(a,a)`:     semiSlice,
		`(a,a,)`:    semiSlice,
		`(a,a...,)`: semiSlice,
	})
}

func TestArrayType(t *testing.T) {
	remaining(t, ArrayType, map[string][][]*Token{
		`[1]int`: semiSlice,
		`[a]int`: semiSlice,
	})
}

func TestAssignment(t *testing.T) {
	remaining(t, Assignment, map[string][][]*Token{`a = 1`: semiSlice})
}

func TestAssignOp(t *testing.T) {
	remaining(t, AssignOp, map[string][][]*Token{
		`+=`: {{}}, `-=`: {{}}, `*=`: {{}}, `/=`: {{}},
		`%=`: {{}}, `&=`: {{}}, `|=`: {{}}, `^=`: {{}},
		`<<=`: {{}}, `>>=`: {{}}, `&^=`: {{}}, `=`: {{}},
	})
}

func TestBasicLit(t *testing.T) {
	remaining(t, BasicLit, map[string][][]*Token{
		``:    empty,
		`1`:   semiSlice,
		`1.1`: semiSlice,
		`1i`:  semiSlice,
		`'a'`: semiSlice,
		`"a"`: semiSlice,
		`a`:   empty,
		`_`:   empty,
	})
}

func TestBinaryOp(t *testing.T) {
	remaining(t, BinaryOp, map[string][][]*Token{
		`==`: [][]*Token{{}},
		`+`:  [][]*Token{{}},
		`*`:  [][]*Token{{}},
		`||`: [][]*Token{{}},
		`&&`: [][]*Token{{}},
		`1`:  empty,
	})
}

func TestBlock(t *testing.T) {
	remaining(t, Block, map[string][][]*Token{
		`{a()}`:     semiSlice,
		`{a();}`:    {{semi}, {semi}},
		`{a();b()}`: semiSlice,
	})
}

func TestBreakStmt(t *testing.T) {
	remaining(t, BreakStmt, map[string][][]*Token{
		`break a`: {{semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`a`:       empty,
	})
}

func TestChannelType(t *testing.T) {
	remaining(t, ChannelType, map[string][][]*Token{
		`chan int`:   semiSlice,
		`<-chan int`: semiSlice,
		`chan<- int`: semiSlice,
		`int`:        empty,
	})
}

func TestCompositeLit(t *testing.T) {
	remaining(t, CompositeLit, map[string][][]*Token{
		`T{1}`: semiSlice,
		`T{foo: "bar", baz: "quux",}`: [][]*Token{{semi}, {semi}, {semi}, {semi}},
	})
}

func TestConstDecl(t *testing.T) {
	remaining(t, ConstDecl, map[string][][]*Token{
		`const a = 1`:                         semiSlice,
		`const a int = 1`:                     semiSlice,
		`const (a = 1)`:                       semiSlice,
		`const (a int = 1)`:                   semiSlice,
		`const (a, b int = 1, 2; c int = 3;)`: semiSlice,
	})
}

func TestConstSpec(t *testing.T) {
	remaining(t, ConstSpec, map[string][][]*Token{
		`a = 1`:     semiSlice,
		`a int = 1`: semiSlice,
		`a, b int = 1, 2`: [][]*Token{
			{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi},
		},
	})
}

func TestContinueStmt(t *testing.T) {
	remaining(t, ContinueStmt, map[string][][]*Token{
		`continue a`: {{semi, {tok: token.IDENT, lit: `a`}}, {semi}},
		`a`:          empty,
	})
}

func TestConversion(t *testing.T) {
	remaining(t, Conversion, map[string][][]*Token{
		`float(1)`:   semiSlice,
		`(int)(5,)`:  semiSlice,
		`a.a("foo")`: semiSlice,
	})
}

func TestDeclaration(t *testing.T) {
	remaining(t, Declaration, map[string][][]*Token{
		`const a = 1`: semiSlice,
		`type a int`:  semiSlice,
		`var a int`:   semiSlice,
	})
}

func TestElement(t *testing.T) {
	remaining(t, Element, map[string][][]*Token{
		`1+1`:   [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.ADD}}, {semi}},
		`{1+1}`: semiSlice,
	})
}

func TestElementList(t *testing.T) {
	remaining(t, ElementList, map[string][][]*Token{
		`1:1`: [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}}, {semi}},
		`1:1, 1:1`: [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}, {tok: token.INT, lit: `1`},
			{tok: token.COMMA}, {tok: token.INT, lit: `1`}, {tok: token.COLON}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}, {tok: token.INT, lit: `1`}, {tok: token.COMMA}},
			{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}},
			{semi}},
	})
}

func TestEllipsisArrayType(t *testing.T) {
	remaining(t, EllipsisArrayType, map[string][][]*Token{`[...]int`: semiSlice})
}

func TestEmptyStmt(t *testing.T) {
	remaining(t, EmptyStmt, map[string][][]*Token{`1`: {{semi, {tok: token.INT, lit: `1`}}}})
}

func TestExprCaseClause(t *testing.T) {
	remaining(t, ExprCaseClause, map[string][][]*Token{
		`default: a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}, {}, {}},
	})
}

func TestExpression(t *testing.T) {
	remaining(t, Expression, map[string][][]*Token{
		`1`: semiSlice,
		`1+1`: [][]*Token{
			{semi, {tok: token.INT, lit: `1`}, {tok: token.ADD}},
			{semi},
		},
		`1+-1`: [][]*Token{
			{semi, {tok: token.INT, lit: `1`}, {tok: token.SUB}, {tok: token.ADD}},
			{semi},
		},
	})
}

func TestExpressionList(t *testing.T) {
	remaining(t, ExpressionList, map[string][][]*Token{
		`1`:   semiSlice,
		`1,1`: [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestExpressionStmt(t *testing.T) {
	remaining(t, ExpressionStmt, map[string][][]*Token{`1`: semiSlice})
}

func TestExprSwitchCase(t *testing.T) {
	remaining(t, ExprSwitchCase, map[string][][]*Token{
		`case 1`:    {{semi}},
		`case 1, 1`: {{semi, {tok: token.INT, lit: `1`}, {tok: token.COMMA}}, {semi}},
		`default`:   {{}},
	})
}

func TestExprSwitchStmt(t *testing.T) {
	remaining(t, ExprSwitchStmt, map[string][][]*Token{
		`switch{}`:                                               semiSlice,
		`switch{default:}`:                                       {{semi}, {semi}},
		`switch a := 1; {}`:                                      semiSlice,
		`switch a {}`:                                            semiSlice,
		`switch a := 1; a {}`:                                    semiSlice,
		`switch a := 1; a {default:}`:                            {{semi}, {semi}},
		`switch {case true:}`:                                    {{semi}, {semi}},
		`switch {case true: a()}`:                                semiSlice,
		`switch{ case true: a(); case false: b(); default: c()}`: {{semi}, {semi}, {semi}, {semi}},
	})
}

func TestFallthroughStmt(t *testing.T) {
	remaining(t, FallthroughStmt, map[string][][]*Token{
		`fallthrough`: semiSlice,
		`a`:           empty,
	})
}

func TestFieldDecl(t *testing.T) {
	remaining(t, FieldDecl, map[string][][]*Token{
		`a,a int`:    [][]*Token{{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `a`}, {tok: token.COMMA}}, {semi}},
		`*int`:       semiSlice,
		`*int "foo"`: [][]*Token{{semi, {tok: token.STRING, lit: `"foo"`}}, {semi}},
	})
}

func TestFunctionType(t *testing.T) {
	remaining(t, FunctionType, map[string][][]*Token{
		`func()`:             semiSlice,
		`func()()`:           [][]*Token{{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`func(int, int) int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}}, {semi}},
	})
}

func TestGoStmt(t *testing.T) {
	remaining(t, GoStmt, map[string][][]*Token{`go a()`: {{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}}})
}

func TestGotoStmt(t *testing.T) {
	remaining(t, GotoStmt, map[string][][]*Token{`goto a`: semiSlice, `goto`: empty, `a`: empty})
}

func TestIdentifierList(t *testing.T) {
	remaining(t, IdentifierList, map[string][][]*Token{
		`a`:   semiSlice,
		`a,a`: [][]*Token{{semi, {tok: token.IDENT, lit: `a`}, {tok: token.COMMA}}, {semi}},
		`1`:   empty,
		`_`:   semiSlice,
	})
}

func TestIfStmt(t *testing.T) {
	remaining(t, IfStmt, map[string][][]*Token{
		`if a {}`:             {{semi}, {semi}},
		`if a := false; a {}`: {{semi}, {semi}},
		`if a {} else {}`: {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi}, {semi}, {semi}, {semi}},
		`if a {} else if b {}`: {
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`}, {tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`}, {tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}},
			{semi}, {semi}, {semi}, {semi}},
		`if a := false; a {} else if b {} else {}`: {
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`},
				{tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `b`},
				{tok: token.IF, lit: `if`}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.ELSE, lit: `else`}},
			{semi}, {semi}, {semi}, {semi}, {semi}, {semi}, {semi}, {semi},
		},
		`if a { b() }`: {{semi}},
	})
}

func TestIncDecStmt(t *testing.T) {
	remaining(t, IncDecStmt, map[string][][]*Token{
		`a++`: semiSlice,
		`a--`: semiSlice,
	})
}

func TestIndex(t *testing.T) {
	remaining(t, Index, map[string][][]*Token{
		`[a]`: semiSlice,
		`[1]`: semiSlice,
	})
}

func TestInterfaceType(t *testing.T) {
	remaining(t, InterfaceType, map[string][][]*Token{
		`interface{}`:       semiSlice,
		`interface{a}`:      semiSlice,
		`interface{a()}`:    semiSlice,
		`interface{a();a;}`: semiSlice,
	})
}

func TestKey(t *testing.T) {
	remaining(t, Key, map[string][][]*Token{
		`a`:     [][]*Token{{semi}, {semi}},
		`1+a`:   [][]*Token{{semi, {tok: token.IDENT, lit: `a`}, {tok: token.ADD}}, {semi}},
		`"foo"`: semiSlice,
	})
}

func TestKeyedElement(t *testing.T) {
	remaining(t, KeyedElement, map[string][][]*Token{
		`1`:   semiSlice,
		`1:1`: [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.COLON}}, {semi}},
	})
}

func TestLabeledStmt(t *testing.T) {
	remaining(t, LabeledStmt, map[string][][]*Token{
		`a: var b int`: {{semi}, {semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.VAR, lit: `var`}}},
	})
}

func TestLiteral(t *testing.T) {
	remaining(t, Literal, map[string][][]*Token{
		`1`:    semiSlice,
		`T{1}`: semiSlice,
		`a`:    empty,
		`_`:    empty,
	})
}

func TestLiteralType(t *testing.T) {
	remaining(t, LiteralType, map[string][][]*Token{
		`struct{}`:    semiSlice,
		`[1]int`:      semiSlice,
		`[...]int`:    semiSlice,
		`[]int`:       semiSlice,
		`map[int]int`: semiSlice,
		`a.a`:         [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
	})
}

func TestLiteralValue(t *testing.T) {
	remaining(t, LiteralValue, map[string][][]*Token{
		`{1}`:           semiSlice,
		`{0: 1, 1: 2,}`: semiSlice,
	})
}

func TestMapType(t *testing.T) {
	remaining(t, MapType, map[string][][]*Token{`map[int]int`: semiSlice})
}

func TestMethodExpr(t *testing.T) {
	remaining(t, MethodExpr, map[string][][]*Token{
		`1`:        empty,
		`a.a`:      semiSlice,
		`a.a.a`:    [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`(a).a`:    semiSlice,
		`(a.a).a`:  semiSlice,
		`(*a.a).a`: semiSlice,
	})
}

func TestMethodSpec(t *testing.T) {
	remaining(t, MethodSpec, map[string][][]*Token{
		`a()`: [][]*Token{{semi}, {semi, {tok: token.RPAREN}, {tok: token.LPAREN}}},
		`a`:   semiSlice,
		`a.a()`: [][]*Token{{semi, {tok: token.RPAREN}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.LPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
	})
}

func TestMulOp(t *testing.T) {
	remaining(t, MulOp, map[string][][]*Token{
		`*`:  [][]*Token{{}},
		`/`:  [][]*Token{{}},
		`%`:  [][]*Token{{}},
		`<<`: [][]*Token{{}},
		`>>`: [][]*Token{{}},
		`&`:  [][]*Token{{}},
		`&^`: [][]*Token{{}},
		`1`:  empty,
	})
}

func TestOperand(t *testing.T) {
	remaining(t, Operand, map[string][][]*Token{
		`1`:     semiSlice,
		`a.a`:   [][]*Token{{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}, {semi}},
		`(a.a)`: [][]*Token{{semi}, {semi}, {semi}},
	})
}

func TestOperandName(t *testing.T) {
	remaining(t, OperandName, map[string][][]*Token{
		`1`:   empty,
		`_`:   semiSlice,
		`a`:   semiSlice,
		`a.a`: [][]*Token{{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}},
	})
}

func TestPackageName(t *testing.T) {
	remaining(t, PackageName, map[string][][]*Token{
		`a`: semiSlice,
		`1`: empty,
		`_`: empty,
	})
}

func TestParameterDecl(t *testing.T) {
	remaining(t, ParameterDecl, map[string][][]*Token{
		`int`:      semiSlice,
		`...int`:   semiSlice,
		`a, b int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`b... int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}, {tok: token.ELLIPSIS}}, {semi}},
	})
}

func TestParameterList(t *testing.T) {
	remaining(t, ParameterList, map[string][][]*Token{
		`int`:      semiSlice,
		`int, int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}, {tok: token.COMMA}}, {semi}},
		`a, b int, c, d int`: [][]*Token{
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
	remaining(t, Parameters, map[string][][]*Token{
		`(int)`:                semiSlice,
		`(int, int)`:           semiSlice,
		`(int, int,)`:          semiSlice,
		`(a, b int)`:           [][]*Token{{semi}, {semi}},
		`(a, b int, c, d int)`: [][]*Token{{semi}, {semi}, {semi}, {semi}},
	})
}

func TestPointerType(t *testing.T) {
	remaining(t, PointerType, map[string][][]*Token{
		`*int`: semiSlice,
		`int`:  empty,
	})
}

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, map[string][][]*Token{
		`1`: semiSlice,
		`(a.a)("foo")`: [][]*Token{
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi, {tok: token.RPAREN}, {tok: token.STRING, lit: `"foo"`}, {tok: token.LPAREN}},
			{semi}, {semi}, {semi}, {semi},
		},
		`a.a`: [][]*Token{
			{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}, {semi}, {semi},
		},
		`a[1]`: [][]*Token{
			{semi, {tok: token.RBRACK}, {tok: token.INT, lit: `1`}, {tok: token.LBRACK}}, {semi},
		},
		`a[:]`: [][]*Token{
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}}, {semi},
		},
		`a.(int)`: [][]*Token{
			{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `int`}, {tok: token.LPAREN}, {tok: token.PERIOD}}, {semi},
		},
		`a(b...)`: [][]*Token{
			{semi, {tok: token.RPAREN}, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, {tok: token.LPAREN}}, {semi},
		},
		`a(b...)[:]`: [][]*Token{
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK},
				{tok: token.RPAREN}, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, {tok: token.LPAREN}},
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}},
			{semi},
		},
	})
}

func TestQualifiedIdent(t *testing.T) {
	remaining(t, QualifiedIdent, map[string][][]*Token{
		`1`:   empty,
		`a`:   empty,
		`a.a`: semiSlice,
		`_.a`: empty,
		`a._`: semiSlice,
		`_._`: empty,
	})
}

func TestReceiverType(t *testing.T) {
	remaining(t, ReceiverType, map[string][][]*Token{
		`1`:      empty,
		`a.a`:    [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`(a.a)`:  semiSlice,
		`(*a.a)`: semiSlice,
	})
}

func TestRelOp(t *testing.T) {
	remaining(t, RelOp, map[string][][]*Token{
		`==`: [][]*Token{{}},
		`!=`: [][]*Token{{}},
		`>`:  [][]*Token{{}},
		`>=`: [][]*Token{{}},
		`<`:  [][]*Token{{}},
		`<=`: [][]*Token{{}},
		`1`:  empty,
	})
}

func TestResult(t *testing.T) {
	remaining(t, Result, map[string][][]*Token{
		`int`:        semiSlice,
		`(int, int)`: semiSlice,
	})
}

func TestReturnStmt(t *testing.T) {
	remaining(t, ReturnStmt, map[string][][]*Token{
		`return 1`: {{semi, {tok: token.INT, lit: `1`}}, {semi}},
		`return 1, 2`: {{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}, {tok: token.INT, lit: `1`}},
			{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestSelector(t *testing.T) {
	remaining(t, Selector, map[string][][]*Token{
		`1`:  empty,
		`a`:  empty,
		`.a`: semiSlice,
	})
}

func TestSendStmt(t *testing.T) {
	remaining(t, SendStmt, map[string][][]*Token{`a <- 1`: semiSlice})
}

func TestShortVarDecl(t *testing.T) {
	remaining(t, ShortVarDecl, map[string][][]*Token{
		`a := 1`:       semiSlice,
		`a, b := 1, 2`: {{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestSignature(t *testing.T) {
	remaining(t, Signature, map[string][][]*Token{
		`()`:             semiSlice,
		`()()`:           [][]*Token{{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`(int, int) int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}}, {semi}},
	})
}

func TestSimpleStmt(t *testing.T) {
	remaining(t, SimpleStmt, map[string][][]*Token{
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
	remaining(t, Slice, map[string][][]*Token{
		`[:]`:     semiSlice,
		`[1:]`:    semiSlice,
		`[:1]`:    semiSlice,
		`[1:1]`:   semiSlice,
		`[:1:1]`:  semiSlice,
		`[1:1:1]`: semiSlice,
	})
}

func TestSliceType(t *testing.T) {
	remaining(t, SliceType, map[string][][]*Token{`[]int`: semiSlice})
}

func TestStatement(t *testing.T) {
	remaining(t, Statement, map[string][][]*Token{
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
		`if a {}`:     {{semi, {tok: token.RBRACE}, {tok: token.LBRACE}, {tok: token.IDENT, lit: `a`}, {tok: token.IF, lit: `if`}}, {semi}, {semi}},
	})
}

func TestStatementList(t *testing.T) {
	remaining(t, StatementList, map[string][][]*Token{
		`fallthrough`: {{semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {semi}, {}, {}},
		`fallthrough;`: {{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{{tok: token.SEMICOLON, lit: `;`}}, {}, {}},
		`fallthrough; fallthrough;`: {{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}, {tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}, {tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}},
			{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}, {tok: token.SEMICOLON, lit: `;`}},
			{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {{tok: token.SEMICOLON, lit: `;`}}, {},
			{{tok: token.SEMICOLON, lit: `;`}, {tok: token.FALLTHROUGH, lit: `fallthrough`}}, {}},
	})
}

func TestStructType(t *testing.T) {
	remaining(t, StructType, map[string][][]*Token{
		`struct{}`:                      semiSlice,
		`struct{int}`:                   semiSlice,
		`struct{int;}`:                  semiSlice,
		`struct{int;float64;}`:          semiSlice,
		`struct{a int}`:                 semiSlice,
		`struct{a, b int}`:              semiSlice,
		`struct{a, b int;}`:             semiSlice,
		`struct{a, b int; c, d string}`: semiSlice,
	})
}

func TestType(t *testing.T) {
	remaining(t, Type, map[string][][]*Token{
		`a`:        semiSlice,
		`a.a`:      [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:        empty,
		`_`:        semiSlice,
		`(a.a)`:    semiSlice,
		`(((_)))`:  semiSlice,
		`chan int`: semiSlice,
	})
}

func TestTypeAssertion(t *testing.T) {
	remaining(t, TypeAssertion, map[string][][]*Token{
		`.(int)`: semiSlice,
		`1`:      empty,
	})
}

func TestTypeDecl(t *testing.T) {
	remaining(t, TypeDecl, map[string][][]*Token{
		`type a int`:           semiSlice,
		`type (a int)`:         semiSlice,
		`type (a int;)`:        semiSlice,
		`type (a int; b int)`:  semiSlice,
		`type (a int; b int;)`: semiSlice,
	})
}

func TestTypeList(t *testing.T) {
	remaining(t, TypeList, map[string][][]*Token{
		`a`:     semiSlice,
		`a, b`:  {{semi, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
		`a, b,`: {{{tok: token.COMMA}, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {{tok: token.COMMA}}},
	})
}

func TestTypeLit(t *testing.T) {
	remaining(t, TypeLit, map[string][][]*Token{
		`[1]int`:      semiSlice,
		`struct{}`:    semiSlice,
		`*int`:        semiSlice,
		`func()`:      semiSlice,
		`interface{}`: semiSlice,
		`[]int`:       semiSlice,
		`map[int]int`: semiSlice,
		`chan int`:    semiSlice,
	})
}

func TestTypeName(t *testing.T) {
	remaining(t, TypeName, map[string][][]*Token{
		`a`:   semiSlice,
		`a.a`: [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:   empty,
		`_`:   semiSlice,
	})
}

func TestTypeSpec(t *testing.T) {
	remaining(t, TypeSpec, map[string][][]*Token{
		`a int`: semiSlice,
		`a`:     empty,
	})
}

func TestTypeSwitchCase(t *testing.T) {
	remaining(t, TypeSwitchCase, map[string][][]*Token{
		`default`:   {{}},
		`case a`:    semiSlice,
		`case a, b`: {{semi, {tok: token.IDENT, lit: `b`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestUnaryExpr(t *testing.T) {
	remaining(t, UnaryExpr, map[string][][]*Token{
		`1`:  semiSlice,
		`-1`: semiSlice,
		`!a`: semiSlice,
	})
}

func TestUnaryOp(t *testing.T) {
	remaining(t, UnaryOp, map[string][][]*Token{
		`+`:  [][]*Token{{}},
		`-`:  [][]*Token{{}},
		`!`:  [][]*Token{{}},
		`^`:  [][]*Token{{}},
		`*`:  [][]*Token{{}},
		`&`:  [][]*Token{{}},
		`<-`: [][]*Token{{}},
		`1`:  empty,
	})
}

func TestVarDecl(t *testing.T) {
	remaining(t, VarDecl, map[string][][]*Token{
		`var a int`:                 semiSlice,
		`var (a int)`:               semiSlice,
		`var (a int;)`:              semiSlice,
		`var (a, b = 1, 2)`:         semiSlice,
		`var (a, b = 1, 2; c int;)`: semiSlice,
	})
}

func TestVarSpec(t *testing.T) {
	remaining(t, VarSpec, map[string][][]*Token{
		`a int`:       semiSlice,
		`a int = 1`:   [][]*Token{{semi, {tok: token.INT, lit: `1`}, {tok: token.ASSIGN}}, {semi}},
		`a = 1`:       semiSlice,
		`a, b = 1, 2`: [][]*Token{{semi, {tok: token.INT, lit: `2`}, {tok: token.COMMA}}, {semi}},
	})
}

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	defect.DeepEqual(t, tokenParser(toks, token.RPAREN), semiSlice)
}
