package main

func advance() {
	tokens = tokens[1:]
}

type expression interface {
	anExpr()
}

type intLit struct {
	val int
}

func (e *intLit) anExpr() {
	panic("implement me")
}

func parse() []expression {
	var ret []expression
	for len(tokens) > 0 {
		ret = append(ret, parseExpression())
	}
	return ret
}

func parseExpression() expression {
	return parseIntLit()
}

func parseIntLit() expression {

	tok := tokens[0]
	advance()

	return &intLit{
		val: tok.num,
	}
}
