package main

import "fmt"

func codegen(program []statement) {

	offset := 0
	for i := len(locals) - 1; i >= 0; i-- {
		offset += 8
		locals[i].offset = offset
	}

	fmt.Printf(".intel_syntax noprefix\n")

	fmt.Printf("\t.text\n")
	fmt.Printf("\t.globl main\n")
	fmt.Printf("main:\n")
	fmt.Printf("\tpush rbp\n")
	fmt.Printf("\tmov rbp, rsp\n")
	fmt.Printf("\tsub rsp, %d\n", offset)

	for _, c := range program {
		genStmt(c)
	}

	fmt.Printf(".Lreturn.main:\n")
	fmt.Printf("\tmov rsp, rbp\n")
	fmt.Printf("\tpop rbp\n")
	fmt.Printf("\tret\n")
}

func genStmt(stmt statement) {
	switch s := stmt.(type) {
	case *returnStmt:
		if s.child != nil {
			genExpr(s.child)
			fmt.Printf("\tpop rax\n")
		}
		fmt.Printf("\tjmp .Lreturn.main\n")
	case *blockStmt:
		for _, s := range s.stmts {
			genStmt(s)
		}
	case *expressionStmt:
		genExpr(s.child)
	case *assignment:
		genAddr(s.lhs)
		genExpr(s.rhs)
		store()
	}
}

func genExpr(expr expression) {
	switch e := expr.(type) {
	case *expressionList:
		genExpr(e.first)
		for _, e := range e.remain {
			genExpr(e)
		}
	case *intLit:
		fmt.Printf("\tpush %d\n", e.val)
	case *obj:
		genAddr(e)
		load()
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

func load() {
	fmt.Printf("\tpop rax\n")
	fmt.Printf("\tmov rax, [rax]\n")
	fmt.Printf("\tpush rax\n")
}

func store() {
	fmt.Printf("\tpop rdi\n")
	fmt.Printf("\tpop rax\n")
	fmt.Printf("\tmov [rax], rdi\n")
}

func genAddr(expr expression) {
	switch e := expr.(type) {
	case *obj:
		fmt.Printf("\tlea rax, [rbp-%d]\n", e.offset)
		fmt.Printf("\tpush rax\n")
	default:
		panic("not a value")
	}
}
