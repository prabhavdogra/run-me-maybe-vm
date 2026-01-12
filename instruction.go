package main

import "fmt"

type InstructionSet uint8

const (
	InstructionPush InstructionSet = iota
	InstructionPop
	InstructionDup
	InstructionSwap
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
	case InstructionDup:
		return "DUP"
	case InstructionSwap:
		return "SWAP"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", i)
	}
}

func runInstructions(machine *Machine) {
	for insPtr := 0; insPtr < len(machine.instructions); insPtr++ {
		instr := machine.instructions[insPtr]
		switch instr.instructionType {
		case InstructionPush:
			push(machine, instr.value)
		case InstructionPop:
			pop(machine)
		case InstructionDup:
			x := pop(machine)
			push(machine, x)
			push(machine, x)
		case InstructionSwap:
			a := pop(machine)
			b := pop(machine)
			push(machine, a)
			push(machine, b)
		case InstructionAdd:
			a := pop(machine)
			b := pop(machine)
			push(machine, a+b)
		case InstructionSub:
			a := pop(machine)
			b := pop(machine)
			push(machine, a-b)
		case InstructionMul:
			a := pop(machine)
			b := pop(machine)
			push(machine, a*b)
		case InstructionDiv:
			a := pop(machine)
			b := pop(machine)
			push(machine, a/b)
			if b == 0 {
				panic("ERROR: division by zero")
			}
		case InstructionPrint:
			x := pop(machine)
			fmt.Println(x)
		default:
			panic("ERROR: unknown instruction")
		}
	}
}
