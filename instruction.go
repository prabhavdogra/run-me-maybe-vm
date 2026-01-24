package main

import (
	"fmt"
	"io"
	"os"
	"vm/internal/parser"
	"vm/internal/token"
)

type InstructionSet uint8

const (
	InstructionNoOp InstructionSet = iota
	InstructionPush
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
	InstructionHalt
	InstructionIntToStr
)

func populateStringTable(parsedTokens *parser.ParserList) ([]int64, map[int64][]Literal, int64) {
	stringTable := []int64{}
	heap := make(map[int64][]Literal)
	heapPtr := int64(0)

	cur := parsedTokens
	for cur != nil {
		if cur.Value.Type == token.TypePushStr {
			if cur.Next == nil || cur.Next.Value.Type != token.TypeString {
				panic("expected string after push_str") // Should be caught by parser
			}
			strVal := cur.Next.Value.Text

			// Malloc string (len + 1 for null terminator)
			ptr := heapPtr
			heap[ptr] = make([]Literal, len(strVal)+1)
			for i, char := range strVal {
				heap[ptr][i] = CharLiteral(char)
			}
			heap[ptr][len(strVal)] = CharLiteral(0) // Null terminator
			heapPtr++

			stringTable = append(stringTable, ptr)
		}
		cur = cur.Next
	}
	return stringTable, heap, heapPtr
}

func (i InstructionSet) String() string {
	switch i {
	case InstructionNoOp:
		return "NOOP"
	case InstructionPush:
		return "PUSH"
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
	case InstructionIntToStr:
		return "INT_TO_STR"
	default:
		return fmt.Sprintf("UNKNOWN(%d)", i)
	}
}

func runInstructions(machine *Machine) *Machine {
	ctx := &RuntimeContext{Machine: machine}
	for insPtr := 0; insPtr < len(machine.instructions); insPtr++ {
		instr := machine.instructions[insPtr]
		ctx.CurrentInstruction = instr
		if debugMode {
			fmt.Fprintf(os.Stderr, "Line %d: %v, Stack: %+v\n", instr.line, instr.instructionType, ctx.stack)
		}
		switch instr.instructionType {
		case InstructionNoOp:
			// do nothing
		case InstructionPush:
			push(ctx, instr.value)
		case InstructionGetStr:
			idx := int(instr.value.valueInt)
			if idx < 0 || idx >= len(machine.stringTable) {
				panic(ctx.CurrentInstruction.Error("string index out of bounds"))
			}
			ptr := machine.stringTable[idx]
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
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if a.Equal(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpne:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if !a.Equal(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpg:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if a.Greater(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpl:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if a.Less(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmpge:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if a.GreaterOrEqual(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionCmple:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b)
			push(ctx, a)
			if a.LessOrEqual(b) {
				push(ctx, IntLiteral(1))
			} else {
				push(ctx, IntLiteral(0))
			}
		case InstructionAdd:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, a.Add(b))
		case InstructionSub:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, b.Sub(a))
		case InstructionMul:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, a.Mul(b))
		case InstructionDiv:
			a := pop(ctx)
			b := pop(ctx)
			push(ctx, a.Div(b))
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
			value := pop(ctx)
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
				insPtr = target - 1 // -1 because loop will increment
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
				// free(ptr)
				nativeFree(ctx)
			case 5:
				// malloc(size)
				nativeMalloc(ctx)
			case 6:
				// exit(code)
				nativeExit(ctx)
			default:
				panic(ctx.CurrentInstruction.Error(fmt.Sprintf("unknown native syscall ID: %d", syscallID.valueInt)))
			}
		case InstructionHalt:
			insPtr = machine.programSize()
		case InstructionIntToStr:
			value := pop(ctx)
			if value.Type() != LiteralInt {
				panic(ctx.CurrentInstruction.Error("int_to_str expects an integer"))
			}
			s := fmt.Sprintf("%d", value.valueInt)
			ptr := machine.heapPtr
			machine.heap[ptr] = make([]Literal, len(s)+1)
			for i, char := range s {
				machine.heap[ptr][i] = CharLiteral(char)
			}
			machine.heap[ptr][len(s)] = CharLiteral(0) // Null terminator
			machine.heapPtr++
			push(ctx, IntLiteral(ptr))
		default:
			panic("ERROR: unknown instruction")
		}
	}
	return machine
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
	ptr := ptrVal.valueInt

	// Read filename from heap
	if _, ok := ctx.heap[ptr]; !ok {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer for filename"))
	}
	if length > len(ctx.heap[ptr]) {
		panic(ctx.CurrentInstruction.Error("buffer overflow: filename length exceeds allocated size"))
	}
	filenameChars := ctx.heap[ptr][:length]
	filename := ""
	for _, charLit := range filenameChars {
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
	// Try to pop fd first (top of stack)
	// Usage 1: push_str "hello"; get_str 0; push 1; native 1 -> Stack: [ptr, fd]
	// Usage 2: push char; push char; push len; push fd; native 1 -> Stack: [..., char, char, len, fd] ??
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

	if buffer, ok := ctx.heap[int64(ptr.valueInt)]; ok {
		s := ""
		for _, charLit := range buffer {
			if charLit.Type() != LiteralChar {
				continue
			}
			if charLit.valueChar == 0 {
				break
			}
			s += string(charLit.valueChar)
		}
		fmt.Fprint(writer, s)
	} else {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer"))
	}
}

// Read from a file descriptor into a buffer
func nativeRead(ctx *RuntimeContext) {
	// Arguments: [..., fd, len, ptr] (Top is ptr)
	ptrVal := pop(ctx)
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("read buffer pointer must be integer"))
	}
	ptr := ptrVal.valueInt

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

	// Validate Heap Pointer
	if _, ok := ctx.heap[ptr]; !ok {
		panic(ctx.CurrentInstruction.Error("segmentation fault: invalid heap pointer"))
	}
	// Validate Size against allocation
	if length > len(ctx.heap[ptr]) {
		panic(ctx.CurrentInstruction.Error("buffer overflow: read length exceeds allocated size"))
	}

	// Read from Input
	buf := make([]byte, length)
	_, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		panic(ctx.CurrentInstruction.Error(fmt.Sprintf("read error: %v", err)))
	}

	// Store in Heap
	for i, b := range buf {
		ctx.heap[ptr][i] = CharLiteral(rune(b))
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
	if ptrVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("free pointer must be integer"))
	}
	ptr := ptrVal.valueInt

	if _, ok := ctx.heap[ptr]; !ok {
		panic(ctx.CurrentInstruction.Error("double free or invalid heap pointer"))
	}
	delete(ctx.heap, ptr)
}

func nativeMalloc(ctx *RuntimeContext) {
	// Pop size
	sizeVal := pop(ctx)
	if sizeVal.Type() != LiteralInt {
		panic(ctx.CurrentInstruction.Error("malloc size must be integer"))
	}
	size := int(sizeVal.valueInt)

	// Allocate
	ptr := ctx.heapPtr
	ctx.heap[ptr] = make([]Literal, size)
	ctx.heapPtr++

	// Push Pointer
	push(ctx, IntLiteral(ptr))
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
