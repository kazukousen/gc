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

echo OK
