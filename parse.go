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

func consumeToken(tk tokenKind) *token {
	if len(tokens) > 0 && tokens[0].kind == tk {
		tok := tokens[0]
		advance()
		return tok
	}
	return nil
}

type program struct {
	funcs []*function
}

type function struct {
	name      string
	body      statement
	params    []*obj
	locals    []*obj
	stackSize int
}

func (f *function) assignLVarOffsets() {
	offset := 0
	for i := len(f.locals) - 1; i >= 0; i-- {
		offset += 8
		f.locals[i].offset = offset
	}
	f.stackSize = offset
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

type deref struct {
	expression
	child expression
}

type addr struct {
	expression
	child expression
}

var locals []*obj

func createLocalVar(name string) *obj {
	lv := &obj{
		name: name,
	}
	locals = append(locals, lv)
	return lv
}

func findLocalVar(name string) *obj {

	for i := range locals {
		lv := locals[i]
		if lv.name == name {
			return lv
		}
	}

	// TODO: remove creating
	return createLocalVar(name)
}

// TopLevelDecl = FunctionDecl .
func parse() *program {
	ret := &program{
		funcs: make([]*function, 0),
	}
	for len(tokens) > 0 {
		expect("func")
		ret.funcs = append(ret.funcs, parseFunction())
		expect(";")
	}
	return ret
}

// FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
// FunctionName = identifier .
// FunctionBody = Block .
func parseFunction() *function {

	locals = []*obj{}

	tok := consumeToken(tokenKindIdentifier)
	if tok == nil {
		panic("must be an identifier")
	}

	ret := &function{name: tok.val}

	expect("(")
	// Signature = Parameters [ Type ] .
	ret.params = parseParameters()
	// TODO: support void
	parseType()

	if !consume("{") {
		return ret
	}

	ret.body = parseBlockStmt()
	ret.locals = locals

	ret.assignLVarOffsets()

	return ret
}

// Parameters = "(" [ ParameterList ] ")" .
func parseParameters() []*obj {

	var params []*obj

	if consume(")") {
		return params
	}

	params = parseParameterList()
	expect(")")

	return params
}

// ParameterList  = ParameterDecl { "," ParameterDecl } .
func parseParameterList() []*obj {
	var ret []*obj
	for i := 0; i == 0 || consume(","); i++ {
		ret = append(ret, parseParameterDecl()...)
	}
	return ret
}

// ParameterDecl  = [ IdentifierList ] Type .
func parseParameterDecl() []*obj {
	ids := parseIdentifierList()
	parseType()
	return ids
}

func parseType() {
	expect("int")
}

// Statement = ReturnStmt | SimpleStmt .
func parseStatement() statement {

	// return
	if consume("return") {
		// ReturnStmt = "return" [ ExpressionList ] .
		if peek("}") {
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

func parseIdentifierList() []*obj {
	var ret []*obj
	tok := consumeToken(tokenKindIdentifier)
	if tok == nil {
		return ret
	}
	ret = append(ret, createLocalVar(tok.val))

	for consume(",") {
		tok := consumeToken(tokenKindIdentifier)
		if tok == nil {
			panic(fmt.Sprintf("Expect an identifier: %+v", tokens[0]))
		}
		ret = append(ret, createLocalVar(tok.val))
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

// unary = ("+" | "-" | "*" | "&")? unary | primary
func parseUnary() expression {
	switch {
	case consume("+"):
		return parseUnary()
	case consume("-"):
		return &binary{op: "-", lhs: &intLit{val: 0}, rhs: parseUnary()}
	case consume("*"):
		return &deref{child: parseUnary()}
	case consume("&"):
		return &addr{child: parseUnary()}
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
	if tok := consumeToken(tokenKindIdentifier); tok != nil {
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
