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

func consumeIdent() *token {
	if len(tokens) > 0 && tokens[0].kind == tokenKindIdentifier {
		tok := tokens[0]
		advance()
		return tok
	}
	return nil
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

type assignment struct {
	statement
	lhs expression
	rhs expression
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

type obj struct {
	expression
	name   string
	offset int
}

var locals []*obj

func findLocalVar(name string) *obj {

	// TODO: remove creating

	for i := range locals {
		lv := locals[i]
		if lv.name == name {
			return lv
		}
	}

	lv := &obj{
		name: name,
	}
	locals = append(locals, lv)
	return lv
}

func parse() []statement {
	var ret []statement
	for len(tokens) > 0 {
		ret = append(ret, parseStatement())
		expect(";")
	}
	return ret
}

// Statement = ReturnStmt | SimpleStmt .
func parseStatement() statement {

	if consume("return") {
		// ReturnStmt = "return" [ ExpressionList ] .
		if peek(";") {
			return &returnStmt{}
		}
		ret := parseExpressionList()
		return &returnStmt{child: ret}
	}

	return parseSimpleStmt()
}

func parseSimpleStmt() statement {
	expr := parseExpression()

	if consume("=") {
		// Assignment
		return &assignment{lhs: expr, rhs: parseExpression()}
	}

	return &expressionStmt{child: expr}
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

// primary = "(" expression ")" | ident | num
func parsePrimary() expression {

	if consume("(") {
		ret := parseExpression()
		expect(")")
		return ret
	}

	if tok := consumeIdent(); tok != nil {
		lv := findLocalVar(tok.val)
		return lv
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
