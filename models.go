package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Machine struct {
	stack           []Literal
	instructions    []Instruction
	heap            []Literal
	allocations     map[int]int // ptr -> size, for safety checks
	input           io.Reader
	output          io.Writer
	fileDescriptors map[int64]*os.File
	stringTable     []int64
	entrypoint      int
	strStack        []int64 // Stack of pointers to heap
}

type RuntimeContext struct {
	*Machine
	returnStack        []int
	CurrentInstruction Instruction
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
	return fmt.Sprintf("ERROR (%s:%d): %s", filepath.Base(i.fileName), i.line, message)
}

const maxStackSize = 1024
const maxReturnStackSize = 1024
const maxStrStackSize = 1024

var debugMode = false
