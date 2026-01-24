package tests

var explanatoryTest = ProgramTestCase{
	name: "string_test2",
	program: `@imp "stddefs.wm"
		push_str "hello\n"
		push_str "world\n"
		get_str 1     ; Pushes ptr to "hello"
		push STDOUT   ; Push File Descriptor (fd) 1 (Stdout)
		write         ; Execute write(fd, ptr)`,
	expected: []string{"world"},
	additionalFiles: map[string]string{
		"stddefs.wm": `@def STDOUT 1
		@def STDIN 0
		@def open native 0
		@def write native 1
		@def read native 2
		@def close native 3
		@def free native 4
		@def malloc native 5
		@def exit native 6

		@def RONLY 0
		@def WONLY 1
		@def RDWR 2
		@def CREAT 64
		@def EXCL 128
		`,
	},
}

// Test basic string push and write
var stringPushTest = ProgramTestCase{
	name: "string_push",
	program: `push_str "Hello, World!"
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"Hello, World!"},
}

// Test string with escape sequences
var stringEscapeTest = ProgramTestCase{
	name: "string_escape",
	program: `push_str "Line 1\nLine 2\tTabbed"
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"Line 1", "Line 2\tTabbed"},
}

var stringEscapeTestMultiLine = ProgramTestCase{
	name: "string_escape_multiline",
	program: `push_str "Line 1\nLine 2"
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"Line 1", "Line 2"},
}

// Test string duplicates (char level)
// "hi" -> 'h', 'i'. dup -> 'h', 'i', 'i'.
var stringDupTest = ProgramTestCase{
	name: "string_dup",
	program: `push "hi"
		dup
		print
		print
		print
		halt`,
	expected: []string{"CHAR i", "CHAR i", "CHAR h"},
}

// Test string swap (char level)
// "hi" -> 'h', 'i'. swap -> 'i', 'h'.
var stringSwapTest = ProgramTestCase{
	name: "string_swap",
	program: `push "hi"
		swap
		print
		print
		halt`,
	expected: []string{"CHAR h", "CHAR i"},
}

// Test string char equality
var stringEqualTestTrue = ProgramTestCase{
	name: "string_equal_true",
	program: `push "a"
		push "a"
		cmpe
		print
		halt`,
	expected: []string{"INT 1"},
}

// Test print only prints the last character
var stringPrintLastCharTest = ProgramTestCase{
	name: "string_print_last_char",
	program: `push "abc"
		print
		halt`,
	expected: []string{"CHAR c"},
}

// Test write to stderr (fd 2) -> Verified via expectedStderr
var stringWriteStderrTest = ProgramTestCase{
	name: "string_write_stderr",
	program: `push_str "error_msg"
		get_str 0
		push 2
		native 1
		halt`,
	expected:       []string{},            // Stdout should be empty
	expectedStderr: []string{"error_msg"}, // Stderr should contain this
}

// Test error: write with invalid file descriptor type (Parser catches this first)
// Test error: write with invalid file descriptor type
var stringWriteInvalidFDTypeTest = ProgramTestCase{
	name: "write_invalid_fd_type",
	program: `push_str "test"
		get_str 0
		push 1.5
		native 1
		halt`,
	expectedError: "write fd must be integer",
}

// Test error: write with invalid file descriptor value (3) (Runtime check in instruction_helper during generation or instruction.go)
// Wait, the instruction_helper.go actually validates this during generation phase!
// Test error: write string pointer must be integer
var stringWriteInvalidPointerTypeTest = ProgramTestCase{
	name: "write_invalid_pointer_type",
	program: `push 1.5
		push 1
		native 1
		halt`,
	expectedError: "write string pointer must be integer",
}

// Test error: write expects characters on stack (found Int) (Runtime error)
// Test error: write with invalid heap pointer (segmentation fault)
var stringWriteInvalidHeapPointerTest = ProgramTestCase{
	name: "write_invalid_heap_pointer",
	program: `push 999
		push 1
		native 1
		halt`,
	expectedError: "segmentation fault",
}

// Test macro import and usage with write
var stringMacroImportTest = ProgramTestCase{
	name: "string_macro_import",
	program: `@imp "stddefs.wm"
		push_str "Hello, world!\n"
		get_str 0
		push STDOUT
		write
		halt`,
	additionalFiles: map[string]string{
		"stddefs.wm": "@def STDOUT 1\n@def write native 1",
	},
	expected: []string{"Hello, world!"},
}

var stringLengthTest = ProgramTestCase{
	name: "string_length",
	program: `push 0             ; Sentinel 0 at bottom of stack
		push "hello world" ; Push characters
		push 0             ; Initial length accumulator
		jmp length_loop

		length_loop:
			swap           ; Swap accumulator and next character: [..., char, acc] -> [..., acc, char]
			dup            ; Check if char is 0 (sentinel)
			push 0
			cmpe
			nzjmp length_done
			pop            ; Clean up cmpe result (0)
			pop            ; Clean up operand (0)
			pop            ; Clean up operand (char)
			
			; At this point stack is [..., acc]
			; We popped the char, so we consumed it.
			push 1
			add            ; Increment accumulator
			jmp length_loop

		length_done:
			pop            ; Clean up cmpe operand (0)
			pop            ; Clean up cmpe operand (0)
			pop            ; Clean up stack sentinel (0)
			push 11        ; Expected length
			cmpe
			print
			halt`,
	expected: []string{"INT 1"},
}

var stringTests = []ProgramTestCase{
	explanatoryTest,
	stringPushTest,
	stringEscapeTest,
	stringEscapeTestMultiLine,
	stringDupTest,
	stringSwapTest,
	stringEqualTestTrue,
	stringPrintLastCharTest,
	stringWriteStderrTest,
	stringWriteInvalidFDTypeTest,
	stringWriteInvalidPointerTypeTest,
	stringWriteInvalidHeapPointerTest,
	stringMacroImportTest,
	stringLengthTest,
}
