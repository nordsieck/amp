package parser

import "go/token"

func Klein(m Multi) Multi {
	return func(ts [][]*Token) [][]*Token {
		var result [][]*Token
		for _, t := range ts {
			for newT := m([][]*Token{t}); newT != nil; newT = m([][]*Token{t}) {
				t = newT[0]
			}
			result = append(result, t)
		}
		return result
	}
}

func klein(p Parser) Parser {
	return func(t []*Token) []*Token {
		for newT := p(t); newT != nil; newT = p(t) {
			t = newT
		}
		return t
	}
}

func And(ms ...Multi) Multi {
	return func(ts [][]*Token) [][]*Token {
		for _, m := range ms {
			ts = m(ts)
		}
		return ts
	}
}

func and(ps ...Parser) Parser {
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

func Maybe(m Multi) Multi {
	return func(ts [][]*Token) [][]*Token {
		var result [][]*Token
		for _, t := range ts {
			if newT := m([][]*Token{t}); newT == nil {
				result = append(result, t)
			} else {
				result = append(result, newT...)
			}
		}
		return result
	}
}

func maybe(p Parser) Parser {
	return func(t []*Token) []*Token {
		if newT := p(t); newT != nil {
			return newT
		}
		return t
	}
}

func Or(ms ...Multi) Multi {
	return func(ts [][]*Token) [][]*Token {
		var results [][]*Token
		for _, m := range ms {
			results = append(results, m(ts)...)
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
