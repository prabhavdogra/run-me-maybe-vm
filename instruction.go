package main

import (
	"fmt"
	"io"
	"os"
	"time"
	"vm/internal/parser"
	"vm/internal/token"
)

type InstructionSet uint8

const (
	InstructionNoOp InstructionSet = iota
	InstructionPush
	InstructionPushPtr
	InstructionGetStr
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
	InstructionCall
	InstructionRet
	InstructionPushStr
	InstructionPopStr
	InstructionDupStr
	InstructionInDupStr
	InstructionSwapStr
	InstructionInSwapStr
	InstructionCastIntToFloat
	InstructionCastFloatToInt
	InstructionRef
	InstructionDeref
	InstructionHalt
)

func populateStringTable(parsedTokens *parser.ParserList) ([]int64, []Literal) {
	strStack := []int64{}
	heap := []Literal{}

	cur := parsedTokens
	for cur != nil {
		if cur.Value.Type == token.TypePushStr {
			if cur.Next == nil || cur.Next.Value.Type != token.TypeString {
				panic("expected string after push_str") // Should be caught by parser
			}
			strVal := cur.Next.Value.Text

			ptr := int64(len(heap))

			for _, char := range strVal {
				heap = append(heap, CharLiteral(char))
			}
			heap = append(heap, CharLiteral(0))

			strStack = append(strStack, ptr)
		}
		cur = cur.Next
	}
	return strStack, heap
}

func (i InstructionSet) String() string {
	switch i {
	case InstructionNoOp:
		return "NOOP"
	case InstructionPush:
		return "PUSH"
	case InstructionPushPtr:
		return "PUSH_PTR"
	case InstructionGetStr:
		return "GET_STR"
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
	case InstructionCall:
		return "CALL"
	case InstructionRet:
		return "RET"
	case InstructionPushStr:
		return "PUSH_STR"
	case InstructionPopStr:
		return "POP_STR"
	case InstructionDupStr:
		return "DUP_STR"
	case InstructionInDupStr:
		return "INDUP_STR"
	case InstructionSwapStr:
		return "SWAP_STR"
	case InstructionInSwapStr:
		return "INSWAP_STR"
	case InstructionCastIntToFloat:
		return "ITOF"
	case InstructionCastFloatToInt:
		return "FTOI"
	case InstructionRef:
		return "REF"
	case InstructionDeref:
		return "DEREF"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", i)
	}
}

func runInstructions(machine *Machine) *Machine {
	ctx := &RuntimeContext{
		Machine:     machine,
		returnStack: make([]int, 0, maxReturnStackSize),
	}
	// Jump to entrypoint
	insPtr := machine.entrypoint

	for insPtr < len(machine.instructions) {
		instr := machine.instructions[insPtr]
		ctx.CurrentInstruction = instr
		if debugMode {
			fmt.Fprintf(os.Stderr, "Line %d: %v, Stack: %+v\n", instr.line, instr.instructionType, ctx.stack)
		}

		jumped := false

		switch instr.instructionType {
		case InstructionNoOp:
			// do nothing
		case InstructionCall:
			if instr.value.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("call target must be an integer"))
			}
			target := int(instr.value.valueInt)
			if target >= machine.programSize() || target < 0 {
				panic(ctx.CurrentInstruction.Error("call target out of bounds"))
			}
			if len(ctx.returnStack) >= maxReturnStackSize {
				panic(ctx.CurrentInstruction.Error("return stack overflow"))
			}
			ctx.returnStack = append(ctx.returnStack, insPtr+1)
			insPtr = target
			jumped = true
		case InstructionRet:
			if len(ctx.returnStack) == 0 {
				panic(ctx.CurrentInstruction.Error("return stack underflow"))
			}
			retAddr := ctx.returnStack[len(ctx.returnStack)-1]
			ctx.returnStack = ctx.returnStack[:len(ctx.returnStack)-1]
			insPtr = retAddr
			jumped = true
		case InstructionPopStr:
			popStr(ctx)
		case InstructionDupStr:
			if len(ctx.strStack) == 0 {
				panic(ctx.CurrentInstruction.Error("string stack underflow"))
			}
			val := ctx.strStack[len(ctx.strStack)-1]
			pushStr(ctx, val)
		case InstructionInDupStr:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: indup_str requires integer arguments")
			}
			indexDupStr(ctx, instr.value.valueInt)
		case InstructionSwapStr:
			if len(ctx.strStack) < 2 {
				panic(ctx.CurrentInstruction.Error("string stack underflow"))
			}
			a := popStr(ctx)
			b := popStr(ctx)
			pushStr(ctx, a)
			pushStr(ctx, b)
		case InstructionInSwapStr:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: inswap_str requires integer arguments")
			}
			indexSwapStr(ctx, instr.value.valueInt)
		case InstructionCastIntToFloat:
			val := pop(ctx)
			if val.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("itof requires an integer"))
			}
			push(ctx, FloatLiteral(float64(val.valueInt)))
		case InstructionCastFloatToInt:
			val := pop(ctx)
			if val.Type() != LiteralFloat {
				panic(ctx.CurrentInstruction.Error("ftoi requires a float"))
			}
			push(ctx, IntLiteral(int64(val.valueFloat)))
		case InstructionRef:
			val := pop(ctx)
			ptr := int64(len(ctx.heap))
			ctx.heap = append(ctx.heap, val)
			push(ctx, IntLiteral(ptr))
		case InstructionDeref:
			ptrVal := pop(ctx)
			if ptrVal.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("deref requires a pointer (int)"))
			}
			ptr := ptrVal.valueInt
			if ptr < 0 || int(ptr) >= len(ctx.heap) {
				panic(ctx.CurrentInstruction.Error("segmentation fault: invalid pointer"))
			}
			val := ctx.heap[ptr]
			push(ctx, val)
		case InstructionPush:
			push(ctx, instr.value)
		case InstructionPushStr:
			if instr.value.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("push_str value must be integer pointer"))
			}
			pushStr(ctx, instr.value.valueInt)
		case InstructionPushPtr:
			if instr.value.Type() != LiteralInt && instr.value.Type() != LiteralNull {
				panic(ctx.CurrentInstruction.Error("push_ptr requires an integer or NULL value"))
			}
			push(ctx, instr.value)
		case InstructionGetStr:
			idx := int(instr.value.valueInt)
			if idx < 0 || idx >= len(machine.strStack) {
				panic(ctx.CurrentInstruction.Error("string index out of bounds"))
			}
			ptr := machine.strStack[idx]
			push(ctx, IntLiteral(ptr))
		case InstructionPop:
			pop(ctx)
		case InstructionDup:
			x := pop(ctx)
			push(ctx, x)
			push(ctx, x)
		case InstructionInDup:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: indup requires integer arguments")
			}
			indexDup(ctx, instr.value.valueInt)
		case InstructionSwap:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, a)
			push(ctx, b)
		case InstructionInSwap:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: inswap requires integer arguments")
			}
			indexSwap(ctx, instr.value.valueInt)
		case InstructionMod:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Mod(a))
		case InstructionCmpe:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if b.Equal(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpne:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if !b.Equal(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpg:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			if b.Greater(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpl:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			if b.Less(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpge:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if b.GreaterOrEqual(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmple:
			if len(ctx.stack) < 2 {
				panic(ctx.CurrentInstruction.Error("stack underflow"))
			}
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if b.LessOrEqual(a) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionAdd:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Add(a))
		case InstructionSub:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Sub(a))
		case InstructionMul:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Mul(a))
		case InstructionDiv:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Div(a))
		case InstructionJmp:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: jump target must be an integer")
			}
			target := int(instr.value.valueInt)
			if target >= machine.programSize() || target < 0 {
				panic("ERROR: jump target out of bounds")
			}
			insPtr = target
			jumped = true
		case InstructionNzjmp:
			if instr.value.Type() != LiteralInt {
				panic("ERROR: jump target must be an integer")
			}
			value := pop(ctx)
			if value.Type() != LiteralInt {
				panic("ERROR: nzjmp condition value must be an integer")
			}
			if value.valueInt != 0 {
				target := int(instr.value.valueInt)
				if target >= machine.programSize() || target < 0 {
					panic("ERROR: jump target out of bounds")
				}
				insPtr = target
				jumped = true
			}
		case InstructionZjmp:
			value := pop(ctx)
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
				insPtr = target
				jumped = true
			}
		case InstructionPrint:
			value := pop(ctx)
			fmt.Println(value)
		case InstructionNative:
			syscallID := instr.value
			if syscallID.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("native syscall ID must be integer"))
			}

			switch syscallID.valueInt {
			case 0:
				// open(flags, len, ptr)
				nativeOpen(ctx)
			case 1:
				// write(len, fd, char...)
				nativeWrite(ctx)
			case 2:
				// read(ptr, len, fd)
				nativeRead(ctx)
			case 3:
				// close(fd)
				nativeClose(ctx)
			case 4:
				// malloc(size)
				nativeMalloc(ctx)
			case 5:
				// realloc(ptr, size)
				nativeRealloc(ctx)
			case 6:
				// free(ptr)
				nativeFree(ctx)
			case 10:
				// time
				nativeTime(ctx)
			case 60:
				// exit(code)
				nativeExit(ctx)
			case 90:
				// 90: strcmp
				nativeStrcmp(ctx)
			case 91:
				// 91: strcpy
				nativeStrcpy(ctx)
			case 92:
				// 92: memcpy
				nativeMemcpy(ctx)
			case 93:
				// 93: strcat
				nativeStrcat(ctx)
			case 94:
				// 94: strlen
				nativeStrlen(ctx)
			case 99:
				// 99: int_to_str
				nativeIntToStr(ctx)
			case 100:
				// 100: assert
				nativeAssert(ctx)
			default:
				panic(ctx.CurrentInstruction.Error(fmt.Sprintf("unknown native syscall ID: %d", syscallID.valueInt)))
			}
		case InstructionHalt:
			insPtr = machine.programSize()
			jumped = true
		default:
			panic(ctx.CurrentInstruction.Error(fmt.Sprintf("unknown instruction type: %d", instr.instructionType)))
		}
		if !jumped {
			insPtr++
		}
	}
	return machine
}

// Native function ID 99: int_to_str
// Stack inputs: [int]
// Stack output: [ptr] (pointer to new string)
func nativeIntToStr(ctx *RuntimeContext) {
	value := pop(ctx)
	if value.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("int_to_str expects an integer"))
	}
	s := fmt.Sprintf("%d", value.valueInt)
	ptr := len(ctx.heap)
	for _, char := range s {
		ctx.heap = append(ctx.heap, CharLiteral(char))
	}
	ctx.heap = append(ctx.heap, CharLiteral(0))
	push(ctx, IntLiteral(int64(ptr)))
}

// Open a file
func nativeOpen(ctx *RuntimeContext) {
	// Pop flags
	flagsVal := pop(ctx)
	if flagsVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("open flags must be integer"))
	}
	flags := int(flagsVal.valueInt)

	// Translate VM flags to OS flags
	// VM: RONLY=0, WONLY=1, RDWR=2, CREAT=64, EXCL=128
	osFlags := 0
	if flags&0x3 == 0 {
		osFlags |= os.O_RDONLY
	} else if flags&0x3 == 1 {
		osFlags |= os.O_WRONLY
	} else if flags&0x3 == 2 {
		osFlags |= os.O_RDWR
	}

	if flags&64 != 0 {
		osFlags |= os.O_CREATE
	}
	if flags&128 != 0 {
		osFlags |= os.O_EXCL
	}

	// Pop filename length
	lenVal := pop(ctx)
	if lenVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("open filename length must be integer"))
	}
	length := int(lenVal.valueInt)

	// Pop filename pointer
	ptrVal := pop(ctx)
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("open filename pointer must be integer"))
	}

	// Read filename from heap
	ptr := int(ptrVal.valueInt)
	if ptr < 0 || ptr+length > len(ctx.heap) {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer for filename"))
	}

	filename := ""
	for i := 0; i < length; i++ {
		charLit := ctx.heap[ptr+i]
		if charLit.Type() != LiteralChar {
			panic(ctx.CurrentInstruction.Error("filename must be a string of characters"))
		}
		filename += string(charLit.valueChar)
	}

	// Open the file
	file, err := os.OpenFile(filename, osFlags, 0644)
	if err != nil {
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("failed to open file %s: %v", filename, err)))
	}

	// 0: Standard Input (stdin)
	// 1: Standard Output (stdout)
	// 2: Standard Error (stderr)

	// Find lowest available file descriptor
	fd := int64(3)
	for {
		if _, ok := ctx.fileDescriptors[fd]; !ok {
			break
		}
		fd++
	}

	ctx.fileDescriptors[fd] = file
	push(ctx, IntLiteral(fd))
}

// Write to a file descriptor
func nativeWrite(ctx *RuntimeContext) {
	fd := pop(ctx) // FD
	if fd.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("write fd must be integer"))
	}

	ptr := pop(ctx) // Ptr
	if ptr.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("write string pointer must be integer"))
	}

	var writer io.Writer
	if fd.valueInt == 1 {
		writer = ctx.output
	} else if fd.valueInt == 2 {
		writer = os.Stderr
	} else {
		if file, ok := ctx.fileDescriptors[int64(fd.valueInt)]; ok {
			writer = file
		} else {
			panic(ctx.CurrentInstruction.Error(fmt.Sprintf("unknown file descriptor %d", fd.valueInt)))
		}
	}

	ptrIdx := int(ptr.valueInt)
	if ptrIdx < 0 || ptrIdx >= len(ctx.heap) {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer"))
	}

	s := ""
	for i := ptrIdx; i < len(ctx.heap); i++ {
		charLit := ctx.heap[i]
		if charLit.Type() == LiteralInt {
			s += string(rune(charLit.valueInt))
			continue
		}
		if charLit.Type() != LiteralChar {
			continue
		}
		if charLit.valueChar == 0 {
			break
		}
		s += string(charLit.valueChar)
	}
	fmt.Fprint(writer, s)
	push(ctx, IntLiteral(int64(len(s))))
}

// Read from a file descriptor into a buffer
func nativeRead(ctx *RuntimeContext) {
	// Arguments: [..., fd, len, ptr] (Top is ptr)
	ptrVal := pop(ctx)
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("read buffer pointer must be integer"))
	}

	lenVal := pop(ctx)
	if lenVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("read length must be integer"))
	}
	length := int(lenVal.valueInt)

	fdVal := pop(ctx)
	if fdVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("read fd must be integer"))
	}
	fd := int64(fdVal.valueInt)

	var reader io.Reader
	if fd == 0 {
		reader = ctx.input
	} else {
		if file, ok := ctx.fileDescriptors[fd]; ok {
			reader = file
		} else {
			panic(ctx.CurrentInstruction.Error(fmt.Sprintf("read error: invalid file descriptor %d", fd)))
		}
	}

	// Validate Heap Pointer and Size
	ptr := int(ptrVal.valueInt)

	// Safety check against allocation size if tracked
	if allocSize, ok := ctx.allocations[ptr]; ok {
		if length > allocSize {
			panic(ctx.CurrentInstruction.Error("buffer overflow: read length exceeds allocated size"))
		}
	} else {
		// Fallback strictly to heap bounds
		if ptr < 0 || ptr+length > len(ctx.heap) {
			panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer or length"))
		}
	}

	// Read from Input
	buf := make([]byte, length)
	_, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("read error: %v", err)))
	}

	// Store in Heap
	for i, b := range buf {
		ctx.heap[ptr+i] = CharLiteral(rune(b))
	}
}

// Close a file descriptor
func nativeClose(ctx *RuntimeContext) {
	// Pop file descriptor ID
	fdVal := pop(ctx)
	if fdVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("close file descriptor must be integer"))
	}
	fd := int64(fdVal.valueInt)

	// Check if it's a valid custom file descriptor
	if fd < 3 { // 0, 1, 2 are stdin, stdout, stderr - cannot close
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("cannot close standard file descriptor %d", fd)))
	}

	file, ok := ctx.fileDescriptors[fd]
	if !ok {
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("invalid file descriptor %d", fd)))
	}

	err := file.Close()
	if err != nil {
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("failed to close file %d: %v", fd, err)))
	}

	delete(ctx.fileDescriptors, fd)
}

// Free a heap pointer
func nativeFree(ctx *RuntimeContext) {
	// Pop ptr
	ptrVal := pop(ctx)
	if ptrVal.Type() == LiteralNull {
		return
	}
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("free pointer must be integer"))
	}
	ptr := int(ptrVal.valueInt)

	// Check allocations map
	if _, ok := ctx.allocations[ptr]; !ok {
		panic(ctx.CurrentInstruction.Error("double free or invalid heap pointer"))
	}
	delete(ctx.allocations, ptr)
	// We don't shrink the heap slice, just mark as freed in allocations map
}

func nativeMalloc(ctx *RuntimeContext) {
	// Pop size
	sizeVal := pop(ctx)
	if sizeVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("malloc size must be integer"))
	}
	size := int(sizeVal.valueInt)

	// Allocate
	ptr := len(ctx.heap)
	for i := 0; i < size; i++ {
		ctx.heap = append(ctx.heap, CharLiteral(0))
	}

	// Track allocation
	ctx.allocations[ptr] = size

	// Push Pointer
	push(ctx, IntLiteral(int64(ptr)))
}

func nativeExit(ctx *RuntimeContext) {
	// Pop exit code
	codeVal := pop(ctx)
	if codeVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("exit code must be integer"))
	}
	code := int(codeVal.valueInt)
	os.Exit(code)
}

func nativeStrcmp(ctx *RuntimeContext) {
	ptr2Val := pop(ctx)
	ptr1Val := pop(ctx)

	if ptr1Val.Type() != LiteralInt || ptr2Val.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("strcmp pointers must be integer"))
	}

	ptr1 := ptr1Val.valueInt
	ptr2 := ptr2Val.valueInt

	s1 := getStringFromHeap(ctx, ptr1)
	s2 := getStringFromHeap(ctx, ptr2)

	if s1 == s2 {
		push(ctx, IntLiteral(1))
	} else {
		push(ctx, IntLiteral(0))
	}
}

// Helper to get string from heap
func getStringFromHeap(ctx *RuntimeContext, ptr int64) string {
	if int(ptr) < 0 || int(ptr) >= len(ctx.heap) {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer"))
	}
	s := ""
	for i := int(ptr); i < len(ctx.heap); i++ {
		charLit := ctx.heap[i]
		if charLit.Type() != LiteralChar {
			continue
		}
		if charLit.valueChar == 0 {
			break
		}
		s += string(charLit.valueChar)
	}
	return s
}

func nativeStrcpy(ctx *RuntimeContext) {
	srcPtrVal := pop(ctx)
	destPtrVal := pop(ctx)

	if srcPtrVal.Type() != LiteralInt || destPtrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("strcpy pointers must be integer"))
	}

	srcPtr := int(srcPtrVal.valueInt)
	destPtr := int(destPtrVal.valueInt)

	if srcPtr < 0 || srcPtr >= len(ctx.heap) {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid source pointer"))
	}
	if destPtr < 0 || destPtr >= len(ctx.heap) {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid destination pointer"))
	}

	for i := 0; srcPtr+i < len(ctx.heap); i++ {
		charLit := ctx.heap[srcPtr+i]
		if destPtr+i >= len(ctx.heap) {
			ctx.heap = append(ctx.heap, CharLiteral(0))
		}
		ctx.heap[destPtr+i] = charLit
		if charLit.Type() == LiteralChar && charLit.valueChar == 0 {
			break
		}
	}
	push(ctx, destPtrVal)
}

// Native function ID 92: memcpy
// Stack inputs: [dest, src, size] (top is size)
// Stack output: [dest]
func nativeMemcpy(ctx *RuntimeContext) {
	sizeVal := pop(ctx)
	srcPtrVal := pop(ctx)
	destPtrVal := pop(ctx)

	if sizeVal.Type() != LiteralInt || srcPtrVal.Type() != LiteralInt || destPtrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("memcpy arguments must be integer"))
	}

	size := int(sizeVal.valueInt)
	srcPtr := int(srcPtrVal.valueInt)
	destPtr := int(destPtrVal.valueInt)

	if srcPtr < 0 || srcPtr+size > len(ctx.heap) { // Strict bound check for src
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid source range"))
	}

	// Extend dest if needed
	if destPtr+size > len(ctx.heap) {
		required := (destPtr + size) - len(ctx.heap)
		for k := 0; k < required; k++ {
			ctx.heap = append(ctx.heap, CharLiteral(0))
		}
	}

	// Copy
	for i := 0; i < size; i++ {
		ctx.heap[destPtr+i] = ctx.heap[srcPtr+i]
	}

	push(ctx, destPtrVal)
}

func nativeRealloc(ctx *RuntimeContext) {
	sizeVal := pop(ctx)
	ptrVal := pop(ctx)

	if sizeVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("realloc size must be integer"))
	}
	if ptrVal.Type() != LiteralInt && ptrVal.Type() != LiteralNull {
		panic(ctx.CurrentInstruction.Error("realloc pointer must be integer or NULL"))
	}

	size := int(sizeVal.valueInt)

	// Check for NULL (Strictly LiteralNull)
	isNull := ptrVal.Type() == LiteralNull

	if isNull { // NULL check -> behaves like malloc
		// Allocate new
		newPtr := len(ctx.heap)
		for i := 0; i < size; i++ {
			ctx.heap = append(ctx.heap, CharLiteral(0))
		}
		ctx.allocations[newPtr] = size
		push(ctx, IntLiteral(int64(newPtr)))
		return
	}

	ptr := int(ptrVal.valueInt)

	// Check existing allocation
	oldSize, ok := ctx.allocations[ptr]
	if !ok {
		panic(ctx.CurrentInstruction.Error("realloc: invalid heap pointer"))
	}

	if size <= oldSize {
		// Shrinking or same size: Just update allocation size (oversimplified)
		// Usually we'd want to reuse usage.
		ctx.allocations[ptr] = size
		push(ctx, ptrVal)
		return
	}

	// Expand: Allocate new, copy, free old (simple implementation)
	newPtr := len(ctx.heap)
	for i := 0; i < size; i++ {
		ctx.heap = append(ctx.heap, CharLiteral(0))
	}
	ctx.allocations[newPtr] = size

	// Copy data
	for i := 0; i < oldSize; i++ {
		ctx.heap[newPtr+i] = ctx.heap[ptr+i]
	}

	// 'Free' old (remove from allocations)
	delete(ctx.allocations, ptr)

	push(ctx, IntLiteral(int64(newPtr)))
}

func nativeTime(ctx *RuntimeContext) {
	now := time.Now().Unix()
	push(ctx, IntLiteral(now))
}

func nativeStrcat(ctx *RuntimeContext) {
	srcPtrVal := pop(ctx)
	destPtrVal := pop(ctx)

	if srcPtrVal.Type() != LiteralInt || destPtrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("strcat pointers must be integer"))
	}

	srcPtr := int(srcPtrVal.valueInt)
	destPtr := int(destPtrVal.valueInt)

	// Get length of dest
	sDest := getStringFromHeap(ctx, int64(destPtr))
	destLen := len(sDest)

	// Append pointer
	appendPtr := destPtr + destLen

	// Copy src to dest end
	// Note: We need to ensure heap is large enough or extend it.
	// getStringFromHeap reads until null.

	sSrc := getStringFromHeap(ctx, int64(srcPtr))

	for i, char := range sSrc {
		// Check bounds/grow
		target := appendPtr + i
		if target >= len(ctx.heap) {
			ctx.heap = append(ctx.heap, CharLiteral(char))
		} else {
			ctx.heap[target] = CharLiteral(char)
		}
	}
	// Null terminate
	target := appendPtr + len(sSrc)
	if target >= len(ctx.heap) {
		ctx.heap = append(ctx.heap, CharLiteral(0))
	} else {
		ctx.heap[target] = CharLiteral(0)
	}

	push(ctx, destPtrVal)
}

func nativeStrlen(ctx *RuntimeContext) {
	ptrVal := pop(ctx)
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("strlen pointer must be integer"))
	}
	s := getStringFromHeap(ctx, ptrVal.valueInt)
	push(ctx, IntLiteral(int64(len(s))))
}

func nativeAssert(ctx *RuntimeContext) {
	val := pop(ctx)
	if val.Type() == LiteralInt {
		if val.valueInt == 0 {
			panic(ctx.CurrentInstruction.Error("assertion failed"))
		}
	}
}
