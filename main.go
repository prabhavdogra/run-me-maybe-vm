package main

import (
	"vm/internal/lexer"
	"vm/internal/parser"
)

func main() {
	lex := lexer.Init("test.wm").Lex()
	// lex.Print()
	parsedTokens := parser.Init(lex)
	parsedTokens.Print()
	instruction := generateInstructions(parsedTokens)
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: instruction,
	}
	writeProgram(loadedMachine, "program.bin")
	loadedMachine = readProgram("program.bin")
	runInstructions(loadedMachine)
	printStack(loadedMachine)
}
