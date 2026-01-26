package tests

var StringEscapeTest = ProgramTestCase{
	name: "string_escapes",
	program: `@imp "stddefs.cmm"
		push_str "Hello\nWorld"
		get_str 0
		push STDOUT
		write
		push_str "\tTabbed"
		get_str 1
		push STDOUT
		write
		push_str "\"Quoted\""
		get_str 2
		push STDOUT
		write
		push_str "\\"
		get_str 3
		push STDOUT
		write`,
	expected: []string{
		"Hello",
		"World\tTabbed\"Quoted\"\\",
	},
	additionalFiles: StdDefs,
}
