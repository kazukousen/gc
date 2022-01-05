#!/bin/bash

assert() {
  expected="$1"
  input="$2"

  ./gc "${input}" > tmp.s || exit
  if [ "$?" != 0 ]; then
    echo "$input => compile failed"
    exit 1
  fi
  cc -o tmp tmp.s
  ./tmp
  actual="$?"

  if [ "$actual" = "$expected" ]; then
    echo "$input => $actual" "OK!"
  else
    echo "$input => $expected expected, but got $actual"
    exit 1
  fi
}

assert 0 '0'
assert 42 '42'
assert 48 '42 + 3+ 4-1'
assert 12 '5 * 6 / 2 + 5 - 8'
assert 10 '5 * -6 / 2 + -5 + 30'
assert 160 '5 * (-6 / (2 + -5) + 30)'
assert 1 '100 > 50'
assert 0 '100 < 50'
assert 1 '100 == 100'
assert 0 '100 == 50'
assert 1 '100 != 50'

echo OK
