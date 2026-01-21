package tests

// Test basic string push and print
var stringPushTest = ProgramTestCase{
	name: "string_push",
	program: `push "Hello, World!"
	print
	halt`,
	expected: []string{"Hello, World!"},
}

// Test string with escape sequences
var stringEscapeTest = ProgramTestCase{
	name: "string_escape",
	program: `push "Line 1\nLine 2\tTabbed"
	print
	halt`,
	expected: []string{"Line 1\nLine 2\tTabbed"}, // Note: internal/token/helper.go handles parsing \n to actual newline
}

// Note: The above expected string might be split into multiple lines by the test runner if it splits by \n.
// The test runner does: lines := strings.Split(out, "\n")
// So "Line 1\nLine 2" becomes two outputs: "Line 1", "Line 2".
// We need to adjust the expected output for multi-line strings.

var stringEscapeTestMultiLine = ProgramTestCase{
	name: "string_escape_multiline",
	program: `push "Line 1\nLine 2"
	print
	halt`,
	expected: []string{"Line 1", "Line 2"},
}

// Test string dup
var stringDupTest = ProgramTestCase{
	name: "string_dup",
	program: `push "copy me"
	dup
	print
	print
	halt`,
	expected: []string{"copy me", "copy me"},
}

// Test string swap
var stringSwapTest = ProgramTestCase{
	name: "string_swap",
	program: `push "first"
	push "second"
	swap
	print
	print
	halt`,
	expected: []string{"first", "second"},
}

// Test string equality (true)
var stringEqualTestTrue = ProgramTestCase{
	name: "string_equal_true",
	program: `push "test"
	push "test"
	cmpe
	print
	halt`,
	expected: []string{"INT 1"},
}

var stringTests = []ProgramTestCase{
	stringPushTest,
	stringEscapeTestMultiLine,
	stringDupTest,
	stringSwapTest,
	stringEqualTestTrue,
}
