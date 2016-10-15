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

func TestConversion(t *testing.T) {
	remaining(t, Conversion, map[string][][]*Token{
		`float(1)`:   semiSlice,
		`(int)(5,)`:  semiSlice,
		`a.a("foo")`: semiSlice,
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
		`1,1`: semiSlice,
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

func TestIdentifierList(t *testing.T) {
	remaining(t, IdentifierList, map[string][][]*Token{
		`a`:   semiSlice,
		`a,a`: semiSlice,
		`1`:   empty,
		`_`:   semiSlice,
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
			{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
			{semi}, {semi}, {semi},
		},
		`a[1]`: [][]*Token{
			{semi, {tok: token.RBRACK}, {tok: token.INT, lit: `1`}, {tok: token.LBRACK}},
			{semi}, {semi},
		},
		`a[:]`: [][]*Token{
			{semi, {tok: token.RBRACK}, {tok: token.COLON}, {tok: token.LBRACK}},
			{semi},
		},
		`a.(int)`: [][]*Token{
			{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `int`}, {tok: token.LPAREN}, {tok: token.PERIOD}},
			{semi},
		},
		`a(b...)`: [][]*Token{
			{semi, {tok: token.RPAREN}, {tok: token.ELLIPSIS}, {tok: token.IDENT, lit: `b`}, {tok: token.LPAREN}},
			{semi},
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

func TestSelector(t *testing.T) {
	remaining(t, Selector, map[string][][]*Token{
		`1`:  empty,
		`a`:  empty,
		`.a`: semiSlice,
	})
}

func TestSignature(t *testing.T) {
	remaining(t, Signature, map[string][][]*Token{
		`()`:             semiSlice,
		`()()`:           [][]*Token{{semi, {tok: token.RPAREN}, {tok: token.LPAREN}}, {semi}},
		`(int, int) int`: [][]*Token{{semi, {tok: token.IDENT, lit: `int`}}, {semi}},
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

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	defect.DeepEqual(t, tokenParser(toks, token.RPAREN), semiSlice)
}
