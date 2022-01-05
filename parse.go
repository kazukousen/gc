package main

import (
	"fmt"
)

func advance() {
	tokens = tokens[1:]
}

func peek(s string) bool {
	return len(tokens) > 0 && tokens[0].val == s
}

func consume(s string) bool {
	if peek(s) {
		advance()
		return true
	}
	return false
}

func expect(s string) {
	if consume(s) {
		return
	}
	panic(fmt.Sprintf("Unexpected token: %+v. want: %s", tokens[0], s))
}

// Statement

type statement interface {
	aStmt()
}

type returnStmt struct {
	statement
	child expression
}

type expressionStmt struct {
	statement
	child expression
}

type expressionList struct {
	expression
	first  expression
	remain []expression
}

// Expressions

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

func parse() []statement {
	var ret []statement
	for len(tokens) > 0 {
		ret = append(ret, parseStatement())
		expect(";")
	}
	return ret
}

// Statement = ReturnStmt | ExpressionStmt .
func parseStatement() statement {

	if consume("return") {
		// ReturnStmt = "return" [ ExpressionList ] .
		if peek(";") {
			return &returnStmt{}
		}
		ret := parseExpressionList()
		return &returnStmt{child: ret}
	}

	// ExpressionStmt = Expression .
	ret := parseExpression()
	return &expressionStmt{child: ret}
}

// ExpressionList = Expression { "," Expression } .
func parseExpressionList() expression {
	first := parseExpression()

	var remain []expression
	for consume(",") {
		remain = append(remain, parseExpression())
	}

	if len(remain) == 0 {
		return first
	}

	return &expressionList{first: first, remain: remain}
}

func parseExpression() expression {
	return parseRel()
}

// rel = add ((">=" | "<=" | "==" | "!=") add)*
func parseRel() expression {
	ret := parseAdd()
	for {
		switch {

		case consume("<"):
			ret = &binary{op: "<", lhs: ret, rhs: parseAdd()}
		case consume(">"):
			ret = &binary{op: "<", lhs: parseAdd(), rhs: ret}
		case consume("<="):
			ret = &binary{op: "<=", lhs: ret, rhs: parseAdd()}
		case consume(">="):
			ret = &binary{op: "<=", lhs: parseAdd(), rhs: ret}
		case consume("=="):
			ret = &binary{op: "==", lhs: ret, rhs: parseAdd()}
		case consume("!="):
			ret = &binary{op: "!=", lhs: ret, rhs: parseAdd()}
		default:
			return ret
		}
	}

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

// primary = "(" expression ")" | num
func parsePrimary() expression {

	if consume("(") {
		ret := parseExpression()
		expect(")")
		return ret
	}

	return parseIntLit()
}

func parseIntLit() expression {

	tok := tokens[0]
	advance()

	return &intLit{
		val: tok.num,
	}
}
