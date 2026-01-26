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
	expectedError: "deref requires a pointer (int)",
}

var DerefErrorBoundsTest = ProgramTestCase{
	name: "deref_error_bounds",
	program: `push 9999
		deref
		halt`,
	expectedError: "segmentation fault: invalid pointer",
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
		@imp "std.wm"
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
		"std.wm": `
		convert:
			push 0
			cmpg
			zjmp _not_neg
			push '-'
			ref
			push 1 
			native 1
			pop
			swap
			push 0
			swap
			sub
			swap
		_not_neg:
			pop
			push 9
			cmpl
			zjmp _lessthannine
			pop
			dup
			push 10
			div
			call convert 
		_lessthannine:
			pop
			push 10
			mod
			push 48
			add
			ref
			push 1 
			native 1
			ret

		printint:
			call convert
			pop
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
