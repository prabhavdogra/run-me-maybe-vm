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
	{
		name: "macro single @imp",
		program: `@imp "stddefs.tash"
				@imp "linuxsyscalls.tash"
				@def N 100
				push N
				print

				push STDIN
				push STDOUT
				push STDERR
				print
				print
				print

				; syscall numbers
				push SYS_READ
				push SYS_WRITE
				push SYS_OPEN
				push SYS_EXIT
				print
				print
				print
				print`,
		additionalFiles: map[string]string{
			"stddefs.tash": `@def STDERR 2
						@def STDOUT 1
						@def STDIN 0`,
			"linuxsyscalls.tash": `@def SYS_READ 0
						@def SYS_WRITE 1
						@def SYS_OPEN 2
						@def SYS_EXIT 60`,
		},
		expected: []string{"INT 100", "INT 2", "INT 1", "INT 0", "INT 60", "INT 2", "INT 1", "INT 0"},
	},
	{
		name: "macro instruction @def",
		program: `@def PUSHTWO "push 2"
				PUSHTWO
				print`,
		expected: []string{"INT 2"},
	},
	{
		name: "macro missing def @imp (error)",
		program: `@imp "test2.tash"
					print`,
		additionalFiles: map[string]string{
			"test2.tash": `push N`,
		},
		expectedError: "expected integer, float, or char value after 'push' instruction, but found label 'N'",
	},
	{
		name: "duplicate @def (error)",
		program: `@imp "test2.tash"
					@def X 10`,
		additionalFiles: map[string]string{
			"test2.tash": `@def X 5`,
		},
		expectedError: "duplicate macro definition found for macro ''",
	},
}
