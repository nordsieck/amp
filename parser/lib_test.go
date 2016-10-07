package parser

import (
	"go/scanner"
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

var (
	semicolon       = &Token{tok: token.SEMICOLON, lit: "\n"}
	semicolonSlice  = []*Token{semicolon}
	semicolonSlice2 = [][]*Token{semicolonSlice}
	fail            = []*Token(nil)
)

func newScanner(src string) *scanner.Scanner {
	var s scanner.Scanner
	fset := token.NewFileSet()
	file := fset.AddFile("", fset.Base(), len(src))
	s.Init(file, []byte(src), nil, scanner.ScanComments)

	return &s
}

func Scan(s *scanner.Scanner) []*Token {
	var t []*Token
	for _, tok, lit := s.Scan(); tok != token.EOF; _, tok, lit = s.Scan() {
		t = append(t, &Token{tok: tok, lit: lit})
	}
	reverse(t)
	return t
}

func TestKlein(t *testing.T) {
	remaining(t, Klein(BasicLit), map[string][][]*Token{
		`1`:     semicolonSlice2,
		`1 1 1`: semicolonSlice2,
		`1 a`:   [][]*Token{{semicolon, {tok: token.IDENT, lit: `a`}}},
		`1 a 1`: [][]*Token{{semicolon, {tok: token.INT, lit: `1`}, {tok: token.IDENT, lit: `a`}}},
	})

	toks := [][]*Token{
		{semicolon, {tok: token.INT, lit: `1`}, {tok: token.INT, lit: `1`}},
		{semicolon, {tok: token.INT, lit: `1`}},
		{semicolon, {tok: token.IDENT, lit: `a`}},
	}
	defect.DeepEqual(t, Klein(BasicLit)(toks), [][]*Token{semicolonSlice, semicolonSlice,
		{semicolon, {tok: token.IDENT, lit: `a`}}})
}

// func Testklein(t *testing.T) {
// 	for raw, left := range map[string][]*Token{
// 		`1`:     semicolonSlice,
// 		`1 1 1`: semicolonSlice,
// 		`1 a`:   []*Token{semicolon, {tok: token.IDENT, lit: `a`}},
// 		`1 a 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}, {tok: token.IDENT, lit: `a`}},
// 	} {
// 		toks := Scan(newScanner(raw))
// 		defect.DeepEqual(t, klein(basicLit)(toks), left)
// 	}
// }

// func Testand(t *testing.T) {
// 	for raw, left := range map[string][]*Token{
// 		`1`:     fail,
// 		`1 1`:   semicolonSlice,
// 		`1 1 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}},
// 	} {
// 		toks := Scan(newScanner(raw))
// 		defect.DeepEqual(t, and(basicLit, basicLit)(toks), left)
// 	}
// }

func TestMaybe(t *testing.T) {
	remaining(t, Maybe(Each(Basic(token.INT))), map[string][][]*Token{
		`1`:   semicolonSlice2,
		`a`:   [][]*Token{{semicolon, {tok: token.IDENT, lit: `a`}}},
		`_`:   [][]*Token{{semicolon, {tok: token.IDENT, lit: `_`}}},
		`1 1`: [][]*Token{{semicolon, {tok: token.INT, lit: `1`}}},
	})

	toks := [][]*Token{
		{semicolon, {tok: token.INT, lit: `1`}, {tok: token.INT, lit: `2`}},
		{semicolon, {tok: token.INT, lit: `3`}},
		{semicolon, {tok: token.IDENT, lit: `_`}},
	}
	defect.DeepEqual(t, Maybe(Each(Basic(token.INT)))(toks), [][]*Token{
		{semicolon, {tok: token.INT, lit: `1`}},
		semicolonSlice,
		{semicolon, {tok: token.IDENT, lit: `_`}},
	})
}

// func Testmaybe(t *testing.T) {
// 	for raw, left := range map[string][]*Token{
// 		`1`:   semicolonSlice,
// 		`1 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}},
// 		`a`:   []*Token{semicolon, {tok: token.IDENT, lit: `a`}},
// 	} {
// 		toks := Scan(newScanner(raw))
// 		defect.DeepEqual(t, maybe(basicLit)(toks), left)
// 	}
// }

func TestOr(t *testing.T) {
	remaining(t, Or(Each(Basic(token.IDENT)), Each(Basic(token.INT))), map[string][][]*Token{
		`a`:   semicolonSlice2,
		`1`:   semicolonSlice2,
		`a.a`: [][]*Token{{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
	})

	remaining(t, Or(Each(and(Basic(token.IDENT), Basic(token.PERIOD), Basic(token.IDENT))), Each(Basic(token.IDENT))), map[string][][]*Token{
		`a`:   [][]*Token{semicolonSlice},
		`a.a`: [][]*Token{semicolonSlice, {semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`a.a.a`: [][]*Token{
			{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}},
			{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}},
		},
		`1`: [][]*Token(nil),
	})
}

func TestEach(t *testing.T) {
	remaining(t, Each(Basic(token.INT)), map[string][][]*Token{
		`1`: semicolonSlice2,
		`a`: [][]*Token(nil),
		`_`: [][]*Token(nil),
	})
}

func TestBasic(t *testing.T) {
	consumes(t, Basic(token.INT), map[string]bool{
		`1`:   true,
		`a`:   false,
		`1.1`: false,
	})
}

func consumes(t *testing.T, p Parser, m map[string]bool) {
	for raw, match := range m {
		toks := Scan(newScanner(raw))
		if match {
			defect.DeepEqual(t, p(toks), semicolonSlice)
		} else {
			defect.DeepEqual(t, p(toks), fail)
		}
	}
}

func remaining(t *testing.T, m Multi, mp map[string][][]*Token) {
	for raw, left := range mp {
		toks := [][]*Token{Scan(newScanner(raw))}
		defect.DeepEqual(t, m(toks), left)
	}
}
