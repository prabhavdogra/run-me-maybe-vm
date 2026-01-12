package main

var program = []Instruction{
	pushIns(5),
	pushIns(10),
	pushIns(10),
	pushIns(2),
	mulIns(),
	dupIns(),
	printIns(),
}

func push(machine *Machine, value int64) {
	if stackSize >= maxStackSize {
		panic("ERROR: stack overflow")
	}
	stack[stackSize] = value
	stackSize++
}

func pop(machine *Machine) int64 {
	if stackSize == 0 {
		panic("ERROR: stack underflow")
	}
	stackSize--
	return stack[stackSize]
}

func main() {
	loadedMachine := &Machine{}
	loadedMachine.instructions = program
	writeProgram(loadedMachine, "program.bin")
	readProgram("program.bin")
	runInstructions(loadedMachine)
}
