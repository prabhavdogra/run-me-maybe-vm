package tests

// Common macro definitions for tests
var StdDefs = map[string]string{
	"stddefs.wm": `@def STDIN 0
	@def STDOUT 1
	@def STDERR 2

	@def RONLY 0
	@def WONLY 1
	@def RDWR 2
	@def CREAT 64
	@def EXCL 128

	@def open native 0
	@def write native 1
	@def read native 2
	@def close native 3
	@def malloc native 4
	@def realloc native 5
	@def free native 6
	@def time native 10
	@def exit native 60
	@def strcmp native 90
	@def strcpy native 91
	@def memcpy native 92
	@def strcat native 93
	@def strlen native 94
	@def int_to_str native 99
	@def assert native 100
	`,
	"std.wm": `
	print_newline:
		push '\n'
		ref
		push 1
		native 1
		pop
		ret

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
		ret`,
}
