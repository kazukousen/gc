package main

import (
	"fmt"
)

type typeKind int

const (
	typeKindInt typeKind = iota
	typeKindBool
	typeKindPtr
)

type typ struct {
	kind typeKind
	base *typ
	size int
}

func newType(kind typeKind, size int) *typ {
	return &typ{kind: kind, size: size}
}

var (
	typeKindMap = map[string]typeKind{
		"int":  typeKindInt,
		"bool": typeKindBool,
	}
	typeKindSize = map[string]int{
		"int":  8,
		"bool": 8,
	}
)

func newLiteralType(s string) *typ {
	return newType(typeKindMap[s], typeKindSize[s])
}

func pointerTo(base *typ) *typ {
	ty := newType(typeKindPtr, 8)
	ty.base = base
	return ty
}

func addType(n interface{}) {
	if n == nil {
		return
	}
	if n, ok := n.(interface{ getType() *typ }); ok && n.getType() != nil {
		return
	}

	switch n := n.(type) {
	case *returnStmt:
		addType(n.child)
		return
	case *blockStmt:
		for _, stmt := range n.stmts {
			addType(stmt)
		}
		return
	case *ifStmt:
		addType(n.init)
		addType(n.cond)
		addType(n.then)
		addType(n.els)
		return
	case *forStmt:
		addType(n.init)
		addType(n.cond)
		addType(n.post)
		addType(n.body)
		return
	case *expressionStmt:
		addType(n.child)
		return
	case *assignment:
		if se := n.rhs.singleMultiValuedExpression(); se != nil {
			addType(se)
			for i, e := range se.multiValues() {
				n.lhs[i].setType(e.getType())
			}
		} else {
			for i, e := range n.rhs {
				addType(e)
				n.lhs[i].setType(e.getType())
			}
		}
		return
	case *intLit:
		n.setType(newLiteralType("int"))
		return
	case *binary:
		addType(n.lhs)
		addType(n.rhs)
		switch n.op {
		case "+", "-", "*", "/":
			n.setType(n.lhs.getType())
		case "==", "!=", "<", "<=":
			n.setType(newLiteralType("bool"))
		}
		return
	case *obj:
		return
	case *deref:
		addType(n.child)
		n.setType(n.child.getType())
		return
	case *addr:
		addType(n.child)
		ct := n.child.getType()
		ty := pointerTo(ct)
		n.setType(ty)
	case *funcCall:
		for _, arg := range n.args {
			addType(arg)
		}
		if len(n.target.results) > 0 {
			n.setType(n.target.results[0].getType())
		}
		return
	default:
		panic(fmt.Sprintf("Unsupported type: %T", n))
	}
}
