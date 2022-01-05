package main

import "fmt"

func codegen(program []expression) {

	fmt.Printf(".intel_syntax noprefix\n")

	fmt.Printf("\t.text\n")
	fmt.Printf("\t.globl main\n")
	fmt.Printf("main:\n")

	for _, c := range program {
		genExpr(c)
		fmt.Printf("\tpop rax\n")
	}

	fmt.Printf("\tret\n")
}

func genExpr(expr expression) {
	switch e := expr.(type) {
	case *intLit:
		fmt.Printf("\tmov rax, %d\n", e.val)
		fmt.Printf("\tpush rax\n")
		return
	}
}
