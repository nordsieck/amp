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

func TestExpression(t *testing.T) {
	remaining(t, Expression, map[string][][]*Token{
		`1`: semiSlice,
	})
}

func TestUnaryExpr(t *testing.T) {
	remaining(t, UnaryExpr, map[string][][]*Token{
		`1`: semiSlice,
	})
}

func TestPrimaryExpr(t *testing.T) {
	remaining(t, PrimaryExpr, map[string][][]*Token{
		`1`: semiSlice,
	})
}

func TestOperand(t *testing.T) {
	remaining(t, Operand, map[string][][]*Token{
		`1`:     semiSlice,
		`a.a`:   [][]*Token{{semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}, {semi}},
		`(a.a)`: semiSlice,
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

func TestLiteral(t *testing.T) {
	remaining(t, Literal, map[string][][]*Token{
		`1`: semiSlice,
		`a`: empty,
		`_`: empty,
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

func TestIdentifierList(t *testing.T) {
	remaining(t, IdentifierList, map[string][][]*Token{
		`a`:   semiSlice,
		`a,a`: semiSlice,
		`1`:   empty,
		`_`:   semiSlice,
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

func TestReceiverType(t *testing.T) {
	remaining(t, ReceiverType, map[string][][]*Token{
		`1`:      empty,
		`a.a`:    [][]*Token{{semi}, {semi, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}}},
		`(a.a)`:  semiSlice,
		`(*a.a)`: semiSlice,
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

func TestPackageName(t *testing.T) {
	remaining(t, PackageName, map[string][][]*Token{
		`a`: semiSlice,
		`1`: empty,
		`_`: empty,
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

func TestAddOp(t *testing.T) {
	remaining(t, AddOp, map[string][][]*Token{
		`+`: [][]*Token{{}},
		`-`: [][]*Token{{}},
		`|`: [][]*Token{{}},
		`&`: [][]*Token{{}},
		`1`: empty,
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

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	defect.DeepEqual(t, tokenParser(toks, token.RPAREN), semiSlice)
}
