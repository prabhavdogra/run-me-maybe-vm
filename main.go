package main

import "vm/lexer"

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

func swap(machine *Machine, index int) {
	if index < 0 || index >= len(machine.stack) {
		panic("ERROR: swap index out of bounds")
	}
	topIndex := len(machine.stack) - 1
	machine.stack[topIndex], machine.stack[index] = machine.stack[index], machine.stack[topIndex]
}

func main() {
	lexer.Lex()
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: program,
	}
	writeProgram(loadedMachine, "program.bin")
	loadedMachine = readProgram("program.bin")
	runInstructions(loadedMachine)
	printStack(loadedMachine)
}
