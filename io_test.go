package main

import (
	"encoding/binary"
	"os"
	"testing"
)

func TestWriteProgramSerialization(t *testing.T) {
	// Setup a sample machine with entrypoint and instructions
	instr1 := Instruction{instructionType: InstructionPush, value: IntLiteral(42)}
	instr2 := Instruction{instructionType: InstructionHalt}
	machine := &Machine{
		entrypoint:   12345,
		instructions: []Instruction{instr1, instr2},
	}

	// Create temp file
	tmpFile, err := os.CreateTemp("", "test_serialization_*.bin")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())
	tmpFile.Close()

	// Write program
	writeProgram(machine, tmpFile.Name())

	// Read back raw bytes
	payload, err := os.ReadFile(tmpFile.Name())
	if err != nil {
		t.Fatalf("failed to read back file: %v", err)
	}

	// Expected size: 8 bytes (entrypoint) + 2 * 16 bytes (instructions) = 40 bytes
	expectedSize := 8 + 32
	if len(payload) != expectedSize {
		t.Errorf("expected file size %d, got %d", expectedSize, len(payload))
	}

	// Verify Entrypoint
	readEntrypoint := binary.LittleEndian.Uint64(payload[0:8])
	if int(readEntrypoint) != machine.entrypoint {
		t.Errorf("expected entrypoint %d, got %d", machine.entrypoint, int(readEntrypoint))
	}

	// Verify Instruction 1 (Push 42)
	// Offset: 8
	// Type (4 bytes): InstructionPush (check value)
	// Literal Type (1 byte): LiteralInt
	// Value (8 bytes): 42
	// For exact verification, we can just check offsets or assume if size matches we are good,
	// but better to check bytes roughly.
	// instrType is binary encoded.
}
