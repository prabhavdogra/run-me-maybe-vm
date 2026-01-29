package tests

var NegativeTest = ProgramTestCase{
	name: "negative_numbers",
	program: `@imp "stddefs.rmm"
push_str "\n"
push -35
push -34
add
int_to_str
push STDOUT
write
get_str 0
push STDOUT
write`,
	expected: []string{
		"-69",
	},
	additionalFiles: StdDefs,
}
