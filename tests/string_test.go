package tests

// Test basic string push and write
var stringPushTest = ProgramTestCase{
	name: "string_push",
	program: `push "Hello, World!"
	write 1 13
	halt`,
	expected: []string{"Hello, World!"},
}

// Test string with escape sequences
var stringEscapeTest = ProgramTestCase{
	name: "string_escape",
	program: `push "Line 1\nLine 2\tTabbed"
	write 1 20
	halt`,
	expected: []string{"Line 1\nLine 2\tTabbed"},
}

var stringEscapeTestMultiLine = ProgramTestCase{
	name: "string_escape_multiline",
	program: `push "Line 1\nLine 2"
	write 1 13
	halt`,
	expected: []string{"Line 1", "Line 2"},
}

// Test string duplicates (char level)
// "hi" -> 'h', 'i'. dup -> 'h', 'i', 'i'.
var stringDupTest = ProgramTestCase{
	name: "string_dup",
	program: `push "hi"
	dup
	write 1 3
	halt`,
	expected: []string{"hii"},
}

// Test string swap (char level)
// "hi" -> 'h', 'i'. swap -> 'i', 'h'.
var stringSwapTest = ProgramTestCase{
	name: "string_swap",
	program: `push "hi"
	swap
	write 1 2
	halt`,
	expected: []string{"ih"},
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
	program: `push "error_msg"
	write 2 9
	halt`,
	expected:       []string{},            // Stdout should be empty
	expectedStderr: []string{"error_msg"}, // Stderr should contain this
}

// Test error: write with invalid file descriptor type (Parser catches this first)
var stringWriteInvalidFDTypeTest = ProgramTestCase{
	name: "write_invalid_fd_type",
	program: `push 1.5
write 1.5 1
halt`,
	expectedError: "expected integer value (fd) after 'write' instruction",
}

// Test error: write with invalid file descriptor value (3) (Runtime check in instruction_helper during generation or instruction.go)
// Wait, the instruction_helper.go actually validates this during generation phase!
var stringWriteInvalidFDValueTest = ProgramTestCase{
	name: "write_invalid_fd_value",
	program: `push 3
write 3 1
halt`,
	expectedError: "invalid integer value for write fd",
}

// Test error: write expects characters on stack (found Int) (Runtime error)
var stringWriteInvalidStackTest = ProgramTestCase{
	name: "write_invalid_stack_content",
	program: `push 123
write 1 1
halt`,
	expectedError: "write expects characters on stack",
}

// Test macro import and usage with write
var stringMacroImportTest = ProgramTestCase{
	name: "string_macro_import",
	program: `@imp "stddefs.tash"
		push "Hello, world!\n"
		write STDOUT 14
		halt`,
	additionalFiles: map[string]string{
		"stddefs.tash": "@def STDOUT 1",
	},
	expected: []string{"Hello, world!"},
}

var stringTests = []ProgramTestCase{
	stringPushTest,
	stringEscapeTestMultiLine,
	stringDupTest,
	stringSwapTest,
	stringEqualTestTrue,
	stringPrintLastCharTest,
	stringWriteStderrTest,
	stringWriteInvalidFDTypeTest,
	stringWriteInvalidFDValueTest,
	stringWriteInvalidStackTest,
	stringMacroImportTest,
}
