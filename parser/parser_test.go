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

func TestArguments(t *testing.T) {
	remaining(t, Arguments, map[string][][]*Token{
		`()`:        semiSlice,
		`(a)`:       semiSlice,
		`(a,a)`:     semiSlice,
		`(a,a,)`:    semiSlice,
		`(a,a...,)`: semiSlice,
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

func TestConversion(t *testing.T) {
	remaining(t, Conversion, map[string][][]*Token{
		`float(1)`:   semiSlice,
		`(int)(5,)`:  semiSlice,
		`a.a("foo")`: semiSlice,
	})
}

func TestExpression(t *testing.T) {
	remaining(t, Expression, map[string][][]*Token{
		`1`: semiSlice,
	})
}

func TestExpressionList(t *testing.T) {
	remaining(t, ExpressionList, map[string][][]*Token{
		`1`:   semiSlice,
		`1,1`: semiSlice,
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

func TestLiteral(t *testing.T) {
	remaining(t, Literal, map[string][][]*Token{
		`1`: semiSlice,
		`a`: empty,
		`_`: empty,
	})
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

func TestSelector(t *testing.T) {
	remaining(t, Selector, map[string][][]*Token{
		`1`:  empty,
		`a`:  empty,
		`.a`: semiSlice,
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

func TestType(t *testing.T) {
	remaining(t, Type, map[string][][]*Token{
		`a`:       semiSlice,
		`a.a`:     [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:       empty,
		`_`:       semiSlice,
		`(a.a)`:   semiSlice,
		`(((_)))`: semiSlice,
	})
}

func TestTypeAssertion(t *testing.T) {
	remaining(t, TypeAssertion, map[string][][]*Token{
		`.(int)`: semiSlice,
		`1`:      empty,
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
