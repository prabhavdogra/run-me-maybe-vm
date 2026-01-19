package tests

var label2 = ProgramTestCase{
	name: "label",
	program: `jmp again
		again:
		push 5
		jmp main

		main:
		push 10
		add
		print
		`,
	expected: []string{"INT 15"},
}
