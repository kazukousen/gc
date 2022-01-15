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

assert 0 'return 0'
assert 42 'return 42'
assert 48 'return 42 + 3+ 4-1'
assert 12 'return 5 * 6 / 2 + 5 - 8'
assert 10 'return 5 * -6 / 2 + -5 + 30'
assert 160 'return 5 * (-6 / (2 + -5) + 30)'
assert 1 'return 100 > 50'
assert 0 'return 100 < 50'
assert 1 'return 100 == 100'
assert 0 'return 100 == 50'
assert 1 'return 100 != 50'

assert 2 '1; return 2;'
assert 1 'return 1; 2;'
assert 3 'return 1, 3;'
assert 2 '(1 + 3) / 2; return'

echo "local variables"
echo ""
echo ""
assert 5 'a = 5; return a;'

echo OK
