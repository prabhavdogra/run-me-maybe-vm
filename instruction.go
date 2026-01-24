package main

import (
	"fmt"
	"io"
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
			case 0:
				nativeOpen(machine, instr)
			case 1:
				nativeWrite(machine, instr)
			case 2:
				nativeRead(machine, instr)
			case 3:
				nativeClose(machine, instr)
			case 4:
				nativeFree(machine, instr)
			case 5:
				nativeMalloc(machine, instr)
			case 6:
				nativeExit(machine, instr)
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

func nativeOpen(machine *Machine, instr Instruction) {
	// Pop flags
	flagsVal := pop(machine)
	if flagsVal.Type() != LiteralInt {
		panic(instr.Error("open flags must be integer"))
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
	lenVal := pop(machine)
	if lenVal.Type() != LiteralInt {
		panic(instr.Error("open filename length must be integer"))
	}
	length := int(lenVal.valueInt)

	// Pop filename pointer
	ptrVal := pop(machine)
	if ptrVal.Type() != LiteralInt {
		panic(instr.Error("open filename pointer must be integer"))
	}
	ptr := ptrVal.valueInt

	// Read filename from heap
	if _, ok := machine.heap[ptr]; !ok {
		panic(instr.Error("segmentation fault: invalid heap pointer for filename"))
	}
	if length > len(machine.heap[ptr]) {
		panic(instr.Error("buffer overflow: filename length exceeds allocated size"))
	}
	filenameChars := machine.heap[ptr][:length]
	filename := ""
	for _, charLit := range filenameChars {
		if charLit.Type() != LiteralChar {
			panic(instr.Error("filename must be a string of characters"))
		}
		filename += string(charLit.valueChar)
	}

	// Open the file
	file, err := os.OpenFile(filename, osFlags, 0644)
	if err != nil {
		panic(instr.Error(fmt.Sprintf("failed to open file %s: %v", filename, err)))
	}

	// 0: Standard Input (stdin)
	// 1: Standard Output (stdout)
	// 2: Standard Error (stderr)

	// Find lowest available file descriptor
	fd := int64(3)
	for {
		if _, ok := machine.fileDescriptors[fd]; !ok {
			break
		}
		fd++
	}

	machine.fileDescriptors[fd] = file
	push(machine, IntLiteral(fd))
}

func nativeWrite(machine *Machine, instr Instruction) {
	lenVal := pop(machine)
	if lenVal.Type() != LiteralInt {
		panic(instr.Error("write length must be integer"))
	}
	length := int(lenVal.valueInt)

	fdVal := pop(machine)
	if fdVal.Type() != LiteralInt {
		panic(instr.Error("write fd must be integer"))
	}
	fd := int64(fdVal.valueInt)

	var writer io.Writer
	if fd == 1 {
		writer = machine.output
	} else if fd == 2 {
		writer = os.Stderr
	} else {
		if file, ok := machine.fileDescriptors[fd]; ok {
			writer = file
		} else {
			panic(instr.Error(fmt.Sprintf("unknown file descriptor %d", fd)))
		}
	}

	s := ""
	for i := 0; i < length; i++ {
		val := pop(machine)
		if val.Type() != LiteralChar {
			panic(instr.Error("write expects characters on stack"))
		}
		s = string(val.valueChar) + s
	}

	fmt.Fprint(writer, s)
}

func nativeRead(machine *Machine, instr Instruction) {
	// Arguments: [..., fd, len, ptr] (Top is ptr)
	ptrVal := pop(machine)
	if ptrVal.Type() != LiteralInt {
		panic(instr.Error("read buffer pointer must be integer"))
	}
	ptr := ptrVal.valueInt

	lenVal := pop(machine)
	if lenVal.Type() != LiteralInt {
		panic(instr.Error("read length must be integer"))
	}
	length := int(lenVal.valueInt)

	fdVal := pop(machine)
	if fdVal.Type() != LiteralInt {
		panic(instr.Error("read fd must be integer"))
	}
	fd := int64(fdVal.valueInt)

	var reader io.Reader
	if fd == 0 {
		reader = machine.input
	} else {
		if file, ok := machine.fileDescriptors[fd]; ok {
			reader = file
		} else {
			panic(instr.Error(fmt.Sprintf("read error: invalid file descriptor %d", fd)))
		}
	}

	// Validate Heap Pointer
	if _, ok := machine.heap[ptr]; !ok {
		panic(instr.Error("segmentation fault: invalid heap pointer"))
	}
	// Validate Size against allocation
	if length > len(machine.heap[ptr]) {
		panic(instr.Error("buffer overflow: read length exceeds allocated size"))
	}

	// Read from Input
	buf := make([]byte, length)
	_, err := reader.Read(buf)
	if err != nil && err != io.EOF {
		panic(instr.Error(fmt.Sprintf("read error: %v", err)))
	}

	// Store in Heap
	for i, b := range buf {
		machine.heap[ptr][i] = CharLiteral(rune(b))
	}
}

func nativeClose(machine *Machine, instr Instruction) {
	// Pop file descriptor ID
	fdVal := pop(machine)
	if fdVal.Type() != LiteralInt {
		panic(instr.Error("close file descriptor must be integer"))
	}
	fd := int64(fdVal.valueInt)

	// Check if it's a valid custom file descriptor
	if fd < 3 { // 0, 1, 2 are stdin, stdout, stderr - cannot close
		panic(instr.Error(fmt.Sprintf("cannot close standard file descriptor %d", fd)))
	}

	file, ok := machine.fileDescriptors[fd]
	if !ok {
		panic(instr.Error(fmt.Sprintf("invalid file descriptor %d", fd)))
	}

	err := file.Close()
	if err != nil {
		panic(instr.Error(fmt.Sprintf("failed to close file %d: %v", fd, err)))
	}

	delete(machine.fileDescriptors, fd)
}

func nativeFree(machine *Machine, instr Instruction) {

	// Pop ptr
	ptrVal := pop(machine)
	if ptrVal.Type() != LiteralInt {
		panic(instr.Error("free pointer must be integer"))
	}
	ptr := ptrVal.valueInt

	if _, ok := machine.heap[ptr]; !ok {
		panic(instr.Error("double free or invalid heap pointer"))
	}
	delete(machine.heap, ptr)
}

func nativeMalloc(machine *Machine, instr Instruction) {
	// Pop size
	sizeVal := pop(machine)
	if sizeVal.Type() != LiteralInt {
		panic(instr.Error("malloc size must be integer"))
	}
	size := int(sizeVal.valueInt)

	// Allocate
	ptr := machine.heapPtr
	machine.heap[ptr] = make([]Literal, size)
	machine.heapPtr++

	// Push Pointer
	push(machine, IntLiteral(ptr))
}

func nativeExit(machine *Machine, instr Instruction) {
	// Pop exit code
	codeVal := pop(machine)
	if codeVal.Type() != LiteralInt {
		panic(instr.Error("exit code must be integer"))
	}
	code := int(codeVal.valueInt)
	os.Exit(code)
}
