package tests

// Common macro definitions for tests
var StdDefs = map[string]string{
	"stddefs.cmm": `@def STDIN 0
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
	@def pow native 8
	@def close native 3
	@def malloc native 4
	@def realloc native 5
	@def free native 6
	@def scanf native 7
	@def time native 10
	@def exit native 60
	@def strcmp native 90
	@def strcpy native 91
	@def memcpy native 92
	@def strcat native 93
	@def strlen native 94
	@def float_to_str native 98
	@def int_to_str native 99
	@def assert native 100
	`,
	"std.cmm": `print_newline:
    push '\n'
    ref
    push 1
    native 1
    pop
    ret

convert:
    dup
    push 0
    cmpl
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

printfloat:
    dup
    ftoi 
    dup
    inswap 0
    swap
    itof
    sub_f
    mov r0 top
    pop
    call printint
    push '.'
    ref
    push 1
    native 1
    push r0
    push 8
    push 10
    native 8
    itof
    mul_f
    ftoi
    push 0
    cmpg
    zjmp _notneg
    pop
    dup
    push 2
    mul
    sub 
    jmp _endneg
    _notneg:
    pop
    _endneg:

    call printint
    pop
    ret

_strcmp_not_equal:
    push 0
    ret

stringcmp:
    _loop:
    push r0
    deref
    push r1
    deref
    inswap 0
    dup
    push '\0'
    cmpe
    nzjmp _c0_is_null
    swap
    cmpe
    zjmp _strcmp_not_equal
    push r0
    push 1
    add
    mov r0 top
    push r1
    push 1
    add
    mov r1 top
    jmp _loop

_c0_is_null:
    pop
    push '\0'
    cmpe 
    nzjmp _equal
    push 0
    ret

_equal:
    push 1
    ret

stringlen:
    push '\0'
    mov r0 top
    push 0
    swap
    _strlen_loop:
        dup
        deref
        push r0
        cmpe
        nzjmp _strlen_end
        push 1
        add
        swap
        push 1
        add
        swap
        jmp _strlen_loop
    _strlen_end:
        pop
        ret
`,
}
