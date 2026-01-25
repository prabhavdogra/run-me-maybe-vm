package tests

var explanatoryTest = ProgramTestCase{
	name: "string_test2",
	program: `@imp "stddefs.wm"
		push_str "hello\n"
		push_str "world\n"
		get_str 1     ; Pushes ptr to "hello"
		push STDOUT   ; Push File Descriptor (fd) 1 (Stdout)
		write         ; Execute write(fd, ptr)`,
	expected:        []string{"world"},
	additionalFiles: StdDefs,
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
	additionalFiles: StdDefs,
	expected:        []string{"Hello, world!"},
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

var intToStrTest = ProgramTestCase{
	name: "int_to_str",
	program: `@imp "stddefs.wm"
		push 12345
		int_to_str
		push STDOUT
		write
		halt`,
	expected:        []string{"12345"},
	additionalFiles: StdDefs,
}

var stringPopStrTest = ProgramTestCase{
	name: "string_pop_str",
	program: `push_str "keep"
		push_str "discard"
		pop_str
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"keep"},
}

var stringDupStrTest = ProgramTestCase{
	name: "string_dup_str",
	program: `push_str "hello"
		dup_str
		pop_str
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"hello"},
}

var stringSwapStrTest = ProgramTestCase{
	name: "string_swap_str",
	program: `push_str "first"
		push_str "second"
		swap_str
		get_str 1
		push 1
		native 1
		get_str 0
		push 1
		native 1
		halt`,
	expected: []string{"firstsecond"},
}

var stringInDupStrTest = ProgramTestCase{
	name: "string_indup_str",
	program: `push_str "bottom"
		push_str "top"
		indup_str 0
		get_str 2
		push 1
		native 1
		halt`,
	expected: []string{"bottom"},
}

var stringInSwapStrTest = ProgramTestCase{
	name: "string_inswap_str",
	program: `push_str "A"
		push_str "B"
		push_str "C"
		inswap_str 0
		get_str 0
		push 1
		native 1
		get_str 1
		push 1
		native 1
		get_str 2
		push 1
		native 1
		halt`,
	expected: []string{"CBA"},
}

var testFibTest = ProgramTestCase{
	name: "t_fib",
	program: `@imp "stddefs.wm"
	@def N 30 

	push_str "\n"
	push_str "# of iterations "

	push N 
	push 1
	push 1
	push 0

	loop:
	inswap 0
	dup
	push 0
	cmpe
	nzjmp end
	pop
	inswap 0

	indup 2
	inswap 1
	pop
	dup
	inswap 2
	pop
	indup 1
	indup 2
	add
	swap

	int_to_str
	push STDOUT
	write


	inswap 0
	push 1
	sub
	inswap 0

	get_str 0
	push STDOUT
	write

	jmp loop 

	end:

	push N
	get_str 1
	push STDOUT
	write
	int_to_str
	push STDOUT
	write

	get_str 0
	push STDOUT
	write`,
	expected: []string{
		"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55", "89", "144", "233", "377",
		"610", "987", "1597", "2584", "4181", "6765", "10946", "17711", "28657", "46368", "75025",
		"121393", "196418", "317811", "514229", "# of iterations 30",
	},
	additionalFiles: StdDefs,
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
	intToStrTest,
	stringPopStrTest,
	stringDupStrTest,
	stringSwapStrTest,
	stringInDupStrTest,
	stringInSwapStrTest,
	testFibTest,
}
