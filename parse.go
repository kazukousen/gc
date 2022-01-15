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
	children []expression
}

type blockStmt struct {
	statement
	stmts []statement
}

type ifStmt struct {
	statement
	init statement
	cond expression
	then statement
	els  statement
}

type forStmt struct {
	statement
	cond expression
	init statement
	post statement

	body statement
}

type expressionStmt struct {
	statement
	child expression
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

type funcCall struct {
	expression
	name string
	args []expression
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

	// return
	if consume("return") {
		// ReturnStmt = "return" [ ExpressionList ] .
		if peek(";") {
			return &returnStmt{}
		}
		ret := parseExpressionList()
		return &returnStmt{children: ret}
	}

	// block
	if consume("{") {
		return parseBlockStmt()
	}

	// if
	if consume("if") {
		return parseIfStmt()
	}

	// for
	if consume("for") {
		return parseForStmt()
	}

	return parseSimpleStmt()
}

// Block = "{" StatementList "}" .
func parseBlockStmt() statement {
	var stmts []statement
	for !consume("}") {
		stmts = append(stmts, parseStatement())
	}
	return &blockStmt{stmts: stmts}
}

// IfStmt = "if" [ SimpleStmt ";" ] Expression Block [ "else" ( IfStmt | Block ) ] .
func parseIfStmt() statement {
	var cond expression
	var init statement
	tmp := parseSimpleStmt()
	if t, ok := tmp.(*expressionStmt); ok {
		cond = t.child
	} else {
		init = tmp
		cond = parseExpression()
	}

	expect("{")
	then := parseBlockStmt()

	ret := &ifStmt{
		init: init,
		cond: cond,
		then: then,
	}

	if !consume("else") {
		return ret
	}

	if consume("{") {
		ret.els = parseBlockStmt()
	} else if consume("if") {
		ret.els = parseIfStmt()
	}

	return ret
}

// ForStmt    = "for" [ Condition | ForClause ] Block .
// Condition  = Expression .
// ForClause  = [ InitStmt ] ";" [ Condition ] ";" [ PostStmt ] .
// InitStmt   = SimpleStmt .
// PostStmt   = SimpleStmt .
func parseForStmt() statement {
	if consume("{") {
		return &forStmt{body: parseBlockStmt()}
	}

	var cond expression
	var init statement
	var post statement
	if !consume(";") {

		tmp := parseSimpleStmt()
		if t, ok := tmp.(*expressionStmt); ok {
			cond = t.child
			expect("{")
			return &forStmt{
				cond: cond,
				body: parseBlockStmt(),
			}
		}

		init = tmp
		expect(";")
	}

	if !consume(";") {
		cond = parseExpression()
		expect(";")
	}

	if !consume("{") {
		post = parseSimpleStmt()
		expect("{")
	}

	return &forStmt{
		cond: cond,
		init: init,
		post: post,
		body: parseBlockStmt(),
	}
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
func parseExpressionList() []expression {

	ret := []expression{parseExpression()}

	for consume(",") {
		ret = append(ret, parseExpression())
	}

	return ret
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

// PrimaryExpr = Operand .
func parsePrimary() expression {

	expr := parseOperand()

	return expr
}

// Operand = Literal | identifier [ Arguments ] | "(" Expression ")" .
func parseOperand() expression {
	if consume("(") {
		ret := parseExpression()
		expect(")")
		return ret
	}

	// identifier
	if tok := consumeIdent(); tok != nil {
		lv := findLocalVar(tok.val)

		if consume("(") {
			return parseArguments(tok.val)
		}

		return lv
	}

	// Literal
	return parseIntLit()
}

// Arguments = "(" [ ExpressionList [ "..." ] [ "," ] ] ")" .
func parseArguments(name string) expression {

	ret := &funcCall{name: name}

	if consume(")") {
		return ret
	}

	ret.args = parseExpressionList()
	expect(")")

	return ret
}

func parseIntLit() expression {

	tok := tokens[0]
	advance()

	return &intLit{
		val: tok.num,
	}
}
