package tests

var RefDerefIntTest = ProgramTestCase{
	name: "ref_deref_int",
	program: `push 42
		ref
		deref
		print
		halt`,
	expected: []string{"INT 42"},
}

var RefDerefFloatTest = ProgramTestCase{
	name: "ref_deref_float",
	program: `push 3.14
		ref
		deref
		print
		halt`,
	expected: []string{"FLOAT 3.140000"},
}

var RefDerefCharTest = ProgramTestCase{
	name: "ref_deref_char",
	program: `push 'X'
		ref
		deref
		print
		halt`,
	expected: []string{"CHAR X"},
}

// Test multiple references / heap expansion
var HeapExpansionTest = ProgramTestCase{
	name: "heap_expansion",
	program: `push 1
		ref
		push 2
		ref
		deref
		print
		deref
		print
		halt`,
	expected: []string{"INT 2", "INT 1"},
}

// Test error handling: Dereferencing non-pointer
var DerefErrorTypeTest = ProgramTestCase{
	name: "deref_error_type",
	// Wait, is 10 a valid pointer? In our simplified model, yes, if heap has 11 elements.
	// But mostly it will be out of bounds or valid index.
	// Let's force an invalid type by using float, which cannot be a pointer index.
	program: `push 3.14
		deref
		halt`,
	expectedError: "deref requires a pointer",
}

var DerefErrorBoundsTest = ProgramTestCase{
	name: "deref_error_bounds",
	program: `push 9999
		deref
		halt`,
	expectedError: "deref requires a pointer",
}

// Complex case: Pointer arithmetic (if allowed? indices are ints, so yes)
// But we don't have explicit pointer types, just ints.
// "push 0; deref" should get first heap element.

var PointerArithmeticTest = ProgramTestCase{
	name: "pointer_arithmetic",
	program: `push 100
		ref         ; heap[0] = 100, stack = [0]
		push 200
		ref         ; heap[1] = 200, stack = [0, 1]
		pop         ; stack = [0]
		push 1
		add         ; stack = [1] -> we are manually calculating pointer to heap[1]
		deref
		print
		halt`,
	expected: []string{"INT 200"},
}

var PrintIntTest = ProgramTestCase{
	name: "print_int",
	program: `
		@imp "std.rmm"
		entrypoint main

		main:
			push 69420 
			call printint
			push_str "\n"
			get_str 0
			push 1 
			native 1
			pop`,
	additionalFiles: map[string]string{
		"std.rmm": `
		print_newline:
			push '\n'
			ref
			push 1
			native 1
			pop
			ret

		convert:
		dup              ; [n, n]
		push 0           ; [n, n, 0]
		cmpl             ; [n, n, 0, res]
		zjmp _not_neg    ; [n, n, 0] if res==0 (n>=0)
		; Fall-through: n < 0
		push '-'
		ref
		push 1 
		native 1
		pop
		swap             ; [n, 0] -> [0, n]
		push 0
		swap
		sub              ; [n_abs]
		jmp _after_neg
		_not_neg:
		; n >= 0
		_after_neg:
		dup              ; [n, n]
		push 9           ; [n, n, 9]
		cmpg             ; [n, n, 9, res]
		zjmp _lessthannine  ; [n, n, 9] if res==0 (n<=9)
		; Fall-through: n > 9
		; Fall-through: n > 9
		dup
		push 10
		div
		call convert
		jmp _after_recursive
		_lessthannine:
		; n <= 9
		_after_recursive:
		push 10
		mod
		push 48
		add
		ref
		push 1 
		native 1
		pop
		ret

		printint:
			call convert
			ret
		`,
	},
	expected: []string{"69420"},
}

var pointerTests = []ProgramTestCase{
	RefDerefIntTest,
	RefDerefFloatTest,
	RefDerefCharTest,
	HeapExpansionTest,
	DerefErrorTypeTest,
	DerefErrorBoundsTest,
	PointerArithmeticTest,
	PrintIntTest,
}
