package main

import (
	"strings"
)

var in string
var tokens []*token

type tokenKind int

const (
	// https://go.dev/ref/spec#Tokens
	// There are four classes: identifiers, keywords, operators and punctuation, and literals
	tokenKindLiteral tokenKind = iota
	tokenKindOperator
	tokenKindKeyword
	tokenKindIdentifier
	tokenKindType
)

type token struct {
	kind tokenKind
	val  string
	num  int // for int
}

func tokenize() {
	for len(in) > 0 {

		if in[0] == ' ' {
			in = in[1:]
			continue
		}

		if in[0] == '\n' {
			in = in[1:]
			autoInsertSemicolon()
			continue
		}

		if in[0] == ';' {
			in = in[1:]
			addSemicolonToken()
			continue
		}

		if strings.Contains("+-*/()=<>!,{}&:.[]", in[0:1]) {
			if len(in) > 1 && (in[0:2] == "<=" || in[0:2] == ">=" || in[0:2] == "==" || in[0:2] == "!=" || in[0:2] == ":=") {
				tokens = append(tokens, &token{kind: tokenKindOperator, val: in[0:2]})
				in = in[2:]
			} else {
				tokens = append(tokens, &token{kind: tokenKindOperator, val: in[0:1]})
				in = in[1:]
			}
			continue
		}

		if in[0] >= 'a' && in[0] <= 'z' {
			name := in[0:1]
			in = in[1:]
			for len(in) > 0 && (isAlpha() || isDigit()) {
				name += in[0:1]
				in = in[1:]
			}
			tokens = append(tokens, identifierToken(name))
			continue
		}

		if isDigit() {
			tokens = append(tokens, &token{kind: tokenKindLiteral, num: toInt()})
			continue
		}

		panic("unexpected character: " + string(in[0]))
	}

	autoInsertSemicolon()
}

func isDigit() bool {
	return in[0] >= '0' && in[0] <= '9'
}

func isAlpha() bool {
	return (in[0] >= 'a' && in[0] <= 'z') || (in[0] >= 'A' && in[0] <= 'Z')
}

func toInt() int {
	var ret int
	for len(in) > 0 && (in[0] >= '0' && in[0] <= '9') {
		ret = ret*10 + int(in[0]-'0')
		in = in[1:]
	}
	return ret
}

func identifierToken(val string) *token {
	if inKeywords(val) {
		return &token{kind: tokenKindKeyword, val: val}
	}
	if inTypes(val) {
		return &token{kind: tokenKindType, val: val}
	}
	return &token{kind: tokenKindIdentifier, val: val}
}

func inKeywords(val string) bool {
	_, ok := map[string]struct{}{
		"return": {},
		"for":    {},
		"if":     {},
		"else":   {},
	}[val]
	return ok
}

func inTypes(val string) bool {
	_, ok := map[string]struct{}{
		"int":    {},
		"bool":   {},
		"struct": {},
	}[val]
	return ok
}

// a semicolon is automatically inserted into the token stream immediately
// after a line's final token if that token is
//
// * an identifier
// * an integer, floating*point, imaginary, rune, or string literal
// * one of the keywords break, continue, fallthrough, or return
// * one of the operators and punctuation ++, **, ), ], or }
func autoInsertSemicolon() {

	needed := func() bool {

		finalTok := tokens[len(tokens)-1]

		if finalTok.kind == tokenKindLiteral || finalTok.kind == tokenKindIdentifier {
			return true
		}

		if finalTok.kind == tokenKindKeyword &&
			(finalTok.val == "return") {
			return true
		}

		if finalTok.kind == tokenKindOperator &&
			strings.Contains(")}", finalTok.val) {
			return true
		}

		return false
	}()

	if needed {
		addSemicolonToken()
	}
}

func addSemicolonToken() {
	tokens = append(tokens, &token{kind: tokenKindOperator, val: ";"})
}
