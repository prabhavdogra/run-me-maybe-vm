package main

import (
	"vm/internal/lexer"
	"vm/internal/parser"
)

func main() {
	args := parseArgs()
	lex := lexer.Init(args.FileName).Lex()
	if args.DebugMode {
		lex.Print()
	}
	parsedTokens := parser.Init(lex)
	if args.DebugMode {
		parsedTokens.Print()
	}
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
