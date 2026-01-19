package tests

var floatPush = ProgramTestCase{
	name: "float",
	program: `push 3.14
	push 3.15
	print
	print`,
	expected: []string{"FLOAT 3.150000", "FLOAT 3.140000"},
}
