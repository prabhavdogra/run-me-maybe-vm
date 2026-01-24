package main

import (
	"os"
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
		instructions: instructions,
		stack:        make([]Literal, 0, maxStackSize),
		heap:         make(map[int64][]Literal),
		input:        os.Stdin,
		output:       os.Stdout,
	}
	loadedMachine = runInstructions(loadedMachine)
	if debugMode {
		printStack(loadedMachine)
	}
	writeProgram(loadedMachine, "program.bin")
}
