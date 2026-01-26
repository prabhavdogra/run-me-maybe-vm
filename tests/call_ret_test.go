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
		push 2
		cmpl
		nzjmp fib_base_case
		
		; Cleanup comparison operands (2, n)
		pop
		pop

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
		; Cleanup comparison operands (2, n)
		pop
		pop
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
