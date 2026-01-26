package tests

var registerExtendedTests = []ProgramTestCase{
	{
		name: "registers_extended_r13",
		program: `
			mov r13 42
			push r13
			print       ; INT 42
		`,
		expected: []string{"INT 42"},
	},
	{
		name: "registers_extended_all_mov",
		program: `
			mov r4 4
			mov r5 5
			mov r15 15
			push r4
			push r5
			add
			push r15
			add
			print       ; INT 24 (4+5+15)
		`,
		expected: []string{"INT 24"},
	},
}
