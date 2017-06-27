package parser

func tokenSliceToString(ts []*Token) string {
	if len(ts) == 0 {
		return `[]`
	}

	s := `[`
	for i := 0; i < len(ts)-1; i++ {
		s += ts[i].String() + ` `
	}
	s += ts[len(ts)-1].String() + `]`
	return s
}

func rendererSliceToString(r []Renderer) string {
	if len(r) == 0 {
		return `[]`
	}

	s := `[`
	for i := 0; i < len(r)-1; i++ {
		s += string(r[i].Render()) + ` `
	}
	s += string(r[len(r)-1].Render()) + `]`
	return s
}
