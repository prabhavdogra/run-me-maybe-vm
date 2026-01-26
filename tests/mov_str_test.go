package tests

var movStrTests = []ProgramTestCase{
	{
		name: "mov_str_char",
		program: `
			push 'H'
			mov_str
			get_str 0
			deref
			print
		`,
		expected: []string{"CHAR H"},
	},
	{
		name: "mov_str_ptr",
		program: `
			push_str "Hello"
			get_str 0    ; Pointer to "Hello"
			mov_str      ; Push pointer to string stack (idx 1)
			get_str 1
			deref
			print
		`,
		expected: []string{"CHAR H"},
	},
	{
		name: "mov_str_invalid",
		program: `
			push 12.3
			mov_str
		`,
		expectedError: "mov_str requires char or int (pointer)",
	},
}
