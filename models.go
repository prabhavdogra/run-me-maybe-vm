package main

type Machine struct {
	stack        []Literal
	instructions []Instruction
}

type Instruction struct {
	operator        uint8
	instructionType InstructionSet
	value           Literal
}

const maxStackSize = 1024

var debugMode = false
