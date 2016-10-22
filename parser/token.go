package parser

import "go/token"

type Token struct {
	// pos token.Position
	tok token.Token
	lit string
}

func (t *Token) String() string { return "{" + t.tok.String() + " " + t.lit + "}" }

func pop(t *State) *Token {
	if len(*t) == 0 {
		return nil
	}
	top := (*t)[len(*t)-1]
	*t = (*t)[:len(*t)-1]
	return top
}

func reverse(s []*Token) {
	for i := 0; i < len(s)/2; i++ {
		s[i], s[len(s)-1-i] = s[len(s)-1-i], s[i]
	}
}
