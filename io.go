package main

import (
	"encoding/binary"
	"fmt"
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

	buf := make([]byte, 8*len(machine.instructions))
	for i := 0; i < len(machine.instructions); i++ {
		off := i * 8
		binary.LittleEndian.PutUint32(buf[off:off+4], uint32(machine.instructions[i].instructionType))
		binary.LittleEndian.PutUint32(buf[off+4:off+8], uint32(int32(machine.instructions[i].value)))
	}

	if _, err := f.Write(buf); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write to file %s: %v\n", filePath, err)
		os.Exit(1)
	}
}
