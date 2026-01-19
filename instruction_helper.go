package main

import (
	"fmt"
	"strconv"

	"vm/internal/parser"
	"vm/internal/token"
)

type InstructionList []Instruction

// ---- Stack helper functions ----

func push(machine *Machine, value Literal) {
	if len(machine.stack) >= maxStackSize {
		panic("ERROR: stack overflow")
	}
	if value.Type() == LiteralInt {
		machine.stack = append(machine.stack, value)
	} else if value.Type() == LiteralFloat {
		machine.stack = append(machine.stack, value)
	} else if value.Type() == LiteralChar {
		machine.stack = append(machine.stack, value)
	}
}

func pop(machine *Machine) Literal {
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

func pushIntIns(value int64, line int) Instruction {
	return Instruction{instructionType: InstructionPush, value: IntLiteral(value), line: line}
}

func pushFloatIns(value float64, line int) Instruction {
	return Instruction{instructionType: InstructionPush, value: FloatLiteral(value), line: line}
}

func pushCharIns(value rune, line int) Instruction {
	return Instruction{instructionType: InstructionPush, value: CharLiteral(value), line: line}
}

func popIns(line int) Instruction {
	return Instruction{instructionType: InstructionPop, line: line}
}

func dupIns(line int) Instruction {
	return Instruction{instructionType: InstructionDup, line: line}
}

func inDupIns(index int64, line int) Instruction {
	return Instruction{instructionType: InstructionInDup, value: IntLiteral(index), line: line}
}

func swapIns(line int) Instruction {
	return Instruction{instructionType: InstructionSwap, line: line}
}

func inSwapIns(index int64, line int) Instruction {
	return Instruction{instructionType: InstructionInSwap, value: IntLiteral(index), line: line}
}

func addIns(line int) Instruction {
	return Instruction{instructionType: InstructionAdd, line: line}
}

func subIns(line int) Instruction {
	return Instruction{instructionType: InstructionSub, line: line}
}

func mulIns(line int) Instruction {
	return Instruction{instructionType: InstructionMul, line: line}
}

func divIns(line int) Instruction {
	return Instruction{instructionType: InstructionDiv, line: line}
}

func printIns(line int) Instruction {
	return Instruction{instructionType: InstructionPrint, line: line}
}

func cmpeIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmpe, line: line}
}

func cmpneIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmpne, line: line}
}

func cmpgIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmpg, line: line}
}

func cmplIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmpl, line: line}
}

func cmpgeIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmpge, line: line}
}

func cmpleIns(line int) Instruction {
	return Instruction{instructionType: InstructionCmple, line: line}
}

func modIns(line int) Instruction {
	return Instruction{instructionType: InstructionMod, line: line}
}

func jmpIns(target int64, line int) Instruction {
	return Instruction{instructionType: InstructionJmp, value: IntLiteral(target), line: line}
}

func zjmpIns(target int64, line int) Instruction {
	return Instruction{instructionType: InstructionZjmp, value: IntLiteral(target), line: line}
}

func nzjmpIns(target int64, line int) Instruction {
	return Instruction{instructionType: InstructionNzjmp, value: IntLiteral(target), line: line}
}

func haltIns(line int) Instruction {
	return Instruction{instructionType: InstructionHalt, line: line}
}

func noopIns(line int) Instruction {
	return Instruction{instructionType: InstructionNoOp, line: line}
}

func printStack(machine *Machine) {
	fmt.Println("------ STACK")
	for i := 0; i < len(machine.stack); i++ {
		fmt.Printf("[%d]: %s\n", i, machine.stack[i].String())
	}
	fmt.Println("------ END OF STACK")
}

func generateInstructions(parsedTokens *parser.ParserList) InstructionList {
	instructions := []Instruction{}
	cur := parsedTokens
	for cur != nil {
		line := int(cur.Value.Line)
		switch cur.Value.Type {
		case token.TypeInvalid:
			panic(fmt.Sprintf("ERROR on line %d: invalid token encountered during instruction generation", line))
		case token.TypeNoOp:
			instructions = append(instructions, noopIns(line))
		case token.TypePush:
			if cur.Next.Value.Type == token.TypeInt {
				value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
				if err != nil {
					panic(fmt.Sprintf("ERROR on line %d: invalid integer value for push instruction", line))
				}
				instructions = append(instructions, pushIntIns(value, line))
			} else if cur.Next.Value.Type == token.TypeFloat {
				value, err := strconv.ParseFloat(cur.Next.Value.Text, 64)
				if err != nil {
					panic(fmt.Sprintf("ERROR on line %d: invalid float value for push instruction", line))
				}
				instructions = append(instructions, pushFloatIns(value, line))
			} else if cur.Next.Value.Type == token.TypeChar {
				if len(cur.Next.Value.Text) == 0 {
					panic(fmt.Sprintf("ERROR on line %d: empty character literal", line))
				}
				charValue := rune(cur.Next.Value.Text[0])
				instructions = append(instructions, pushCharIns(charValue, line))
			}
			cur = cur.Next
		case token.TypePop:
			instructions = append(instructions, popIns(line))
		case token.TypeDup:
			instructions = append(instructions, dupIns(line))
		case token.TypeInDup:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("ERROR on line %d: invalid integer value for indup instruction", line))
			}
			cur = cur.Next
			instructions = append(instructions, inDupIns(value, line))
		case token.TypeSwap:
			instructions = append(instructions, swapIns(line))
		case token.TypeInSwap:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("ERROR on line %d: invalid integer value for inswap instruction", line))
			}
			cur = cur.Next
			instructions = append(instructions, inSwapIns(value, line))
		case token.TypeAdd:
			instructions = append(instructions, addIns(line))
		case token.TypeSub:
			instructions = append(instructions, subIns(line))
		case token.TypeMul:
			instructions = append(instructions, mulIns(line))
		case token.TypeDiv:
			instructions = append(instructions, divIns(line))
		case token.TypeCmpe:
			instructions = append(instructions, cmpeIns(line))
		case token.TypeCmpne:
			instructions = append(instructions, cmpneIns(line))
		case token.TypeCmpg:
			instructions = append(instructions, cmpgIns(line))
		case token.TypeCmpl:
			instructions = append(instructions, cmplIns(line))
		case token.TypeCmpge:
			instructions = append(instructions, cmpgeIns(line))
		case token.TypeCmple:
			instructions = append(instructions, cmpleIns(line))
		case token.TypeMod:
			instructions = append(instructions, modIns(line))
		case token.TypeJmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("ERROR on line %d: invalid integer value for jmp instruction", line))
			}
			cur = cur.Next
			instructions = append(instructions, jmpIns(value, line))
		case token.TypeZjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("ERROR on line %d: invalid integer value for zjmp instruction", line))
			}
			cur = cur.Next
			instructions = append(instructions, zjmpIns(value, line))
		case token.TypeNzjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(fmt.Sprintf("ERROR on line %d: invalid integer value for nzjmp instruction", line))
			}
			cur = cur.Next
			instructions = append(instructions, nzjmpIns(value, line))
		case token.TypePrint:
			instructions = append(instructions, printIns(line))
		case token.TypeInt:
			panic(fmt.Sprintf("ERROR on line %d: unexpected standalone integer token encountered during instruction generation", line))
		case token.TypeLabelDefinition:
			panic(fmt.Sprintf("ERROR on line %d: unexpected label definition token encountered during instruction generation", line))
		case token.TypeLabel:
			panic(fmt.Sprintf("ERROR on line %d: unexpected label token encountered during instruction generation", line))
		case token.TypeHalt:
			instructions = append(instructions, haltIns(line))
		default:
			panic(fmt.Sprintf("ERROR on line %d: unknown token type encountered during instruction generation", line))
		}
		cur = cur.Next
	}
	return instructions
}

func (il InstructionList) Print() {
	for i, instr := range il {
		fmt.Printf("[%d]: Type=%s", i, instr.instructionType.String())
		if instr.value.Type() == LiteralInt {
			fmt.Printf(", ValueInt=%d", instr.value.valueInt)
		}
		if instr.value.Type() == LiteralFloat {
			fmt.Printf(", ValueFloat=%f", instr.value.valueFloat)
		}
		if instr.value.Type() == LiteralChar {
			fmt.Printf(", ValueChar=%c", instr.value.valueChar)
		}
		fmt.Println()
	}
}
