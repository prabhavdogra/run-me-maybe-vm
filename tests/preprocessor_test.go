package tests

// PreprocessorTests contains test cases for @imp and @def directives
var PreprocessorTests = []ProgramTestCase{
	{
		name: "basic @def and @imp",
		program: `@imp "macros.cmm"
					push KEY
					print
					halt`,
		additionalFiles: map[string]string{
			"macros.cmm": "@def KEY 123",
		},
		expected: []string{"INT 123"},
	},
	{
		name: "multiple @def in one file",
		program: `@imp "defs.cmm"
					push FIRST
					push SECOND
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"defs.cmm": `@def FIRST 10
						@def SECOND 20`,
		},
		expected: []string{"INT 30"},
	},
	{
		name: "nested @imp",
		program: `@imp "middle.cmm"
					push VALUE
					push MULTIPLIER
					mul
					print
					halt`,
		additionalFiles: map[string]string{
			"base.cmm": "@def VALUE 42",
			"middle.cmm": `@imp "base.cmm"
							@def MULTIPLIER 2`,
		},
		expected: []string{"INT 84"},
	},
	{
		name: "@def with expression value",
		program: `@imp "config.cmm"
					push MAX
					push 50
					sub
					print
					halt`,
		additionalFiles: map[string]string{
			"config.cmm": "@def MAX 100",
		},
		expected: []string{"INT 50"},
	},
	{
		name: "multiple @imp in one file",
		program: `@imp "constants1.cmm"
					@imp "constants2.cmm"
					push A
					push B
					add
					print
					halt`,
		additionalFiles: map[string]string{
			"constants1.cmm": "@def A 5",
			"constants2.cmm": "@def B 7",
		},
		expected: []string{"INT 12"},
	},
	{
		name: "macro single @imp",
		program: `@imp "stddefs.cmm"
				@imp "linuxsyscalls.cmm"
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
			"stddefs.cmm": `@def STDERR 2
						@def STDOUT 1
						@def STDIN 0`,
			"linuxsyscalls.cmm": `@def SYS_READ 0
						@def SYS_WRITE 1
						@def SYS_OPEN 2
						@def SYS_EXIT 60`,
		},
		expected: []string{"INT 100", "INT 2", "INT 1", "INT 0", "INT 60", "INT 2", "INT 1", "INT 0"},
	},
	{
		name: "preprocessor_entrypoint_at_end",
		program: `
			@imp "std.cmm"
			
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
			@imp "nested.cmm"
			main:
			call nested_func
			print
			halt
			entrypoint main
		`,
		additionalFiles: map[string]string{
			"nested.cmm": `
				@imp "std.cmm"
				nested_func:
				push 10
				ret
			`,
			"std.cmm": StdDefs["std.cmm"],
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
		program: `@imp "test2.cmm"
					print`,
		additionalFiles: map[string]string{
			"test2.cmm": `push N`,
		},
		expectedError: "expected integer, float, char, or string value after 'push' instruction, but found label 'N'",
	},
	{
		name: "duplicate @def (error)",
		program: `@imp "test2.cmm"
					@def X 10`,
		additionalFiles: map[string]string{
			"test2.cmm": `@def X 5`,
		},
		expectedError: "duplicate macro definition found for macro 'X'",
	},
}
