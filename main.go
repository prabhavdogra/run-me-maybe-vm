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
	// preprocess strings into Heap
	stringTable, heap := populateStringTable(parsedTokens)
	loadedMachine := &Machine{
		stack:           []Literal{},
		instructions:    instructions,
		heap:            heap,
		allocations:     make(map[int]int),
		input:           os.Stdin,
		output:          os.Stdout,
		fileDescriptors: make(map[int64]*os.File),
		stringTable:     stringTable,
	}

	loadedMachine = runInstructions(loadedMachine)
	if debugMode {
		printStack(loadedMachine)
	}
	writeProgram(loadedMachine, "program.bin")
}
