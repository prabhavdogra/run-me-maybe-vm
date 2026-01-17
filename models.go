package main

type Machine struct {
	stack        []int64
	instructions []Instruction
}

type Instruction struct {
	operator        uint8
	instructionType InstructionSet
	value           int64
}

const maxStackSize = 1024

var debugMode = false
