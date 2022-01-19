package main

import (
	"fmt"
)

var funcName string

func codegen(prog *program) {

	fmt.Printf(".intel_syntax noprefix\n")

	fmt.Printf("\t.text\n")
	fmt.Printf("\tjmp main\n")

	for _, f := range prog.funcs {

		funcName = f.name

		fmt.Printf("\t.globl %s\n", funcName)
		fmt.Printf("%s:\n", funcName)
		fmt.Printf("\tpush rbp\n")
		fmt.Printf("\tmov rbp, rsp\n")

		fmt.Printf("\tsub rsp, %d\n", f.stackSize)

		genStmt(f.body)

		fmt.Printf(".Lreturn.%s:\n", funcName)
		fmt.Printf("\tmov rsp, rbp\n")
		fmt.Printf("\tpop rbp\n")
		fmt.Printf("\tret\n")
	}
}

var labelCnt = 0

func genStmt(stmt statement) {
	switch s := stmt.(type) {
	case *returnStmt:
		if s.child != nil {
			genStmt(s.child)
		}
		fmt.Printf("\tjmp .Lreturn.%s\n", funcName)
	case *blockStmt:
		for _, s := range s.stmts {
			genStmt(s)
		}
	case *ifStmt:
		labelCnt++
		cnt := labelCnt

		if s.init != nil {
			genStmt(s.init)
		}

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
		for i := range s.rhs {
			genExpr(s.rhs[i])
		}
		for i := len(s.lhs) - 1; i >= 0; i-- {
			genAddr(s.lhs[i])
			store()
		}
	}
}

func genExpr(expr expression) {
	switch e := expr.(type) {
	case *funcCall:
		fmt.Printf("\tsub rsp, %d\n", e.target.resultsSize)
		for i := len(e.args) - 1; i >= 0; i-- {
			genExpr(e.args[i])
		}
		fmt.Printf("\tcall %s\n", e.name)
		fmt.Printf("\tadd rsp, %d\n", e.target.paramsSize)
	case *intLit:
		fmt.Printf("\tpush %d\n", e.val)
	case *obj:
		genAddr(e)
		load(e.ty)
	case *deref:
		genExpr(e.child)
		load(e.ty)
	case *addr:
		genAddr(e.child)
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

func load(ty *typ) {
	fmt.Printf("\tpop rax\n")
	if ty.size == 1 {
		fmt.Printf("\tmovzx rax, byte ptr [rax]\n")
	} else {
		fmt.Printf("\tmov rax, [rax]\n")
	}
	fmt.Printf("\tpush rax\n")
}

func store() {
	fmt.Printf("\tpop rdi\n")
	fmt.Printf("\tpop rax\n")
	fmt.Printf("\tmov [rdi], rax\n")
}

func genAddr(expr expression) {
	switch e := expr.(type) {
	case *obj:
		fmt.Printf("\tlea rax, [rbp%+d]\n", e.offset)
		fmt.Printf("\tpush rax\n")
	default:
		panic("not a value")
	}
}
