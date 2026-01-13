package main

var program = []Instruction{
	pushIns(5),
	pushIns(10),
	cmpgIns(),
	zjmpIns(5),
	haltIns(),
	pushIns(42),
}

func push(machine *Machine, value int64) {
	if len(machine.stack) >= maxStackSize {
		panic("ERROR: stack overflow")
	}
	machine.stack = append(machine.stack, value)
}

func pop(machine *Machine) int64 {
	if len(machine.stack) == 0 {
		panic("ERROR: stack underflow")
	}
	value := machine.stack[len(machine.stack)-1]
	machine.stack = machine.stack[:len(machine.stack)-1]
	return value
}

func main() {
	loadedMachine := &Machine{
		stack:        make([]int64, 0, maxStackSize),
		instructions: program,
	}
	writeProgram(loadedMachine, "program.bin")
	loadedMachine = readProgram("program.bin")
	runInstructions(loadedMachine)
}
