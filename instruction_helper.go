package main

import (
	"fmt"
	"strconv"

	"vm/internal/parser"
	"vm/internal/token"
)

type InstructionList []Instruction

// InstructionContext contains metadata for instruction creation
type InstructionContext struct {
	Line      int
	FileName  string
	Character int // Reserved for future use
}

// Error formats an error message with file and line context
func (ctx InstructionContext) Error(message string) string {
	return fmt.Sprintf("ERROR (%s:%d): %s", ctx.FileName, ctx.Line, message)
}

// ---- Stack helper functions ----

func push(ctx *RuntimeContext, value Literal) {
	if len(ctx.stack) >= maxStackSize {
		panic(ctx.CurrentInstruction.Error("stack overflow"))
	}
	if value.Type() == LiteralInt {
		ctx.stack = append(ctx.stack, value)
	} else if value.Type() == LiteralFloat {
		ctx.stack = append(ctx.stack, value)
	} else if value.Type() == LiteralChar {
		ctx.stack = append(ctx.stack, value)
	} else if value.Type() == LiteralString {
		ctx.stack = append(ctx.stack, value)
	} else if value.Type() == LiteralNull {
		ctx.stack = append(ctx.stack, value)
	}
}

func pop(ctx *RuntimeContext) Literal {
	if len(ctx.stack) == 0 {
		panic(ctx.CurrentInstruction.Error("stack underflow"))
	}
	value := ctx.stack[len(ctx.stack)-1]
	ctx.stack = ctx.stack[:len(ctx.stack)-1]
	return value
}
func indexSwap(ctx *RuntimeContext, index int64) {
	if index < 0 || int(index) >= len(ctx.stack) {
		panic(ctx.CurrentInstruction.Error("index out of bounds for swap"))
	}
	targetIdx := int(index)
	topIdx := len(ctx.stack) - 1

	temp := ctx.stack[targetIdx]
	ctx.stack[targetIdx] = ctx.stack[topIdx]
	ctx.stack[topIdx] = temp
}

func indexDup(ctx *RuntimeContext, index int64) {
	if index < 0 || int(index) >= len(ctx.stack) {
		panic(ctx.CurrentInstruction.Error("index out of bounds for dup"))
	}
	targetIdx := int(index)
	push(ctx, ctx.stack[targetIdx])
}

func (machine *Machine) programSize() int {
	return len(machine.instructions)
}

func pushStr(ctx *RuntimeContext, val int64) {
	if len(ctx.strStack) >= maxStrStackSize {
		panic(ctx.CurrentInstruction.Error("string stack overflow"))
	}
	ctx.strStack = append(ctx.strStack, val)
}

func popStr(ctx *RuntimeContext) int64 {
	if len(ctx.strStack) == 0 {
		panic(ctx.CurrentInstruction.Error("string stack underflow"))
	}
	val := ctx.strStack[len(ctx.strStack)-1]
	ctx.strStack = ctx.strStack[:len(ctx.strStack)-1]
	return val
}

func indexDupStr(ctx *RuntimeContext, index int64) {
	if index < 0 || int(index) >= len(ctx.strStack) {
		panic(ctx.CurrentInstruction.Error("index out of bounds for indup_str"))
	}
	targetIdx := int(index)
	pushStr(ctx, ctx.strStack[targetIdx])
}

func indexSwapStr(ctx *RuntimeContext, index int64) {
	if index < 0 || int(index) >= len(ctx.strStack) {
		panic(ctx.CurrentInstruction.Error("index out of bounds for inswap_str"))
	}
	targetIdx := int(index)
	topIdx := len(ctx.strStack) - 1

	temp := ctx.strStack[targetIdx]
	ctx.strStack[targetIdx] = ctx.strStack[topIdx]
	ctx.strStack[topIdx] = temp
}

// ---- Instruction helper functions ----

func pushIntIns(value int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPush, value: IntLiteral(value), line: ctx.Line, fileName: ctx.FileName}
}

func pushPtrIns(value int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPushPtr, value: IntLiteral(value), line: ctx.Line, fileName: ctx.FileName}
}

func pushNullIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPushPtr, value: NullLiteral(), line: ctx.Line, fileName: ctx.FileName}
}

func pushFloatIns(value float64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPush, value: FloatLiteral(value), line: ctx.Line, fileName: ctx.FileName}
}

func pushCharIns(value rune, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPush, value: CharLiteral(value), line: ctx.Line, fileName: ctx.FileName}
}

func popIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPop, line: ctx.Line, fileName: ctx.FileName}
}

func dupIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionDup, line: ctx.Line, fileName: ctx.FileName}
}

func inDupIns(index int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionInDup, value: IntLiteral(index), line: ctx.Line, fileName: ctx.FileName}
}

func swapIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionSwap, line: ctx.Line, fileName: ctx.FileName}
}

func inSwapIns(index int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionInSwap, value: IntLiteral(index), line: ctx.Line, fileName: ctx.FileName}
}

func getStrIns(index int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionGetStr, value: IntLiteral(index), line: ctx.Line, fileName: ctx.FileName}
}

func addIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionAdd, line: ctx.Line, fileName: ctx.FileName}
}

func subIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionSub, line: ctx.Line, fileName: ctx.FileName}
}

func mulIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionMul, line: ctx.Line, fileName: ctx.FileName}
}

func divIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionDiv, line: ctx.Line, fileName: ctx.FileName}
}

func printIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPrint, line: ctx.Line, fileName: ctx.FileName}
}

func cmpeIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmpe, line: ctx.Line, fileName: ctx.FileName}
}

func cmpneIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmpne, line: ctx.Line, fileName: ctx.FileName}
}

func cmpgIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmpg, line: ctx.Line, fileName: ctx.FileName}
}

func cmplIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmpl, line: ctx.Line, fileName: ctx.FileName}
}

func cmpgeIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmpge, line: ctx.Line, fileName: ctx.FileName}
}

func cmpleIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCmple, line: ctx.Line, fileName: ctx.FileName}
}

func modIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionMod, line: ctx.Line, fileName: ctx.FileName}
}

func jmpIns(target int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionJmp, value: IntLiteral(target), line: ctx.Line, fileName: ctx.FileName}
}

func zjmpIns(target int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionZjmp, value: IntLiteral(target), line: ctx.Line, fileName: ctx.FileName}
}

func nzjmpIns(target int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionNzjmp, value: IntLiteral(target), line: ctx.Line, fileName: ctx.FileName}
}

func haltIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionHalt, line: ctx.Line, fileName: ctx.FileName}
}

func nativeIns(id int64, ctx InstructionContext) Instruction {
	return Instruction{
		instructionType: InstructionNative,
		value:           IntLiteral(id),
		line:            ctx.Line,
		fileName:        ctx.FileName,
	}
}

func callIns(label int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCall, value: IntLiteral(label), line: ctx.Line, fileName: ctx.FileName}
}

func retIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionRet, line: ctx.Line, fileName: ctx.FileName}
}

func popStrIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionPopStr, line: ctx.Line, fileName: ctx.FileName}
}

func dupStrIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionDupStr, line: ctx.Line, fileName: ctx.FileName}
}

func inDupStrIns(index int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionInDupStr, value: IntLiteral(index), line: ctx.Line, fileName: ctx.FileName}
}

func swapStrIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionSwapStr, line: ctx.Line, fileName: ctx.FileName}
}

func inSwapStrIns(index int64, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionInSwapStr, value: IntLiteral(index), line: ctx.Line, fileName: ctx.FileName}
}

func castIntToFloatIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCastIntToFloat, line: ctx.Line, fileName: ctx.FileName}
}

func castFloatToIntIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionCastFloatToInt, line: ctx.Line, fileName: ctx.FileName}
}

func refIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionRef, line: ctx.Line, fileName: ctx.FileName}
}

func derefIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionDeref, line: ctx.Line, fileName: ctx.FileName}
}

func movStrIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionMovStr, line: ctx.Line, fileName: ctx.FileName}
}

func indexIns(val rune, ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionIndex, value: CharLiteral(val), line: ctx.Line, fileName: ctx.FileName}
}

func noopIns(ctx InstructionContext) Instruction {
	return Instruction{instructionType: InstructionNoOp, line: ctx.Line, fileName: ctx.FileName}
}

func printStack(machine *Machine) {
	fmt.Println("------ STACK")
	for i := 0; i < len(machine.stack); i++ {
		fmt.Printf("[%d]: %s\n", i, machine.stack[i].String())
	}
	fmt.Println("------ END OF STACK")
}

func generateInstructions(parsedTokens *parser.ParserList) (InstructionList, int) {
	instructions := []Instruction{}
	cur := parsedTokens
	entrypointIndex := 0

	for cur != nil {
		ctx := InstructionContext{
			Line:     int(cur.Value.Line),
			FileName: cur.Value.FileName,
		}

		switch cur.Value.Type {
		case token.TypeInvalid:
			panic(ctx.Error("invalid token encountered during instruction generation"))
		case token.TypeNoOp:
			instructions = append(instructions, noopIns(ctx))
		case token.TypeCall:
			if cur.Next.Value.Type != token.TypeInt {
				panic(ctx.Error("expected integer (label address) after call"))
			}
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for call instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, callIns(value, ctx))
		case token.TypeRet:
			instructions = append(instructions, retIns(ctx))
		case token.TypeEntrypoint:
			if cur.Next.Value.Type != token.TypeInt {
				panic(ctx.Error("expected integer (label address) after entrypoint"))
			}
			val, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for entrypoint"))
			}
			entrypointIndex = int(val)
			cur = cur.Next
		case token.TypePush:
			if cur.Next.Value.Type == token.TypeInt {
				value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
				if err != nil {
					panic(ctx.Error("invalid integer value for push instruction"))
				}
				instructions = append(instructions, pushIntIns(value, ctx))
			} else if cur.Next.Value.Type == token.TypeFloat {
				value, err := strconv.ParseFloat(cur.Next.Value.Text, 64)
				if err != nil {
					panic(ctx.Error("invalid float value for push instruction"))
				}
				instructions = append(instructions, pushFloatIns(value, ctx))
			} else if cur.Next.Value.Type == token.TypeChar {
				if len(cur.Next.Value.Text) == 0 {
					panic(ctx.Error("empty character literal"))
				}
				charValue := rune(cur.Next.Value.Text[0])
				instructions = append(instructions, pushCharIns(charValue, ctx))
			} else if cur.Next.Value.Type == token.TypeString {
				strVal := cur.Next.Value.Text
				for _, char := range strVal {
					instructions = append(instructions, pushCharIns(char, ctx))
				}
			} else if cur.Next.Value.Type == token.TypeNull {
				instructions = append(instructions, pushNullIns(ctx))
			}
			cur = cur.Next
		case token.TypePushPtr:
			if cur.Next.Value.Type == token.TypeInt {
				value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
				if err != nil {
					panic(ctx.Error("invalid integer value for push_ptr instruction"))
				}
				instructions = append(instructions, pushPtrIns(value, ctx))
			} else if cur.Next.Value.Type == token.TypeNull {
				instructions = append(instructions, pushNullIns(ctx))
			}
			cur = cur.Next
		case token.TypePop:
			instructions = append(instructions, popIns(ctx))
		case token.TypeDup:
			instructions = append(instructions, dupIns(ctx))
		case token.TypeInDup:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for indup instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, inDupIns(value, ctx))
		case token.TypeSwap:
			instructions = append(instructions, swapIns(ctx))
		case token.TypeInSwap:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for inswap instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, inSwapIns(value, ctx))
		case token.TypeAdd:
			instructions = append(instructions, addIns(ctx))
		case token.TypeSub:
			instructions = append(instructions, subIns(ctx))
		case token.TypeMul:
			instructions = append(instructions, mulIns(ctx))
		case token.TypeDiv:
			instructions = append(instructions, divIns(ctx))
		case token.TypeCmpe:
			instructions = append(instructions, cmpeIns(ctx))
		case token.TypeCmpne:
			instructions = append(instructions, cmpneIns(ctx))
		case token.TypeCmpg:
			instructions = append(instructions, cmpgIns(ctx))
		case token.TypeCmpl:
			instructions = append(instructions, cmplIns(ctx))
		case token.TypeCmpge:
			instructions = append(instructions, cmpgeIns(ctx))
		case token.TypeCmple:
			instructions = append(instructions, cmpleIns(ctx))
		case token.TypeMod:
			instructions = append(instructions, modIns(ctx))
		case token.TypeJmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for jmp instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, jmpIns(value, ctx))
		case token.TypeZjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for zjmp instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, zjmpIns(value, ctx))
		case token.TypeNzjmp:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for nzjmp instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, nzjmpIns(value, ctx))
		case token.TypeNative:
			id, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for native function ID"))
			}
			cur = cur.Next
			instructions = append(instructions, nativeIns(id, ctx))
		case token.TypePrint:
			instructions = append(instructions, printIns(ctx))
		case token.TypePushStr:
			cur = cur.Next
		case token.TypeGetStr:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for get_str instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, getStrIns(value, ctx))
		case token.TypeInt:
			panic(ctx.Error("unexpected standalone integer token encountered during instruction generation"))
		case token.TypeLabelDefinition:
			panic(ctx.Error("unexpected label definition token encountered during instruction generation"))
		case token.TypeLabel:
			panic(ctx.Error("unexpected label token encountered during instruction generation"))
		case token.TypeNull:
			instructions = append(instructions, pushNullIns(ctx))
		case token.TypeHalt:
			instructions = append(instructions, haltIns(ctx))
		case token.TypePopStr:
			instructions = append(instructions, popStrIns(ctx))
		case token.TypeDupStr:
			instructions = append(instructions, dupStrIns(ctx))
		case token.TypeInDupStr:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for indup_str instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, inDupStrIns(value, ctx))
		case token.TypeSwapStr:
			instructions = append(instructions, swapStrIns(ctx))
		case token.TypeInSwapStr:
			value, err := strconv.ParseInt(cur.Next.Value.Text, 10, 64)
			if err != nil {
				panic(ctx.Error("invalid integer value for inswap_str instruction"))
			}
			cur = cur.Next
			instructions = append(instructions, inSwapStrIns(value, ctx))
		case token.TypeCastIntToFloat:
			instructions = append(instructions, castIntToFloatIns(ctx))
		case token.TypeCastFloatToInt:
			instructions = append(instructions, castFloatToIntIns(ctx))
		case token.TypeRef:
			instructions = append(instructions, refIns(ctx))
		case token.TypeDeref:
			instructions = append(instructions, derefIns(ctx))
		case token.TypeMovStr:
			instructions = append(instructions, movStrIns(ctx))
		case token.TypeIndex:
			if cur.Next == nil {
				panic(ctx.Error("expected char after index"))
			}
			if cur.Next.Value.Type == token.TypeChar {
				if len(cur.Next.Value.Text) == 0 {
					panic(ctx.Error("empty character literal for index"))
				}
				charValue := rune(cur.Next.Value.Text[0])
				instructions = append(instructions, indexIns(charValue, ctx))
			} else {
				panic(ctx.Error("expected char after index"))
			}
			cur = cur.Next
		default:
			panic(ctx.Error("unknown token type encountered during instruction generation"))
		}
		cur = cur.Next
	}
	return instructions, entrypointIndex
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
