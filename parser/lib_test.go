package parser

import (
	"go/scanner"
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

var (
	semicolon      = &Token{tok: token.SEMICOLON, lit: "\n"}
	semicolonSlice = []*Token{semicolon}
	fail           = []*Token(nil)
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
	for raw, left := range map[string][]*Token{
		`1`:     semicolonSlice,
		`1 1 1`: semicolonSlice,
		`1 a`:   []*Token{semicolon, {tok: token.IDENT, lit: `a`}},
		`1 a 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}, {tok: token.IDENT, lit: `a`}},
	} {
		toks := Scan(newScanner(raw))
		defect.DeepEqual(t, Klein(BasicLit)(toks), left)
	}
}

func TestAnd(t *testing.T) {
	for raw, left := range map[string][]*Token{
		`1`:     fail,
		`1 1`:   semicolonSlice,
		`1 1 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}},
	} {
		toks := Scan(newScanner(raw))
		defect.DeepEqual(t, And(BasicLit, BasicLit)(toks), left)
	}
}

func TestMaybe(t *testing.T) {
	for raw, left := range map[string][]*Token{
		`1`:   semicolonSlice,
		`1 1`: []*Token{semicolon, {tok: token.INT, lit: `1`}},
		`a`:   []*Token{semicolon, {tok: token.IDENT, lit: `a`}},
	} {
		toks := Scan(newScanner(raw))
		defect.DeepEqual(t, Maybe(BasicLit)(toks), left)
	}
}

func TestOr(t *testing.T) {
	for raw, left := range map[string][][]*Token{
		`a`:   [][]*Token{semicolonSlice},
		`1`:   [][]*Token{semicolonSlice},
		`a.a`: [][]*Token{{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
	} {
		toks := [][]*Token{Scan(newScanner(raw))}
		defect.DeepEqual(t, Or(Basic(token.IDENT), Basic(token.INT))(toks), left)
	}

	for raw, left := range map[string][][]*Token{
		`a`:   [][]*Token{semicolonSlice},
		`a.a`: [][]*Token{semicolonSlice, {semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}}},
		`a.a.a`: [][]*Token{
			{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}},
			{semicolon, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}, {tok: token.IDENT, lit: "a"}, {tok: token.PERIOD}},
		},
		`1`: [][]*Token(nil),
	} {
		toks := [][]*Token{Scan(newScanner(raw))}
		defect.DeepEqual(t, Or(And(Basic(token.IDENT), Basic(token.PERIOD), Basic(token.IDENT)), Basic(token.IDENT))(toks), left)
	}
}

func TestEach(t *testing.T) {
	for raw, left := range map[string][][]*Token{
		`1`: [][]*Token{semicolonSlice},
		`a`: [][]*Token(nil),
		`_`: [][]*Token(nil),
	} {
		toks := [][]*Token{Scan(newScanner(raw))}
		defect.DeepEqual(t, Each(Basic(token.INT))(toks), left)
	}
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
