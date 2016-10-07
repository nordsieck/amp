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

func TestTokenParser(t *testing.T) {
	toks := [][]*Token{
		{semi, {tok: token.RPAREN}},
		{semi, {tok: token.RPAREN}, {tok: token.IDENT, lit: `a`}, {tok: token.PERIOD}},
	}
	defect.DeepEqual(t, tokenParser(toks, token.RPAREN), semiSlice)
}
