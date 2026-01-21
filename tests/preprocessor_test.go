package tests

// PreprocessorTests contains test cases for @imp and @def directives
var PreprocessorTests = []ProgramTestCase{
	{
		name: "basic @def and @imp",
		program: `@imp "macros.wm"
					push KEY
					print
					halt`,
		additionalFiles: map[string]string{
			"macros.wm": "@def KEY 123",
		},
		expected: []string{"INT 123"},
	},
	{
		name: "multiple @def in one file",
		program: `@imp "defs.wm"
					push FIRST
					push SECOND
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"defs.wm": `@def FIRST 10
						@def SECOND 20`,
		},
		expected: []string{"INT 30"},
	},
	{
		name: "nested @imp",
		program: `@imp "middle.wm"
					push VALUE
					push MULTIPLIER
					mul
					print
					halt`,
		additionalFiles: map[string]string{
			"base.wm": "@def VALUE 42",
			"middle.wm": `@imp "base.wm"
							@def MULTIPLIER 2`,
		},
		expected: []string{"INT 84"},
	},
	{
		name: "@def with expression value",
		program: `@imp "config.wm"
					push MAX
					push 50
					sub
					print
					halt`,
		additionalFiles: map[string]string{
			"config.wm": "@def MAX 100",
		},
		expected: []string{"INT 50"},
	},
	{
		name: "multiple @imp in one file",
		program: `@imp "constants1.wm"
					@imp "constants2.wm"
					push A
					push B
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"constants1.wm": "@def A 5",
			"constants2.wm": "@def B 7",
		},
		expected: []string{"INT 12"},
	},
}
