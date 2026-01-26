package tests

var fibRecTest = ProgramTestCase{
	name: "fib_recursive",
	program: `
	@def N 10
	@imp "stddefs.wm"
	entrypoint main

	fib:
		; Arg n is on stack
		dup
		dup
		push 2
		cmpl             ; Stack: [n, n, 2, res]
		push 2
		swap
		nzjmp fib_base_case  ; Stack: [n, n, 2]
		
		; Cleanup comparison operands (n, 2)
		pop
		pop              ; [n]

		; Recursive case: fib(n-1) + fib(n-2)
		dup
		push 1
		sub
		call fib
		
		swap
		push 2
		sub
		call fib
		
		add
		ret

	fib_base_case:
		; Cleanup comparison operands (n, 2)
		pop
		pop              ; [n]
		ret

	main:
		push N
		call fib
		int_to_str
		push STDOUT
		write
		halt
	`,
	expected:        []string{"55"},
	additionalFiles: StdDefs,
}
