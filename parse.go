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
	return parseAdd()
}

// add = mul (("*" | "/") mul)*
func parseAdd() expression {
	ret := parseMul()
	for {
		switch {
		case consume("+"):
			ret = &binary{op: "+", lhs: ret, rhs: parseMul()}
		case consume("-"):
			ret = &binary{op: "-", lhs: ret, rhs: parseMul()}
		default:
			return ret
		}
	}
}

// mul = unary (("*" | "/") unary)*
func parseMul() expression {
	ret := parseUnary()
	for {
		switch {
		case consume("*"):
			ret = &binary{op: "*", lhs: ret, rhs: parseUnary()}
		case consume("/"):
			ret = &binary{op: "/", lhs: ret, rhs: parseUnary()}
		default:
			return ret
		}
	}
}

// unary = ("+" | "-")? unary | primary
func parseUnary() expression {
	switch {
	case consume("+"):
		return parseUnary()
	case consume("-"):
		return &binary{op: "-", lhs: &intLit{val: 0}, rhs: parseUnary()}
	default:
		return parsePrimary()
	}
}

func parsePrimary() expression {
	return parseIntLit()
}

func parseIntLit() expression {

	tok := tokens[0]
	advance()

	return &intLit{
		val: tok.num,
	}
}
