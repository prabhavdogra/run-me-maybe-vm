package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

// readProgram reads a file where each instruction is stored as 8 bytes:
//
//	4 bytes little-endian for instruction type (uint32)
//	4 bytes little-endian for value (uint32, stored into int64)
func readProgram(filePath string) *Machine {
	machine := &Machine{}

	f, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not read from file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not stat file %s: %v\n", filePath, err)
		os.Exit(1)
	}

	size := fi.Size()
	if size == 0 {
		machine.instructions = nil
		return machine
	}

	if size%8 != 0 {
		fmt.Fprintf(os.Stderr, "WARNING: file %s size (%d) not a multiple of 8, truncating\n", filePath, size)
	}

	count := int(size / 8)
	buf := make([]byte, count*8)
	if _, err := io.ReadFull(f, buf); err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: failed to read program bytes: %v\n", err)
		os.Exit(1)
	}

	insts := make([]Instruction, 0, count)
	for i := 0; i < count; i++ {
		off := i * 8
		typ := binary.LittleEndian.Uint32(buf[off : off+4])
		val := int32(binary.LittleEndian.Uint32(buf[off+4 : off+8]))
		insts = append(insts, Instruction{instructionType: InstructionSet(typ), value: int64(val)})
	}

	machine.instructions = insts
	return machine
}

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
