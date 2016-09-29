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

func TestBasicLit(t *testing.T) {
	consumes(t, BasicLit, map[string]bool{
		``:    false,
		`1`:   true,
		`1.1`: true,
		`1i`:  true,
		`'a'`: true,
		`"a"`: true,
		`a`:   false,
	})
}

func TestIdentifierList(t *testing.T) {
	consumes(t, IdentifierList, map[string]bool{
		`a`:   true,
		`a,a`: true,
		`1`:   false,
	})
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
