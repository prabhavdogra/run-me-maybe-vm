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
	if args.DebugMode {
		instructions.Print()
	}
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: instructions,
	}
	runInstructions(loadedMachine)
	writeProgram(loadedMachine, "program.bin")
	printStack(loadedMachine)
}
