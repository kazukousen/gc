package main

import "fmt"

var argRegisters64 = []string{"rdi", "rsi", "rdx", "rcx", "r8", "r9"}

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

var labelCnt = 0

func genStmt(stmt statement) {
	switch s := stmt.(type) {
	case *returnStmt:
		for _, child := range s.children {
			genExpr(child)
		}
		fmt.Printf("\tpop rax\n")
		fmt.Printf("\tjmp .Lreturn.main\n")
	case *blockStmt:
		for _, s := range s.stmts {
			genStmt(s)
		}
	case *ifStmt:
		labelCnt++
		cnt := labelCnt
		genExpr(s.cond)
		fmt.Printf("\tpop rax\n")
		fmt.Printf("\tcmp rax, 0\n")

		if s.els != nil {
			fmt.Printf("\tje .Lelse%d\n", cnt)
			genStmt(s.then)
			fmt.Printf("\tjmp .Lend%d\n", cnt)
			fmt.Printf(".Lelse%d:\n", cnt)
			genStmt(s.els)
			fmt.Printf(".Lend%d:\n", cnt)
		} else {
			fmt.Printf("\tje .Lend%d\n", cnt)
			genStmt(s.then)
			fmt.Printf(".Lend%d:\n", cnt)
		}
	case *forStmt:
		labelCnt++
		cnt := labelCnt
		if s.init != nil {
			genStmt(s.init)
		}
		fmt.Printf(".Lbegin%d:\n", cnt)
		if s.cond != nil {
			genExpr(s.cond)
			fmt.Printf("\tpop rax\n")
			fmt.Printf("\tcmp rax, 0\n")
			fmt.Printf("\tje .Lend%d\n", cnt)
		}
		genStmt(s.body)
		if s.post != nil {
			genStmt(s.post)
		}
		fmt.Printf("\tjmp .Lbegin%d\n", cnt)
		fmt.Printf(".Lend%d:\n", cnt)
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
	case *funcCall:
		for _, arg := range e.args {
			genExpr(arg)
		}

		for i := len(e.args) - 1; i >= 0; i-- {
			fmt.Printf("\tpop %s\n", argRegisters64[i])
		}

		fmt.Printf("\tcall %s\n", e.name)
		fmt.Printf("\tpush rax\n")
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
