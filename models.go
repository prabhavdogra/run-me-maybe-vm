package main

type Machine struct {
	stack        []int
	stackSize    int
	programSize  int
	instructions []Instruction
}

type Instruction struct {
	operator        uint8
	instructionType InstructionSet
	value           int64
}

var ProgramSize = len(program)

const maxStackSize = 1024

var stack [maxStackSize]int64
var stackSize = 0

var debugMode = false
