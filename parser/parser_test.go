package parser

import (
	"go/scanner"
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

var (
	semicolon = []*Token{{tok: token.SEMICOLON, lit: "\n"}}
	fail      = []*Token(nil)
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
	for raw, match := range map[string]bool{
		"1":   true,
		"1.1": true,
		"1i":  true,
		"'a'": true,
		`"a"`: true,
		"a":   false,
	} {
		toks := Scan(newScanner(raw))
		if match {
			defect.DeepEqual(t, BasicLit(toks), semicolon)
		} else {
			defect.DeepEqual(t, BasicLit(toks), fail)
		}
	}
}
