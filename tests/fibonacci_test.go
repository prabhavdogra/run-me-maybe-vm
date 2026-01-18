package tests

var fib = ProgramTestCase{
	name: "fibonacci",
	program: `push 10
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
	expected: []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55"},
}
