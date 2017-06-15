package parser

func tokenSliceToString(ts []*Token) string {
	if len(ts) == 0 {
		return "[]"
	}

	s := "["
	for i := 0; i < len(ts)-1; i++ {
		s += ts[i].String()
		s += ` `
	}
	s += ts[len(ts)-1].String()
	s += "]"
	return s
}
