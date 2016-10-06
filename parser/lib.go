package parser

import "go/token"

func Klein(p Parser) Parser {
	return func(t []*Token) []*Token {
		for newT := p(t); newT != nil; newT = p(t) {
			t = newT
		}
		return t
	}
}

func And(ps ...Parser) Parser {
	return func(t []*Token) []*Token {
		for _, p := range ps {
			newT := p(t)
			if newT == nil {
				return nil
			}
			t = newT
		}
		return t
	}
}

func Maybe(p Parser) Parser {
	return func(t []*Token) []*Token {
		if newT := p(t); newT != nil {
			return newT
		}
		return t
	}
}

func Or(ps ...Parser) Multi {
	return func(ts [][]*Token) [][]*Token {
		var results [][]*Token
		for _, t := range ts {
			for _, p := range ps {
				if newT := p(t); newT != nil {
					results = append(results, newT)
				}
			}
		}
		return results
	}
}

func Each(p Parser) Multi {
	return func(ts [][]*Token) [][]*Token {
		var results [][]*Token
		for _, t := range ts {
			if newT := p(t); newT != nil {
				results = append(results, newT)
			}
		}
		return results
	}
}

func Basic(tok token.Token) Parser {
	return func(t []*Token) []*Token {
		if p := pop(&t); p == nil || p.tok != tok {
			return nil
		}
		return t
	}
}
