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

func printStack(machine *Machine) {
	fmt.Println("Stack contents:")
	for i := 0; i < stackSize+1; i++ {
		fmt.Printf("[%d]: %d\n", i, stack[i])
	}
}
