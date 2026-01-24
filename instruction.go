package main

import (
	"fmt"
	"os"
)

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
	InstructionNative
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
	case InstructionNative:
		return "NATIVE"
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
			if instr.value.Type() != LiteralInt {
				panic("ERROR: indup requires integer arguments")
			}
			indexDup(machine, instr.value.valueInt)
		case InstructionSwap:
			a := pop(machine)
			b := pop(machine)
			push(machine, a)
			push(machine, b)
		case InstructionInSwap:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: inswap requires integer arguments")
			}
			indexSwap(machine, instr.value.valueInt)
		case InstructionMod:
			a := pop(machine)
			b := pop(machine)
			push(machine, b.Mod(a))
		case InstructionCmpe:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a.Equal(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionCmpne:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if !a.Equal(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionCmpg:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a.Greater(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionCmpl:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a.Less(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionCmpge:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a.GreaterOrEqual(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionCmple:
			a := pop(machine)
			b := pop(machine)
			push(machine, b)
			push(machine, a)
			if a.LessOrEqual(b) {
				push(machine, IntLiteral(1))
			} else {
				push(machine, IntLiteral(0))
			}
		case InstructionAdd:
			a := pop(machine)
			b := pop(machine)
			push(machine, a.Add(b))
		case InstructionSub:
			a := pop(machine)
			b := pop(machine)
			push(machine, b.Sub(a))
		case InstructionMul:
			a := pop(machine)
			b := pop(machine)
			push(machine, a.Mul(b))
		case InstructionDiv:
			a := pop(machine)
			b := pop(machine)
			push(machine, a.Div(b))
		case InstructionJmp:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: jump target must be an integer")
			}
			target := int(instr.value.valueInt)
			if target >= machine.programSize() || target < 0 {
				panic("ERROR: jump target out of bounds")
			}
			insPtr = target - 1 // -1 because loop will increment
		case InstructionNzjmp:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: jump target must be an integer")
			}
			value := pop(machine)
			if value.Type() != LiteralInt {
				panic("ERROR: nzjmp condition value must be an integer")
			}
			if value.valueInt != 0 {
				target := int(instr.value.valueInt)
				if target >= machine.programSize() || target < 0 {
					panic("ERROR: jump target out of bounds")
				}
				insPtr = target - 1 // -1 because loop will increment
			}
		case InstructionZjmp:
			value := pop(machine)
			if instr.value.Type() != LiteralInt {
				panic("ERROR: jump target must be an integer")
			}
			if value.Type() != LiteralInt {
				panic("ERROR: zjmp condition value must be an integer")
			}
			if value.valueInt == 0 {
				target := int(instr.value.valueInt)
				if target >= machine.programSize() || target < 0 {
					panic("ERROR: jump target out of bounds")
				}
				insPtr = target - 1 // -1 because loop will increment
			}
		case InstructionPrint:
			value := pop(machine)
			fmt.Println(value)
		case InstructionNative:
			syscallID := instr.value
			if syscallID.Type() != LiteralInt {
				panic(instr.Error("native syscall ID must be integer"))
			}

			switch syscallID.valueInt {
			case 1:
				nativeWrite(machine, instr)
			default:
				panic(instr.Error(fmt.Sprintf("unknown native syscall ID: %d", syscallID.valueInt)))
			}
		case InstructionHalt:
			insPtr = machine.programSize()
		default:
			panic("ERROR: unknown instruction")
		}
	}
	return machine
}

func nativeWrite(machine *Machine, instr Instruction) {
	// Arguments are on stack: [..., fd, length] (Top is length)
	// Pop length
	lenVal := pop(machine)
	if lenVal.Type() != LiteralInt {
		panic(instr.Error("write length must be integer"))
	}
	length := int(lenVal.valueInt)

	// Pop fd
	fdVal := pop(machine)
	if fdVal.Type() != LiteralInt {
		panic(instr.Error("write fd must be integer"))
	}
	fd := int(fdVal.valueInt)

	if fd != 1 && fd != 2 {
		panic(instr.Error("unknown file descriptor"))
	}

	// Pop string characters
	s := ""
	for i := 0; i < length; i++ {
		val := pop(machine)
		if val.Type() != LiteralChar {
			panic(instr.Error("write expects characters on stack"))
		}
		s = string(val.valueChar) + s
	}

	if fd == 1 {
		fmt.Fprint(os.Stdout, s)
	} else {
		fmt.Fprint(os.Stderr, s)
	}
}
