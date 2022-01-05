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
	case *binary:
		genExpr(e.lhs)
		genExpr(e.rhs)
		fmt.Printf("\tpop rdi\n")
		fmt.Printf("\tpop rax\n")
		switch e.op {
		case "+":
			fmt.Printf("\tadd rax, rdi\n")
		case "-":
			fmt.Printf("\tsub rax, rdi\n")
		case "*":
			fmt.Printf("\timul rax, rdi\n")
		case "/":
			fmt.Printf("\tcqo\n")
			fmt.Printf("\tidiv rdi\n")
		case "<":
			fmt.Printf("\tcmp rax, rdi\n")
			fmt.Printf("\tsetl al\n")
			fmt.Printf("\tmovzb rax, al\n")
		case "<=":
			fmt.Printf("\tcmp rax, rdi\n")
			fmt.Printf("\tsetle al\n")
			fmt.Printf("\tmovzb rax, al\n")
		case "==":
			fmt.Printf("\tcmp rax, rdi\n")
			fmt.Printf("\tsete al\n")
			fmt.Printf("\tmovzb rax, al\n")
		case "!=":
			fmt.Printf("\tcmp rax, rdi\n")
			fmt.Printf("\tsetne al\n")
			fmt.Printf("\tmovzb rax, al\n")
		}
		fmt.Printf("\tpush rax\n")
		return
	}
}
