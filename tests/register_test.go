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
			push r16
		`,
		expectedError: "register index out of bounds: r16",
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
	{
		name: "native_pow",
		program: `
			@imp "stddefs.rmm"
			push 3 ; power
			push 2 ; base
			pow
			print  ; INT 8
		`,
		expected:        []string{"INT 8"},
		additionalFiles: StdDefs,
	},
	{
		name: "native_ftoa",
		program: `
			@imp "stddefs.rmm"
			push 3.14159
			float_to_str
			push 1
			native 1 ; write
		`,
		expected:        []string{"3.14159000"},
		additionalFiles: StdDefs,
	},
	{
		name: "registers_mov_imm_types",
		program: `
			mov r0 42
			push r0
			print       ; INT 42
			mov r1 3.14
			push r1
			print       ; FLOAT 3.140000
			mov r2 'Z'
			push r2
			print       ; CHAR Z
		`,
		expected: []string{"INT 42", "FLOAT 3.140000", "CHAR Z"},
	},
	{
		name: "registers_mov_start",
		program: `
			mov r0 100
			push r0
			print       ; INT 100
		`,
		expected: []string{"INT 100"},
	},
}
