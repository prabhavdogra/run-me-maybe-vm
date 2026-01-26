package tests

var scanfTests = []ProgramTestCase{
	{
		name: "scanf_basic",
		program: `@imp "stddefs.wm"
			push 10		; [size]
			malloc      ; mallon(size) 
						; Stack: [ptr]
			dup         ; Stack: [ptr, ptr]
			scanf       ; read "Hello" into ptr
			            ; Stack: [ptr]
			push 1      ; stdout
			            ; Stack: [ptr, 1]
			write       ; writes string at ptr to stdout
		`,
		input:           "Hello",
		expected:        []string{"Hello"},
		additionalFiles: StdDefs,
	},
}
