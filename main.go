package main

import (
	"fmt"
)

type InstructionSet int

const (
	InstructionPush InstructionSet = iota
	InstructionPop
	InstructionAdd
	InstructionSub
	InstructionMul
	InstructionDiv
	InstructionPrint
)

func (i InstructionSet) String() string {
	switch i {
	case InstructionPush:
		return "PUSH"
	case InstructionPop:
		return "POP"
	case InstructionAdd:
		return "ADD"
	case InstructionSub:
		return "SUB"
	case InstructionMul:
		return "MUL"
	case InstructionDiv:
		return "DIV"
	case InstructionPrint:
		return "PRINT"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", i)
	}
}

type Instruction struct {
	operator        int
	instructionType InstructionSet
	value           int
}

var program = []Instruction{
	pushIns(5),
	pushIns(10),
	pushIns(10),
	pushIns(2),
	mulIns(),
	printIns(),
}

var ProgramSize = len(program)

const maxStackSize = 1024

var stack [maxStackSize]int
var stackSize = 0

var debugMode = false

func push(value int) {
	if stackSize >= maxStackSize {
		panic("ERROR: stack overflow")
	}
	stack[stackSize] = value
	stackSize++
}

func pop() int {
	if stackSize == 0 {
		panic("ERROR: stack underflow")
	}
	stackSize--
	return stack[stackSize]
}

func printStack() {
	fmt.Println("Stack contents:")
	for i := 0; i < stackSize+1; i++ {
		fmt.Printf("[%d]: %d\n", i, stack[i])
	}
}

func main() {
	for insPtr := 0; insPtr < ProgramSize; insPtr++ {
		instr := program[insPtr]
		switch instr.instructionType {
		case InstructionPush:
			push(instr.value)
			if debugMode {
				fmt.Println("Pushed", instr.value)
			}
		case InstructionPop:
			x := pop()
			if debugMode {
				fmt.Println("Popped", x)
			}
		case InstructionAdd:
			a := pop()
			b := pop()
			push(a + b)
			if debugMode {
				fmt.Println("Added", a, "+", b)
			}
		case InstructionSub:
			a := pop()
			b := pop()
			push(a - b)
			if debugMode {
				fmt.Println("Subtracted", a, "-", b)
			}
		case InstructionMul:
			a := pop()
			b := pop()
			push(a * b)
			if debugMode {
				fmt.Println("Multiplied", a, "*", b)
			}
		case InstructionDiv:
			a := pop()
			b := pop()
			push(a / b)
			if b == 0 {
				panic("ERROR: division by zero")
			}
			if debugMode {
				fmt.Println("Divided", a, "/", b)
			}
		case InstructionPrint:
			x := pop()
			fmt.Println(x)
			if debugMode {
				fmt.Println("Printed top of stack", x)
			}
		default:
			panic("ERROR: unknown instruction")
		}
	}
	printStack()
}
