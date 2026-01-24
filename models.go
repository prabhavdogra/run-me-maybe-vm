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
	heap            map[int64][]Literal
	heapPtr         int64
	input           io.Reader
	output          io.Writer
	fileDescriptors map[int64]*os.File
	stringTable     []int64
}

type RuntimeContext struct {
	*Machine
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

var debugMode = false
