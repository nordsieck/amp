package parser

import (
	"go/token"
	"testing"
)

func TestBasicLit(t *testing.T) {
	consumes(t, basicLit, map[string]bool{
		``:    false,
		`1`:   true,
		`1.1`: true,
		`1i`:  true,
		`'a'`: true,
		`"a"`: true,
		`a`:   false,
		`_`:   false,
	})
}

func TestIdentifierList(t *testing.T) {
	consumes(t, identifierList, map[string]bool{
		`a`:   true,
		`a,a`: true,
		`1`:   false,
		`_`:   true,
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
	consumes(t, qualifiedIdent, map[string]bool{
		`1`:   false,
		`a`:   false,
		`a.a`: true,
		`_.a`: false,
		`a._`: true,
		`_._`: false,
	})
}

func TestPackageName(t *testing.T) {
	consumes(t, packageName, map[string]bool{
		`a`: true,
		`1`: false,
		`_`: false,
	})
}
