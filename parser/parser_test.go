package parser

import (
	"go/token"
	"testing"
)

func TestBasicLit(t *testing.T) {
	remaining(t, BasicLit, map[string][][]*Token{
		``:    [][]*Token(nil),
		`1`:   semicolonSlice2,
		`1.1`: semicolonSlice2,
		`1i`:  semicolonSlice2,
		`'a'`: semicolonSlice2,
		`"a"`: semicolonSlice2,
		`a`:   [][]*Token(nil),
		`_`:   [][]*Token(nil),
	})
}

func TestIdentifierList(t *testing.T) {
	remaining(t, IdentifierList, map[string][][]*Token{
		`a`:   semicolonSlice2,
		`a,a`: semicolonSlice2,
		`1`:   [][]*Token(nil),
		`_`:   semicolonSlice2,
	})
}

func TestTypeName(t *testing.T) {
	remaining(t, TypeName, map[string][][]*Token{
		`a`:   semicolonSlice2,
		`a.a`: [][]*Token{semicolonSlice, {semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`1`:   [][]*Token(nil),
		`_`:   semicolonSlice2,
	})
}

func TestQualifiedIdent(t *testing.T) {
	remaining(t, QualifiedIdent, map[string][][]*Token{
		`1`:   [][]*Token(nil),
		`a`:   [][]*Token(nil),
		`a.a`: semicolonSlice2,
		`_.a`: [][]*Token(nil),
		`a._`: semicolonSlice2,
		`_._`: [][]*Token(nil),
	})
}

func TestPackageName(t *testing.T) {
	remaining(t, PackageName, map[string][][]*Token{
		`a`: semicolonSlice2,
		`1`: [][]*Token(nil),
		`_`: [][]*Token(nil),
	})
}
