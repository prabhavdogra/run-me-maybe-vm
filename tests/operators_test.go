package tests

// Test addition with integers
var addTest = ProgramTestCase{
	name: "addition",
	program: `push 5
	push 3
	add
	print
	halt`,
	expected: []string{"INT 8"},
}

// Test subtraction with integers
var subTest = ProgramTestCase{
	name: "subtraction",
	program: `push 10
	push 3
	sub
	print
	halt`,
	expected: []string{"INT 7"},
}

// Test multiplication with integers
var mulTest = ProgramTestCase{
	name: "multiplication",
	program: `push 6
	push 7
	mul
	print
	halt`,
	expected: []string{"INT 42"},
}

// Test division with integers
var divTest = ProgramTestCase{
	name: "division",
	program: `push 20
	push 4
	div
	print
	halt`,
	expected: []string{"INT 5"},
}

// Test modulo with integers
var modTest = ProgramTestCase{
	name: "modulo",
	program: `push 17
	push 5
	mod
	print
	halt`,
	expected: []string{"INT 2"},
}

// Test comparison equal (true case)
var cmpeTestTrue = ProgramTestCase{
	name: "compare_equal_true",
	program: `push 5
	push 5
	cmpe
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 5", "INT 5"},
}

// Test comparison equal (false case)
var cmpeTestFalse = ProgramTestCase{
	name: "compare_equal_false",
	program: `push 5
	push 3
	cmpe
	print
	print
	print
	halt`,
	expected: []string{"INT 0", "INT 3", "INT 5"},
}

// Test comparison not equal (true case)
var cmpneTestTrue = ProgramTestCase{
	name: "compare_not_equal_true",
	program: `push 5
	push 3
	cmpne
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 3", "INT 5"},
}

// Test comparison not equal (false case)
var cmpneTestFalse = ProgramTestCase{
	name: "compare_not_equal_false",
	program: `push 5
	push 5
	cmpne
	print
	print
	print
	halt`,
	expected: []string{"INT 0", "INT 5", "INT 5"},
}

// Test comparison greater than (true case)
var cmpgTestTrue = ProgramTestCase{
	name: "compare_greater_true",
	program: `push 5
	push 10
	cmpg
	print
	halt`,
	expected: []string{"INT 0"},
}

// Test comparison greater than (false case)
var cmpgTestFalse = ProgramTestCase{
	name: "compare_greater_false",
	program: `push 5
	push 3
	cmpg
	print
	halt`,
	expected: []string{"INT 1"},
}

// Test comparison less than (true case)
var cmplTestTrue = ProgramTestCase{
	name: "compare_less_true",
	program: `push 10
	push 3
	cmpl
	print
	halt`,
	expected: []string{"INT 0"},
}

// Test comparison less than (false case)
var cmplTestFalse = ProgramTestCase{
	name: "compare_less_false",
	program: `push 3
	push 10
	cmpl
	print
	halt`,
	expected: []string{"INT 1"},
}

// Test comparison greater or equal (true - greater)
var cmpgeTestTrue = ProgramTestCase{
	name: "compare_greater_equal_true",
	program: `push 5
	push 10
	cmpge
	print
	print
	print
	halt`,
	expected: []string{"INT 0", "INT 10", "INT 5"},
}

// Test comparison greater or equal (true - equal)
var cmpgeTestEqual = ProgramTestCase{
	name: "compare_greater_equal_equal",
	program: `push 5
	push 5
	cmpge
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 5", "INT 5"},
}

// Test comparison greater or equal (false)
var cmpgeTestFalse = ProgramTestCase{
	name: "compare_greater_equal_false",
	program: `push 10
	push 3
	cmpge
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 3", "INT 10"},
}

// Test comparison less or equal (true - less)
var cmpleTestTrue = ProgramTestCase{
	name: "compare_less_equal_true",
	program: `push 10
	push 3
	cmple
	print
	print
	print
	halt`,
	expected: []string{"INT 0", "INT 3", "INT 10"},
}

// Test comparison less or equal (true - equal)
var cmpleTestEqual = ProgramTestCase{
	name: "compare_less_equal_equal",
	program: `push 5
	push 5
	cmple
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 5", "INT 5"},
}

// Test comparison less or equal (false)
var cmpleTestFalse = ProgramTestCase{
	name: "compare_less_equal_false",
	program: `push 3
	push 10
	cmple
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 10", "INT 3"},
}

// Test stack operations - dup
var dupTest = ProgramTestCase{
	name: "dup",
	program: `push 42
	dup
	print
	print
	halt`,
	expected: []string{"INT 42", "INT 42"},
}

// Test stack operations - swap
var swapTest = ProgramTestCase{
	name: "swap",
	program: `push 1
	push 2
	swap
	print
	print
	halt`,
	expected: []string{"INT 1", "INT 2"},
}

// Test indexed dup
var indupTest = ProgramTestCase{
	name: "indup",
	program: `push 10
	push 20
	push 30
	indup 2
	print
	print
	print
	print
	halt`,
	expected: []string{"INT 30", "INT 30", "INT 20", "INT 10"},
}

// Test indexed swap
var inswapTest = ProgramTestCase{
	name: "inswap",
	program: `push 10
	push 20
	push 30
	inswap 1
	print
	print
	print
	halt`,
	expected: []string{"INT 20", "INT 30", "INT 10"},
}

// Test jump
var jmpTest = ProgramTestCase{
	name: "jump",
	program: `push 1
	jmp 3
	push 2
	print
	halt`,
	expected: []string{"INT 1"},
}

// Test zero jump (should jump)
var zjmpTestJump = ProgramTestCase{
	name: "zjmp_jump",
	program: `push 1
	push 0
	zjmp 4
	push 2
	print
	halt`,
	expected: []string{"INT 1"},
}

// Test zero jump (should not jump)
var zjmpTestNoJump = ProgramTestCase{
	name: "zjmp_no_jump",
	program: `push 1
	push 5
	zjmp 4
	push 2
	print
	print
	halt`,
	expected: []string{"INT 2", "INT 1"},
}

// Test non-zero jump (should jump)
var nzjmpTestJump = ProgramTestCase{
	name: "nzjmp_jump",
	program: `push 1
	push 5
	nzjmp 4
	push 2
	print
	halt`,
	expected: []string{"INT 1"},
}

// Test non-zero jump (should not jump)
var nzjmpTestNoJump = ProgramTestCase{
	name: "nzjmp_no_jump",
	program: `push 1
	push 0
	nzjmp 4
	push 2
	print
	print
	halt`,
	expected: []string{"INT 2", "INT 1"},
}

// Test complex expression: (5 + 3) * 2
var complexExpr1 = ProgramTestCase{
	name: "complex_expression_1",
	program: `push 5
	push 3
	add
	push 2
	mul
	print
	halt`,
	expected: []string{"INT 16"},
}

// Test complex expression: 20 / (10 - 5)
var complexExpr2 = ProgramTestCase{
	name: "complex_expression_2",
	program: `push 30
	push 10
	sub
	push 5
	div
	print
	halt`,
	expected: []string{"INT 4"},
}

// Test float addition
var floatAddTest = ProgramTestCase{
	name: "float_addition",
	program: `push 3.5
	push 2.5
	add
	print
	halt`,
	expected: []string{"FLOAT 6.000000"},
}

// Test float subtraction
var floatSubTest = ProgramTestCase{
	name: "float_subtraction",
	program: `push 10.5
	push 3.5
	sub
	print
	halt`,
	expected: []string{"FLOAT 7.000000"},
}

// Test float multiplication
var floatMulTest = ProgramTestCase{
	name: "float_multiplication",
	program: `push 2.5
	push 4.0
	mul
	print
	halt`,
	expected: []string{"FLOAT 10.000000"},
}

// Test float modulo
var floatModTest = ProgramTestCase{
	name: "float_modulo",
	program: `push 10.5
	push 3.0
	mod
	print
	halt`,
	expected: []string{"FLOAT 1.500000"},
}

// Test float division
var floatDivTest = ProgramTestCase{
	name: "float_division",
	program: `push 10.0
	push 2.0
	div
	print
	halt`,
	expected: []string{"FLOAT 5.000000"},
}

// Test char push
var charPushTest = ProgramTestCase{
	name: "char_push",
	program: `push 'a'
	print
	halt`,
	expected: []string{"CHAR a"},
}

// Test multiple char push
var charMultipleTest = ProgramTestCase{
	name: "char_multiple",
	program: `push 'H'
	push 'i'
	print
	print
	halt`,
	expected: []string{"CHAR i", "CHAR H"},
}

// Test char with numbers
var charNumberTest = ProgramTestCase{
	name: "char_number",
	program: `push '5'
	push '9'
	print
	print
	halt`,
	expected: []string{"CHAR 9", "CHAR 5"},
}

// Test char equality (true)
var charEqualTestTrue = ProgramTestCase{
	name: "char_equal_true",
	program: `push 'x'
	push 'x'
	cmpe
	print
	print
	print
	halt`,
	expected: []string{"INT 1", "CHAR x", "CHAR x"},
}

// Test char equality (false)
var charEqualTestFalse = ProgramTestCase{
	name: "char_equal_false",
	program: `push 'a'
	push 'b'
	cmpe
	print
	print
	print
	halt`,
	expected: []string{"INT 0", "CHAR b", "CHAR a"},
}

// Test char dup
var charDupTest = ProgramTestCase{
	name: "char_dup",
	program: `push 'Z'
	dup
	print
	print
	halt`,
	expected: []string{"CHAR Z", "CHAR Z"},
}

// Test char swap
var charSwapTest = ProgramTestCase{
	name: "char_swap",
	program: `push 'A'
	push 'B'
	swap
	print
	print
	halt`,
	expected: []string{"CHAR A", "CHAR B"},
}

// Test comments
var commentTest = ProgramTestCase{
	name: "comments",
	program: `; This is a comment
	push 10 ; push ten
	push 5  ; push five
	; Another comment line
	add     ; add them
	print   ; print result
	halt    ; done`,
	expected: []string{"INT 15"},
}

// Test comment only lines
var commentOnlyTest = ProgramTestCase{
	name: "comment_only_lines",
	program: `; Comment at start
	; Another comment
	push 42
	; Mid-program comment
	print
	; Final comment
	halt`,
	expected: []string{"INT 42"},
}

var operatorsTest = []ProgramTestCase{
	// Arithmetic operations
	addTest,
	subTest,
	mulTest,
	divTest,
	floatModTest,
	modTest,

	// Comparison operations
	cmpeTestTrue,
	cmpeTestFalse,
	cmpneTestTrue,
	cmpneTestFalse,
	cmpgTestTrue,
	cmpgTestFalse,
	cmplTestTrue,
	cmplTestFalse,
	cmpgeTestTrue,
	cmpgeTestEqual,
	cmpgeTestFalse,
	cmpleTestTrue,
	cmpleTestEqual,
	cmpleTestFalse,

	// Stack operations
	dupTest,
	swapTest,
	indupTest,
	inswapTest,

	// Jump operations
	jmpTest,
	zjmpTestJump,
	zjmpTestNoJump,
	nzjmpTestJump,
	nzjmpTestNoJump,

	// Complex expressions
	complexExpr1,
	complexExpr2,

	// Float operations
	floatAddTest,
	floatSubTest,
	floatMulTest,
	floatDivTest,

	// Char operations
	charPushTest,
	charMultipleTest,
	charNumberTest,
	charEqualTestTrue,
	charEqualTestFalse,
	charDupTest,
	charSwapTest,

	// Comment tests
	commentTest,
	commentOnlyTest,
}
