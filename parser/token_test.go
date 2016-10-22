package parser

import (
	"go/token"
	"testing"

	"github.com/nordsieck/defect"
)

func TestPop(t *testing.T) {
	s := State{{tok: token.INT, lit: "1"}, {tok: token.ADD, lit: "+"}, {tok: token.INT, lit: "2"}}
	defect.Equal(t, *pop(&s), Token{tok: token.INT, lit: "2"})
	defect.DeepEqual(t, s, State{{tok: token.INT, lit: "1"}, {tok: token.ADD, lit: "+"}})

	s = []*Token{}
	defect.Equal(t, pop(&s), (*Token)(nil))
	defect.DeepEqual(t, s, State{})
}

func TestReverse(t *testing.T) {
	s := State{{tok: token.INT, lit: "1"}, {tok: token.ADD, lit: "+"}, {tok: token.INT, lit: "2"}}
	reverse(s)
	defect.DeepEqual(t, s, State{{tok: token.INT, lit: "2"}, {tok: token.ADD, lit: "+"}, {tok: token.INT, lit: "1"}})
}
