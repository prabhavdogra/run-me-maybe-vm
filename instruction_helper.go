package main

import "fmt"

func pushIns(value int64) Instruction {
	return Instruction{instructionType: InstructionPush, value: value}
}

func popIns() Instruction {
	return Instruction{instructionType: InstructionPop}
}

func dupIns() Instruction {
	return Instruction{instructionType: InstructionDup}
}

func swapIns() Instruction {
	return Instruction{instructionType: InstructionSwap}
}

func addIns() Instruction {
	return Instruction{instructionType: InstructionAdd}
}

func subIns() Instruction {
	return Instruction{instructionType: InstructionSub}
}

func mulIns() Instruction {
	return Instruction{instructionType: InstructionMul}
}

func divIns() Instruction {
	return Instruction{instructionType: InstructionDiv}
}

func printIns() Instruction {
	return Instruction{instructionType: InstructionPrint}
}

func cmpeIns() Instruction {
	return Instruction{instructionType: InstructionCmpe}
}

func cmpneIns() Instruction {
	return Instruction{instructionType: InstructionCmpne}
}

func cmpgIns() Instruction {
	return Instruction{instructionType: InstructionCmpg}
}

func cmplIns() Instruction {
	return Instruction{instructionType: InstructionCmpl}
}

func cmpgeIns() Instruction {
	return Instruction{instructionType: InstructionCmpge}
}

func cmpleIns() Instruction {
	return Instruction{instructionType: InstructionCmple}
}

func modIns() Instruction {
	return Instruction{instructionType: InstructionMod}
}

func jmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionJmp, value: target}
}

func zjmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionZjmp, value: target}
}

func nzjmpIns(target int64) Instruction {
	return Instruction{instructionType: InstructionNzjmp, value: target}
}

func haltIns() Instruction {
	return Instruction{instructionType: InstructionHalt}
}

func noopIns() Instruction {
	return Instruction{instructionType: InstructionNoOp}
}

func printStack(machine *Machine) {
	fmt.Println("------ STACK")
	for i := 0; i < len(machine.stack); i++ {
		fmt.Printf("[%d]: %d\n", i, machine.stack[i])
	}
	fmt.Println("------ END OF STACK")
}

func (machine *Machine) programSize() int {
	return len(machine.instructions)
}
