package tests

// LineNumberTests verifies that panics report the correct file and line number.
var LineNumberTests = []ProgramTestCase{
	{
		name: "panic_runtime_lines",
		program: `push 1
		push 2
		add
		pop
		pop  ; 1st element
		pop  ; Stack underflow here (Line 5)
		halt`,
		expectedError:   "ERROR (main.wm:5): stack underflow",
		additionalFiles: map[string]string{},
	},
	{
		name: "panic_macro_lines",
		program: `@imp "stddefs.wm"
		push "w"
		open
		halt`,
		expectedError:   "ERROR (main.wm:3): open flags must be integer",
		additionalFiles: StdDefs,
	},
	{
		name: "panic_multiline_macro_lines",
		program: `@imp "stddefs.wm"
		
		
		push "w"
		open
		halt`,
		expectedError:   "ERROR (main.wm:5): open flags must be integer",
		additionalFiles: StdDefs,
	},
}
