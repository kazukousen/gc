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
	name       string
	body       statement
	params     []*obj
	results    []*obj
	locals     []*obj
	stackSize  int
	resultSize int
}

func (f *function) assignLVarOffsets() {
	offset := 8
	for i := len(f.params) - 1; i >= 0; i-- {
		offset += 8
		f.params[i].offset = offset
	}
	f.resultSize = offset
	for i := len(f.results) - 1; i >= 0; i-- {
		offset += 8
		f.results[i].offset = offset
	}
	f.resultSize = offset - f.resultSize

	offset = 0
	for i := len(f.locals) - 1; i >= 0; i-- {
		if f.locals[i].offset != 0 {
			continue
		}
		offset += 8
		f.locals[i].offset = -offset
	}
	f.stackSize = offset
}

// Statement

type statement interface {
	aStmt()
}

type returnStmt struct {
	statement
	child    statement
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
	lhs []expression
	rhs []expression
}

// Expressions

type expression interface {
	anExpr()
}

type funcCall struct {
	expression
	name       string
	args       []expression
	resultSize int
}

var callers []*funcCall

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

// temporary sets
var locals []*obj
var results []*obj

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
	return nil
}

// TopLevelDecl = FunctionDecl .
func parse() *program {
	mFuncs := make(map[string]*function)
	ret := &program{
		funcs: make([]*function, 0),
	}
	for len(tokens) > 0 {
		expect("func")
		f := parseFunction()
		ret.funcs = append(ret.funcs, f)
		mFuncs[f.name] = f
		expect(";")
	}

	for _, c := range callers {
		f := mFuncs[c.name]
		c.resultSize = f.resultSize
	}

	return ret
}

// VarDecl = "var" ( VarSpec | "(" { VarSpec ";" } ")" ) .
func parseVarDecl() statement {

	if !consume("(") {
		return parseVarSpec()
	}

	ret := &blockStmt{
		stmts: []statement{},
	}
	for !consume(")") {
		ret.stmts = append(ret.stmts, parseVarSpec())
		expect(";")
	}
	return ret
}

// VarSpec = IdentifierList ( Type [ "=" ExpressionList ] | "=" ExpressionList ) .
func parseVarSpec() statement {
	ids := parseIdentifierList()
	lhs := make([]expression, len(ids))
	for i, id := range ids {
		lhs[i] = id
	}
	if consume("=") {
		rhs := parseExpressionList()
		return &assignment{lhs: lhs, rhs: rhs}
	}

	parseType()
	rhs := make([]expression, len(ids))
	for i := range ids {
		rhs[i] = &intLit{
			val: 0,
		}
	}
	return &assignment{lhs: lhs, rhs: rhs}
}

// FunctionDecl = "func" FunctionName Signature [ FunctionBody ] .
// FunctionName = identifier .
// FunctionBody = Block .
func parseFunction() *function {

	locals = []*obj{}
	results = []*obj{}

	tok := consumeToken(tokenKindIdentifier)
	if tok == nil {
		panic("must be an identifier")
	}

	ret := &function{name: tok.val}

	expect("(")
	// Signature = Parameters [ Type ] .
	ret.params, ret.results = parseSignature()

	results = ret.results
	expect("{")
	ret.body = parseBlockStmt()
	ret.locals = locals

	ret.assignLVarOffsets()

	return ret
}

func parseSignature() ([]*obj, []*obj) {
	params := parseParameters()

	if consume("(") {
		// multiple results
		for i := 0; !consume(")"); i++ {
			if i > 0 {
				expect(",")
			}
			// TODO: identifier
			tok := consumeToken(tokenKindType)
			if tok == nil {
				panic(fmt.Sprintf("Expected a type: %+v", tokens[0]))
			}
			results = append(results, createLocalVar(tok.val))
		}

		return params, results
	}

	if tok := consumeToken(tokenKindType); tok != nil {
		// TODO: identifier
		results = append(results, createLocalVar(tok.val))
		return params, results
	}

	return params, results
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

// Statement = Declaration | ReturnStmt | SimpleStmt .
// Declaration = VarDecl .
func parseStatement() statement {

	// varDeclaration
	if consume("var") {
		return parseVarDecl()
	}

	// return
	if consume("return") {
		// ReturnStmt = "return" [ ExpressionList ] .
		if peek("}") {
			return &returnStmt{}
		}
		ret := parseExpressionList()
		if len(results) == 0 {
			return &returnStmt{children: ret}
		}
		res := make([]expression, len(results))
		for i, v := range results {
			res[i] = v
		}
		return &returnStmt{child: &assignment{lhs: res, rhs: ret}}
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
		expect(";")
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
	expr := parseExpressionList()

	if consume("=") {
		// Assignment
		for i, l := range expr {
			switch l := l.(type) {
			case *obj:
				lv := findLocalVar(l.name)
				if lv == nil {
					panic(fmt.Sprintf("a local variable not declared: %s", l.name))
				}
				expr[i] = lv
			}
		}
		return &assignment{lhs: expr, rhs: parseExpressionList()}
	}

	if consume(":=") {
		// ShortVarDecl
		for i, l := range expr {
			switch l := l.(type) {
			case *obj:
				expr[i] = createLocalVar(l.name)
			}
		}
		return &assignment{lhs: expr, rhs: parseExpressionList()}
	}

	return &expressionStmt{child: expr[0]}
}

// ExpressionList = Expression { "," Expression } .
func parseExpressionList() []expression {

	ret := []expression{parseExpression()}

	for consume(",") {
		ret = append(ret, parseExpression())
	}

	return ret
}

// IdentifierList = identifier { "," identifier } .
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

		if consume("(") {
			return parseArguments(tok.val)
		}

		lv := findLocalVar(tok.val)
		if lv == nil {
			return &obj{name: tok.val}
		}

		return lv
	}

	// Literal
	return parseIntLit()
}

// Arguments = "(" [ ExpressionList [ "..." ] [ "," ] ] ")" .
func parseArguments(name string) expression {

	ret := &funcCall{name: name}

	callers = append(callers, ret)

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
