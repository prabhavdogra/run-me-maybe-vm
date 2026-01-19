package tests

var label = ProgramTestCase{
	name: "label",
	program: `jmp main

		test:
		push 10
		push 4
		add
		print
		jmp end

		main:
		push 5
		push 7
		add
		print
		jmp test

		end:
		push 2
		push 5
		add
		print
		jmp realend

		realend:
		push 15
		print
		jmp theend

		theend:
		push 99
		print`,
	expected: []string{"INT 12", "INT 14", "INT 7", "INT 15", "INT 99"},
}
