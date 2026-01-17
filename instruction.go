package main

import "fmt"

type InstructionSet uint8

const (
	InstructionNoOp InstructionSet = iota
	InstructionPush
	InstructionPop
	InstructionDup
	InstructionInDup
	InstructionSwap
	InstructionInSwap
	InstructionAdd
	InstructionSub
	InstructionMul
	InstructionDiv
	InstructionMod
	InstructionCmpe
	InstructionCmpne
	InstructionCmpg
	InstructionCmpl
	InstructionCmpge
	InstructionCmple
	InstructionZjmp
	InstructionNzjmp
	InstructionJmp
	InstructionPrint
	InstructionHalt
)

func (i InstructionSet) String() string {
	switch i {
	case InstructionNoOp:
		return "NOOP"
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
	case InstructionInDup:
		return "INDUP"
	case InstructionSwap:
		return "SWAP"
	case InstructionInSwap:
		return "INSWAP"
	case InstructionCmpe:
		return "CMPE"
	case InstructionCmpne:
		return "CMPNE"
	case InstructionCmpg:
		return "CMPG"
	case InstructionCmpl:
		return "CMPL"
	case InstructionCmpge:
		return "CMPGE"
	case InstructionCmple:
		return "CMPLE"
	case InstructionJmp:
		return "JMP"
	case InstructionZjmp:
		return "ZJMP"
	case InstructionNzjmp:
		return "NZJMP"
	case InstructionMod:
		return "MOD"
	case InstructionHalt:
		return "HALT"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", i)
	}
}

func runInstructions(machine *Machine) *Machine {
	for insPtr := 0; insPtr < len(machine.instructions); insPtr++ {
		instr := machine.instructions[insPtr]
		switch instr.instructionType {
		case InstructionNoOp:
			// do nothing
		case InstructionPush:
			push(machine, instr.value)
		case InstructionPop:
			pop(machine)
		case InstructionDup:
			x := pop(machine)
			push(machine, x)
			push(machine, x)
		case InstructionInDup:
			indexDup(machine, instr.value)
		case InstructionSwap:
			a := pop(machine)
			b := pop(machine)
			push(machine, a)
			push(machine, b)
		case InstructionInSwap:
			indexSwap(machine, instr.value)
		case InstructionMod:
			a := pop(machine)
			b := pop(machine)
			if b == 0 {
				panic("ERROR: modulo by zero")
			}
			push(machine, a%b)
		case InstructionCmpe:
			a := pop(machine)
			b := pop(machine)
			if a == b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
		case InstructionCmpne:
			a := pop(machine)
			b := pop(machine)
			if a != b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
		case InstructionCmpg:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a > b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
		case InstructionCmpl:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a < b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
		case InstructionCmpge:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a >= b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
		case InstructionCmple:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a <= b {
				push(machine, 1)
			} else {
				push(machine, 0)
			}
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
			if b == 0 {
				panic("ERROR: division by zero")
			}
			push(machine, a/b)
		case InstructionJmp:
			target := int(instr.value)
			if target >= machine.programSize() {
				panic("ERROR: jump target out of bounds")
			}
			insPtr = target - 1 // -1 because loop will increment
		case InstructionNzjmp:
			value := pop(machine)
			if value != 0 {
				target := int(instr.value)
				if target >= machine.programSize() {
					panic("ERROR: jump target out of bounds")
				}
				insPtr = target - 1 // -1 because loop will increment
			}
		case InstructionZjmp:
			value := pop(machine)
			if value == 0 {
				target := int(instr.value)
				if target >= machine.programSize() {
					panic("ERROR: jump target out of bounds")
				}
				insPtr = target - 1 // -1 because loop will increment
			}
		case InstructionPrint:
			value := pop(machine)
			fmt.Println(value)
		case InstructionHalt:
			insPtr = machine.programSize()
		default:
			panic("ERROR: unknown instruction")
		}
	}
	return machine
}
