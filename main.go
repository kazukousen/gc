package main

import "os"

func main() {

	in = os.Args[1]

	tokenize()

	prog := parse()

	codegen(prog)
}
