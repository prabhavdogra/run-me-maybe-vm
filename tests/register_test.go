package tests

var registerTests = []ProgramTestCase{
	{
		name: "registers_basic",
		program: `
			mov r0 123
			push r0
			print       ; INT 123
			mov r1 'A'
			push r1
			print       ; CHAR A
			mov r2 0
			mov r3 999
			push r2
			push r3
			add
			print       ; INT 999
		`,
		expected: []string{"INT 123", "CHAR A", "INT 999"},
	},
	{
		name: "registers_float",
		program: `
			mov r0 3.14
			push r0
			print
		`,
		expected: []string{"FLOAT 3.140000"}, // Float printing format check
	},
	{
		name: "registers_overwrite",
		program: `
			mov r0 10
			push r0
			mov r0 20
			push r0
			add
			print       ; INT 30
		`,
		expected: []string{"INT 30"},
	},
	{
		name: "registers_error_missing_val",
		program: `
			mov r0
		`,
		expectedError: "expected integer, float, char, or top value after register in 'mov' instruction, but found invalid",
	},
	{
		name: "registers_error_invalid_reg",
		program: `
			push r4
		`,
		expectedError: "found label 'r4'",
	},
	{
		name: "registers_mov_top",
		program: `
			push 999
			push 50
			mov r0 top
			            ; Stack: [999] (50 moved to r0)
			push r0
			print       ; INT 50
			print       ; INT 999
		`,
		expected: []string{"INT 50", "INT 999"},
	},
}
