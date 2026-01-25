package tests

var itofTest = ProgramTestCase{
	name: "itof_test",
	program: `push 10
		itof
		push 2.5
		add
		print
		halt`,
	expected: []string{"FLOAT 12.500000"},
}

var itofNegTest = ProgramTestCase{
	name: "itof_neg_test",
	program: `push -5
		itof
		print
		halt`,
	expected: []string{"FLOAT -5.000000"},
}

var ftoiTest = ProgramTestCase{
	name: "ftoi_test",
	program: `push 12.34
		ftoi
		print
		halt`,
	expected: []string{"INT 12"},
}

var ftoiNegTest = ProgramTestCase{
	name: "ftoi_neg_test",
	program: `push -9.99
		ftoi
		print
		halt`,
	expected: []string{"INT -9"},
}

var ftoiRoundTest = ProgramTestCase{
	name: "ftoi_round_test",
	program: `push 0.99
		ftoi
		print
		halt`,
	expected: []string{"INT 0"},
}

// Test error handling
var itofErrorTest = ProgramTestCase{
	name: "itof_error_test",
	program: `push 1.5
		itof
		halt`,
	expectedError: "itof requires an integer",
}

var ftoiErrorTest = ProgramTestCase{
	name: "ftoi_error_test",
	program: `push 5
		ftoi
		halt`,
	expectedError: "ftoi requires a float",
}

var castTests = []ProgramTestCase{
	itofTest,
	itofNegTest,
	ftoiTest,
	ftoiNegTest,
	ftoiRoundTest,
	itofErrorTest,
	ftoiErrorTest,
}
