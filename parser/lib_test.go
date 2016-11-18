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

func result(t *testing.T, p Parser, o Omap) {
	for raw, left := range o {
		toks := [][]*Token{Scan(newScanner(raw))}
		tree, state := p(toks)
		var rendered []string
		for _, t := range tree {
			rendered = append(rendered, string(t.Render()))
		}
		defect.DeepEqual(t, state, left.rest)
		defect.DeepEqual(t, rendered, left.tree)
		defect.Equal(t, len(tree), len(state))
	}
}

func resultState(t *testing.T, s Stator, m map[string][]StateOutput) {
	for raw, result := range m {
		toks := Scan(newScanner(raw))
		state := s([]State{{[]Renderer{e{}}, toks}})
		so := GetStateOutput(state)
		defect.DeepEqual(t, so, result)

		// Use for understanding test output
		// for i := range state {
		// 	defect.DeepEqual(t, len(so[i].s), len(result[i].s))
		// 	defect.DeepEqual(t, so[i].s, result[i].s)
		// 	defect.DeepEqual(t, so[i].t, result[i].t)
		// }
	}
}
