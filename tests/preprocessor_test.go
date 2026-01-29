package tests

// PreprocessorTests contains test cases for @imp and @def directives
var PreprocessorTests = []ProgramTestCase{
	{
		name: "basic @def and @imp",
		program: `@imp "macros.rmm"
					push KEY
					print
					halt`,
		additionalFiles: map[string]string{
			"macros.rmm": "@def KEY 123",
		},
		expected: []string{"INT 123"},
	},
	{
		name: "multiple @def in one file",
		program: `@imp "defs.rmm"
					push FIRST
					push SECOND
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"defs.rmm": `@def FIRST 10
						@def SECOND 20`,
		},
		expected: []string{"INT 30"},
	},
	{
		name: "nested @imp",
		program: `@imp "middle.rmm"
					push VALUE
					push MULTIPLIER
					mul
					print
					halt`,
		additionalFiles: map[string]string{
			"base.rmm": "@def VALUE 42",
			"middle.rmm": `@imp "base.rmm"
							@def MULTIPLIER 2`,
		},
		expected: []string{"INT 84"},
	},
	{
		name: "@def with expression value",
		program: `@imp "config.rmm"
					push MAX
					push 50
					sub
					print
					halt`,
		additionalFiles: map[string]string{
			"config.rmm": "@def MAX 100",
		},
		expected: []string{"INT 50"},
	},
	{
		name: "multiple @imp in one file",
		program: `@imp "constants1.rmm"
					@imp "constants2.rmm"
					push A
					push B
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"constants1.rmm": "@def A 5",
			"constants2.rmm": "@def B 7",
		},
		expected: []string{"INT 12"},
	},
	{
		name: "macro single @imp",
		program: `@imp "stddefs.rmm"
				@imp "linuxsyscalls.rmm"
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
			"stddefs.rmm": `@def STDERR 2
						@def STDOUT 1
						@def STDIN 0`,
			"linuxsyscalls.rmm": `@def SYS_READ 0
						@def SYS_WRITE 1
						@def SYS_OPEN 2
						@def SYS_EXIT 60`,
		},
		expected: []string{"INT 100", "INT 2", "INT 1", "INT 0", "INT 60", "INT 2", "INT 1", "INT 0"},
	},
	{
		name: "preprocessor_entrypoint_at_end",
		program: `
			@imp "stddefs.rmm"
			
			main:
			push 5
			print
			halt
			
			entrypoint main
		`,
		expected:        []string{"INT 5"},
		additionalFiles: StdDefs,
	},
	{
		name: "preprocessor_nested_imports",
		program: `
			@imp "nested.rmm"
			main:
			call nested_func
			print
			halt
			entrypoint main
		`,
		additionalFiles: map[string]string{
			"nested.rmm": `
				@imp "std.rmm"
				nested_func:
				push 10
				ret
			`,
			"std.rmm": StdDefs["std.rmm"],
		},
		expected: []string{"INT 10"},
	},
	{
		name: "macro instruction @def",
		program: `@def PUSHTWO push 2
				PUSHTWO
				print`,
		expected: []string{"INT 2"},
	},
	{
		name: "macro missing def @imp (error)",
		program: `@imp "test2.rmm"
					print`,
		additionalFiles: map[string]string{
			"test2.rmm": `push N`,
		},
		expectedError: "expected integer, float, char, or string value after 'push' instruction, but found label 'N'",
	},
	{
		name: "duplicate @def (error)",
		program: `@imp "test2.rmm"
					@def X 10`,
		additionalFiles: map[string]string{
			"test2.rmm": `@def X 5`,
		},
		expectedError: "duplicate macro definition found for macro 'X'",
	},
}
