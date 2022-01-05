package main

import "strings"

var in string
var tokens []token

type tokenKind int

const (
	tokenKindNumber tokenKind = iota
	tokenKindReserved
)

type token struct {
	kind tokenKind
	val  string
	num  int // for int
}

func isDigit() bool {
	return in[0] >= '0' && in[0] <= '9'
}

func toInt() int {
	var ret int
	for len(in) > 0 && (in[0] >= '0' && in[0] <= '9') {
		ret = ret*10 + int(in[0]-'0')
		in = in[1:]
	}
	return ret
}

func tokenize() {
	for len(in) > 0 {

		if in[0] == ' ' {
			in = in[1:]
			continue
		}

		if strings.Contains("+-*/()", in[0:1]) {
			tokens = append(tokens, token{kind: tokenKindReserved, val: in[0:1]})
			in = in[1:]
			continue
		}

		if isDigit() {
			tokens = append(tokens, token{kind: tokenKindNumber, num: toInt()})
			continue
		}

		panic("unexpected character: " + string(in[0]))
	}
}
