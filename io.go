package main

import (
	"encoding/binary"
	"fmt"
	"math"
	"os"
)

// writeProgram writes instructions to a file using the 8-byte-per-instruction
// format (4-byte type, 4-byte value, both little-endian).
func writeProgram(machine *Machine, filePath string) {
	f, err := os.Create(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not write to file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	// Calculate buffer size: 8 bytes for entrypoint + 16 bytes per instruction
	buf := make([]byte, 8+16*len(machine.instructions))

	// Write entrypoint (8 bytes) at the beginning
	binary.LittleEndian.PutUint64(buf[0:8], uint64(machine.entrypoint))

	for i := 0; i < len(machine.instructions); i++ {
		off := 8 + i*16 // Offset by 8 bytes for entrypoint

		instr := machine.instructions[i]

		// Write instruction type (4 bytes)
		binary.LittleEndian.PutUint32(buf[off:off+4], uint32(instr.instructionType))

		// Write literal type (1 byte at offset 4)
		buf[off+4] = uint8(instr.value.valueType)

		// Write value (8 bytes starting at offset 8)
		switch instr.value.valueType {
		case LiteralInt:
			binary.LittleEndian.PutUint64(buf[off+8:off+16], uint64(instr.value.valueInt))
		case LiteralFloat:
			binary.LittleEndian.PutUint64(buf[off+8:off+16], math.Float64bits(instr.value.valueFloat))
		case LiteralChar:
			binary.LittleEndian.PutUint64(buf[off+8:off+16], uint64(instr.value.valueChar))
		case LiteralNone:
			binary.LittleEndian.PutUint64(buf[off+8:off+16], 0)
		case LiteralPointer:
			binary.LittleEndian.PutUint64(buf[off+8:off+16], uint64(instr.value.valuePtr))
		}
	}

	if _, err := f.Write(buf); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write to file %s: %v\n", filePath, err)
		os.Exit(1)
	}
}
