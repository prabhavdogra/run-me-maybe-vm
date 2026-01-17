package main

import (
	"vm/internal/lexer"
	"vm/internal/parser"
)

func main() {
	lex := lexer.Init("test.wm").Lex()
	// lex.Print()
	parsedTokens := parser.Init(lex)
	// parsedTokens.Print()
	instructions := generateInstructions(parsedTokens)
	instructions.Print()
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: instructions,
	}
	writeProgram(loadedMachine, "program.bin")
	loadedMachine = readProgram("program.bin")
	runInstructions(loadedMachine)
	printStack(loadedMachine)
}
