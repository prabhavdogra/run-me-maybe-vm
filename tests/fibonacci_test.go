package tests

var fib = ProgramTestCase{
	name: "fibonacci",
	program: `@def N 10
	push N
	push 1
	push 1
	push 0

	indup 2
	inswap 1
	pop

	dup
	inswap 2
	pop

	indup 1
	indup 2
	add
	swap
	print

	indup 0
	push 1
	sub
	inswap 0

	nzjmp 4`,
	expected: []string{"INT 0", "INT 1", "INT 1", "INT 2", "INT 3", "INT 5", "INT 8", "INT 13", "INT 21", "INT 34", "INT 55"},
}
