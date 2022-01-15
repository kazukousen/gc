#!/bin/bash

cat <<EOF | gcc -xc -c -o tmp2.o -
int ret3() { return 3; }
int ret5() { return 5; }
int add(int x, int y) { return x + y; }
int sub(int x, int y) { return x - y; }
int add6(int a, int b, int c, int d, int e, int f) {
  return a+b+c+d+e+f;
}
EOF

assert() {
  expected="$1"
  input="$2"

  ./gc "${input}" > tmp.s || exit
  if [ "$?" != 0 ]; then
    echo "$input => compile failed"
    exit 1
  fi
  cc -o tmp tmp.s tmp2.o
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
assert 5 'a = 5; return a;'
assert 3 'foo=3; return foo;'
assert 8 'foo123=3; bar=5; return foo123+bar;'
echo ""

echo "blocks"
echo ""
assert 5 '{a = 5; return a}'
echo ""

echo "if"
echo ""
assert 5 '{ i = 5; if i == 5 { return i}; return 0; }'
assert 5 '{ i = 5; if i == 5 { return i} else { return 3} }'
assert 3 '{ i = 3; if i == 5 { return i} else { return 3} }'
assert 3 '{ i = 3; if i == 5 { return i} else if i == 3 { return 3} }'
assert 3 '{ i = 1; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i }'
assert 4 '{ i = 4; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i }'
echo ""

echo "for"
echo ""
assert 55 '{ i=0; j=0; for i=0; i<=10; i=i+1 { j=i+j }; return j; }'
assert 55 '{ i=0; j=0; for ; i<=10; i=i+1 { j=i+j }; return j; }'
assert 55 '{ i=0; j=0; for ; i<=10; { j=i+j; i=i+1 }; return j; }'
assert 55 '{ i=0; j=0; for i<=10 { j=i+j; i=i+1 }; return j; }'
assert 3 '{ for {return 3;} return 5; }'
echo ""

echo "function call"
echo ""
assert 3 '{return ret3()}'
assert 5 '{return ret5();}'
assert 8 'return add(3, 5);'
assert 3 'return sub(5, 2);'
assert 21 'return add6(1,2,3,4,5,6);'
echo ""

echo OK
