package tests

var registerTests = []ProgramTestCase{
	{
		name: "registers_basic",
		program: `
			mov r0 123
			push r0
			print       ; 123
			mov r1 'A'
			push r1
			print       ; 65
			mov r2 0
			mov r3 999
			push r2
			push r3
			add
			print       ; 999
		`,
		expected: []string{"123", "65", "999"},
	},
	{
		name: "registers_float",
		program: `
			mov r0 3.14
			push r0
			print
		`,
		expected: []string{"3.140000"}, // Float printing format check
	},
	{
		name: "registers_overwrite",
		program: `
			mov r0 10
			push r0
			mov r0 20
			push r0
			add
			print       ; 30
		`,
		expected: []string{"30"},
	},
	{
		name: "registers_error_missing_val",
		program: `
			mov r0
		`,
		expectedError: "expected immediate value after register in mov",
	},
	{
		name: "registers_error_invalid_reg",
		program: `
			push r4
		`,
		expectedError: "unknown register r4", // Parser catches this likely
	},
}
