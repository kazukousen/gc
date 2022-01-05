package main

func advance() {
	tokens = tokens[1:]
}

func consume(s string) bool {
	if len(tokens) > 0 && tokens[0].val == s {
		advance()
		return true
	}
	return false
}

type expression interface {
	anExpr()
}

type intLit struct {
	expression
	val int
}

type binary struct {
	expression
	op  string
	lhs expression
	rhs expression
}

func parse() []expression {
	var ret []expression
	for len(tokens) > 0 {
		ret = append(ret, parseExpression())
	}
	return ret
}

func parseExpression() expression {
	return add()
}

func add() expression {
	ret := primary()
	for {
		switch {
		case consume("+"):
			ret = &binary{op: "+", lhs: ret, rhs: primary()}
		case consume("-"):
			ret = &binary{op: "-", lhs: ret, rhs: primary()}
		default:
			return ret
		}
	}
}

func primary() expression {
	return parseIntLit()
}

func parseIntLit() expression {

	tok := tokens[0]
	advance()

	return &intLit{
		val: tok.num,
	}
}
