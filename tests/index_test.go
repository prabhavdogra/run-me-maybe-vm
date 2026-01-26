package tests

var indexTests = []ProgramTestCase{
	{
		name: "index_modify_char",
		program: `
			push_str "Hello"
			get_str 0    ; ptr
			push 0       ; index
			index 'J'    ; Set "Hello"[0] = 'J' -> "Jello"
			deref        ; Read char at ptr (index 0)
			print        ; Should be 'J' (74)
		`,
		expected: []string{"CHAR J"},
	},
	{
		name: "index_modify_middle",
		program: `
			push_str "Bit"
			get_str 0    ; ptr
			push 1       ; index
			index 'a'    ; Set "Bit"[1] = 'a' -> "Bat"
			push 1
			add          ; ptr + 1
			deref
			print        ; Should be 'a' (97)
		`,
		expected: []string{"CHAR a"},
	},
	{
		name: "index_error_negative",
		program: `
			push_str "Hello"
			get_str 0
			push -1
			index 'X'
		`,
		expectedError: "index cannot be less than 0",
	},
	{
		name: "index_error_oob",
		program: `
			push_str "Hi"
			get_str 0
			push 100
			index 'X'
		`,
		expectedError: "segmentation fault: index out of bounds",
	},
	{
		name: "index_stack_usage",
		program: `
			push_str "Hello"
			get_str 0    ; ptr
			push 0       ; index
			push 'J'     ; val
			index        ; Set "Hello"[0] = 'J' using stack args
			deref        ; Read char at ptr (index 0)
			print        ; Should be 'J' (74)
		`,
		expected: []string{"CHAR J"},
	},
}
