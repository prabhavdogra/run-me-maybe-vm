package main

import "fmt"

type Machine struct {
	stack        []Literal
	instructions []Instruction
}

type Instruction struct {
	operator        uint8
	instructionType InstructionSet
	value           Literal
	length          int
	line            int
	fileName        string
}

func (i Instruction) Error(message string) string {
	return fmt.Sprintf("ERROR (%s:%d): %s", i.fileName, i.line, message)
}

const maxStackSize = 1024

var debugMode = false
