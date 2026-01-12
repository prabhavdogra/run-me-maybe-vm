package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func readProgram(filePath string) *Machine {
	machine := &Machine{}

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not read from file %s: %v\n", filePath, err)
		os.Exit(1)
	}
	defer file.Close()

	const instrSize = 8 // 2 × int32 in file format

	var insts []Instruction
	var buf [instrSize]byte

	for {
		_, err := io.ReadFull(file, buf[:])
		if err == io.EOF {
			break // done reading
		}
		if err == io.ErrUnexpectedEOF {
			fmt.Fprintf(os.Stderr, "WARNING: file %s has incomplete instruction, ignoring\n", filePath)
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: reading file %s: %v\n", filePath, err)
			os.Exit(1)
		}

		typ := int32(binary.LittleEndian.Uint32(buf[0:4]))
		val := int32(binary.LittleEndian.Uint32(buf[4:8]))

		insts = append(insts, Instruction{
			instructionType: InstructionSet(typ),
			value:           int64(val),
		})
	}

	machine.programSize = len(insts)
	if len(insts) > 0 {
		machine.instructions = insts
	} else {
		machine.instructions = nil
	}

	return machine
}

func writeProgram(machine *Machine, filePath string) {
	file, err := os.Create(filePath) // creates or truncates, like "wb"
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Could not write to file %s\n", filePath)
		os.Exit(1)
	}
	defer file.Close()

	var buf bytes.Buffer

	for i := 0; i < len(machine.instructions); i++ {
		err := binary.Write(&buf, binary.LittleEndian, machine.instructions[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "ERROR: Failed to write instruction\n")
			os.Exit(1)
		}
	}

	_, err = file.Write(buf.Bytes())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: Failed to write to file %s\n", filePath)
		os.Exit(1)
	}
}
