package main

import (
	"fmt"
	"strconv"

	"vm/internal/parser"
	"vm/internal/token"
)

type InstructionList []Instruction

// ---- Stack helper functions ----

func push(machine *Machine, value int64) {
	if len(machine.stack) >= maxStackSize {
		panic("ERROR: stack overflow")
	}
	machine.stack = append(machine.stack, value)
}

func pop(machine *Machine) int64 {
	if len(machine.stack) == 0 {
		panic("ERROR: stack underflow")
	}
	value := machine.stack[len(machine.stack)-1]
	machine.stack = machine.stack[:len(machine.stack)-1]
	return value
}

func indexSwap(machine *Machine, index int64) {
	if index < 0 || int(index) >= len(machine.stack) {
		panic("ERROR: index out of bounds for swap")
	}
	tempValue := machine.stack[index]
	machine.stack[index] = pop(machine)
	push(machine, tempValue)
}

func indexDup(machine *Machine, index int64) {
	if index < 0 || int(index) >= len(machine.stack) {
		panic("ERROR: index out of bounds for swap")
	}
	push(machine, machine.stack[index])
}

func (machine *Machine) programSize() int {
	return len(machine.instructions)
}

// ---- Instruction helper functions ----

func pushIns(value int64) Instruction {
	return Instruction{instructionType: InstructionPush, value: value}
}

func popIns() Instruction {
	return Instruction{instructionType: InstructionPop}
}

func dupIns() Instruction {
	return Instruction{instructionType: InstructionDup}
}

func inDupIns(index int64) Instruction {
	return Instruction{instructionType: InstructionInDup, value: index}
}

func swapIns() Instruction {
	return Instruction{instructionType: InstructionSwap}
}

func inSwapIns(index int64) Instruction {
	return Instruction{instructionType: InstructionInSwap, value: index}
}

func addIns() Instruction {
	return Instruction{instructionType: InstructionAdd}
}

func subIns() Instruction {
	return Instruction{instructionType: InstructionSub}
}

func mulIns() Instruction {
	return Instruction{instructionType: InstructionMul}
}

func divIns() Instruction {
	return Instruction{instructionType: InstructionDiv}
}

func printIns() Instruction {
	return Instruction{instructionType: InstructionPrint}
}

func cmpeIns() Instruction {
	return Instruction{instructionType: InstructionCmpe}
}

func cmpneIns() Instruction {
	return Instruction{instructionType: InstructionCmpne}
}

func cmpgIns() Instruction {
	return Instruction{instructionType: InstructionCmpg}
}

func cmplIns() Instruction {
	return Instruction{instructionType: InstructionCmpl}
}

func cmpgeIns() Instruction {
	return Instruction{instructionType: InstructionCmpge}
}

func cmpleIns() Instruction {
	return Instruction{instructionType: InstructionCmple}
}

func modIns() Instruction {
	return Instruction{instructionType: InstructionMod}
}

func jmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionJmp, value: target}
}

func zjmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionZjmp, value: target}
}

func nzjmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionNzjmp, value: target}
}

func haltIns() Instruction {
	return Instruction{instructionType: InstructionHalt}
}

func noopIns() Instruction {
	return Instruction{instructionType: InstructionNoOp}
}

func printStack(machine *Machine) {
	fmt.Println("------ STACK")
	for i := 0; i < len(machine.stack); i++ {
		fmt.Printf("[%d]: %d\n", i, machine.stack[i])
	}
	fmt.Println("------ END OF STACK")
}

func generateInstructions(parsedTokens *parser.ParserList) InstructionList {
	instructions := []Instruction{}
	cur := parsedTokens
	for cur != nil {
		switch cur.Value.Type {
		case token.TypeInvalid:
			panic("ERROR: invalid token encountered during instruction generation")
		case token.TypeNoOp:
			instructions = append(instructions, noopIns())
		case token.TypePush:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for push instruction")
			}
			cur = cur.Next
			instructions = append(instructions, pushIns(value))
		case token.TypePop:
			instructions = append(instructions, popIns())
		case token.TypeDup:
			instructions = append(instructions, dupIns())
		case token.TypeInDup:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for indup instruction")
			}
			cur = cur.Next
			instructions = append(instructions, inDupIns(value))
		case token.TypeSwap:
			instructions = append(instructions, swapIns())
		case token.TypeInSwap:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for inswap instruction")
			}
			cur = cur.Next
			instructions = append(instructions, inSwapIns(value))
		case token.TypeAdd:
			instructions = append(instructions, addIns())
		case token.TypeSub:
			instructions = append(instructions, subIns())
		case token.TypeMul:
			instructions = append(instructions, mulIns())
		case token.TypeDiv:
			instructions = append(instructions, divIns())
		case token.TypeCmpe:
			instructions = append(instructions, cmpeIns())
		case token.TypeCmpne:
			instructions = append(instructions, cmpneIns())
		case token.TypeCmpg:
			instructions = append(instructions, cmpgIns())
		case token.TypeCmpl:
			instructions = append(instructions, cmplIns())
		case token.TypeCmpge:
			instructions = append(instructions, cmpgeIns())
		case token.TypeCmple:
			instructions = append(instructions, cmpleIns())
		case token.TypeMod:
			instructions = append(instructions, modIns())
		case token.TypeJmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for jmp instruction")
			}
			cur = cur.Next
			instructions = append(instructions, jmpIns(value))
		case token.TypeZjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for zjmp instruction")
			}
			cur = cur.Next
			instructions = append(instructions, zjmpIns(value))
		case token.TypeNzjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic("ERROR: invalid integer value for nzjmp instruction")
			}
			cur = cur.Next
			instructions = append(instructions, nzjmpIns(value))
		case token.TypePrint:
			instructions = append(instructions, printIns())
		case token.TypeInt:
			panic("ERROR: unexpected standalone integer token encountered during instruction generation")
		case token.TypeHalt:
			instructions = append(instructions, haltIns())
		default:
			panic("ERROR: unknown token type encountered during instruction generation")
		}
		cur = cur.Next
	}
	return instructions
}

func (il InstructionList) Print() {
	for i, instr := range il {
		fmt.Printf("[%d]: Type=%s, Value=%d\n", i, instr.instructionType.String(), instr.value)
	}
}
