package main

func pushIns(value int) Instruction {
	return Instruction{instructionType: InstructionPush, value: value}
}

func popIns() Instruction {
	return Instruction{instructionType: InstructionPop}
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
