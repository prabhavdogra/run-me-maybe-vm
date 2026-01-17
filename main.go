package main

import (
	"vm/internal/lexer"
	"vm/internal/parser"
)

var program = []Instruction{
	pushIns(14),
	pushIns(14),
	pushIns(14),
	pushIns(14),
	pushIns(5),
	nzjmpIns(7),
	pushIns(15),
	noopIns(),
	noopIns(),
	noopIns(),
	printIns(),
}

func main() {
	lex := lexer.Init("test.wm").Lex()
	parser.Init(lex)
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: program,
	}
	writeProgram(loadedMachine, "program.bin")
	loadedMachine = readProgram("program.bin")
	runInstructions(loadedMachine)
	printStack(loadedMachine)
}
