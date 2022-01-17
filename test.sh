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
assert 5 'func main() int { var a = 5; return a; }'
assert 3 'func main() int { var foo=3; return foo; }'
assert 8 'func main() int { var foo123=3; var bar=5; return foo123+bar; }'
assert 8 'func main() int { var foo123, bar = 3, 5; return foo123+bar; }'
echo ""

echo "blocks"
echo ""
assert 5 'func main() int { {var a = 5; return a} }'
echo ""

echo "if"
echo ""
assert 5 'func main() int { var i = 5; if i == 5 { return i}; return 0;  }'
assert 5 'func main() int { if i := 5; i == 5 { return i}; return 0; }'
assert 5 'func main() int { var i = 5; if i == 5 { return i} else { return 3}  }'
assert 3 'func main() int { var i = 3; if i == 5 { return i} else { return 3}  }'
assert 3 'func main() int { var i = 3; if i == 5 { return i} else if i == 3 { return 3}  }'
assert 3 'func main() int { i := 1; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i }'
assert 4 'func main() int { i := 4; if i == 5 { return i} else if i == 4 { return 4} else { i = 3 } return i }'
echo ""

echo "for"
echo ""
assert 55 'func main() int { i:=0; j:=0; for i=0; i<=10; i=i+1 { j=i+j }; return j; }'
assert 55 'func main() int { i:=0; j:=0; for ; i<=10; i=i+1 { j=i+j }; return j; }'
assert 55 'func main() int { i:=0; j := 0; for ; i<=10; { j=i+j; i=i+1 }; return j; }'
assert 55 'func main() int { i := 0; j := 0; for i<=10 { j=i+j; i=i+1 }; return j; }'
assert 3 'func main() int { for {return 3;} return 5; }'
echo ""

echo "pointer"
echo ""
assert 3 'func main() int { { var x=3; return *&x; } }'
assert 3 'func main() int { { x := 3; var y = &x; z := &y; return **z; } }'
echo ""

echo "function"
echo ""
assert 32 'func main() int { return ret32() }; func ret32() int { return 32; }'
assert 7 'func main() int { return add2(3,4); }; func add2(x int, y int) int { return x+y }'
assert 1 'func main() int { return sub2(4,3); }; func sub2(x int, y int) int { return x-y }'
assert 55 'func main() int { return fib(9); }; func fib(x int) int { if x <= 1 { return 1}; return fib(x-1) + fib(x-2) }'
assert 5 'func main() int {return myFunction(1, myFunction(1, 3))}; func myFunction(a, b int) int {return a + b}'
assert 13 'func myFunction(a, b int) int {return a + b}; func main() int { a := 6; return myFunction(a, 7)}'
echo ""

echo "function"
echo ""
assert 35 'func main() int {a, b := myFunction(3, 4); return a * b}; func myFunction(x, y int) (int, int) { lvar := 5; return x + y, 5 }'
echo ""

echo OK
