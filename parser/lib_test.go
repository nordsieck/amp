package parser

import (
	"go/scanner"
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
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

func remaining(t *testing.T, r Reader, m Tmap) {
	for raw, left := range m {
		toks := [][]*Token{Scan(newScanner(raw))}
		defect.DeepEqual(t, r(toks), left)
	}
}
