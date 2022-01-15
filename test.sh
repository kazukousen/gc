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

assert 0 'func main() int {return 0}'
assert 42 'func main() int {return 42}'
assert 48 'func main() int{return 42 + 3+ 4-1}'
assert 12 'func main() int {return 5 * 6 / 2 + 5 - 8}'
assert 10 'func main() int {return 5 * -6 / 2 + -5 + 30}'
assert 160 'func main() int {return 5 * (-6 / (2 + -5) + 30)}'
assert 1 'func main() int{return 100 > 50}'
assert 0 'func main() int {return 100 < 50}'
assert 1 'func main() int {return 100 == 100}'
assert 0 'func main() int { return 100 == 50 }'
assert 1 'func main() int { return 100 != 50 }'

assert 2 'func main() int { 1; return 2; }'
assert 1 'func main() int { return 1; 2; }'
assert 3 'func main() int { return 1, 3; }'
assert 2 'func main() int { return (1 + 3) / 2 }'

echo "local variables"
echo ""
assert 5 'func main() int { a = 5; return a; }'
assert 3 'func main() int { foo=3; return foo; }'
assert 8 'func main() int { foo123=3; bar=5; return foo123+bar; }'
echo ""

echo "blocks"
echo ""
assert 5 'func main() int { {a = 5; return a} }'
echo ""

echo "if"
echo ""
assert 5 'func main() int { { i = 5; if i == 5 { return i}; return 0; } }'
assert 5 'func main() int { { i = 5; if i == 5 { return i} else { return 3} } }'
assert 3 'func main() int { { i = 3; if i == 5 { return i} else { return 3} } }'
assert 3 'func main() int { { i = 3; if i == 5 { return i} else if i == 3 { return 3} } }'
assert 3 'func main() int { { i = 1; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i } }'
assert 4 'func main() int { { i = 4; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i } }'
echo ""

echo "for"
echo ""
assert 55 'func main() int { { i=0; j=0; for i=0; i<=10; i=i+1 { j=i+j }; return j; } }'
assert 55 'func main() int { { i=0; j=0; for ; i<=10; i=i+1 { j=i+j }; return j; } }'
assert 55 'func main() int { { i=0; j=0; for ; i<=10; { j=i+j; i=i+1 }; return j; } }'
assert 55 'func main() int { { i=0; j=0; for i<=10 { j=i+j; i=i+1 }; return j; } }'
assert 3 'func main() int { { for {return 3;} return 5; } }'
echo ""

echo "function call"
echo ""
assert 3 'func main() int { {return ret3()} }'
assert 5 'func main() int { {return ret5();} }'
assert 8 'func main() int { return add(3, 5); }'
assert 3 'func main() int { return sub(5, 2); }'
assert 21 'func main() int { return add6(1,2,3,4,5,6); }'
echo ""

echo "pointer"
echo ""
assert 3 'func main() int { { x=3; return *&x; } }'
assert 3 'func main() int { { x=3; y=&x; z=&y; return **z; } }'
echo ""

echo "function"
echo ""
assert 32 'func main() int { return ret32() }; func ret32() int { return 32; }'
echo ""

echo OK
