package tests

var StrlenTest = ProgramTestCase{
	name: "strlen",
	program: `@imp "stddefs.wm"
		push_str "Any string here and you'll get the length!\n"
		jmp main
		
		my_strlen:
			push_str "\0"
			push 0
			swap
			loop:
			dup
			get_str 1
			native 90
			swap
			push 1
			add
			swap
			inswap 0
			push 1
			add
			inswap 0
			zjmp loop
			inswap 0
			push 1
			sub
			print
			jmp finish

		main:
			get_str 0
			jmp my_strlen

		finish:
			push 0
			exit`,
	expected: []string{
		"INT 43",
	},
	additionalFiles: StdDefs,
}
