package tests

var StdcmmContent = `print_newline:
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
    push 1
    sub
    jmp _strcmp_end

stringcmp:
    push 1
    _strcmp_loop:
        push r0
        deref
        swap
        push 1
        add
        mov r0 top
        pop
        push r1
        deref
        swap
        push 1
        add
        mov r1 top
        pop
        cmpe
        swap
        pop
        swap
        pop
    zjmp _strcmp_not_equal
    push '\0'
    push r0
    deref
    swap
    pop
    cmpe
    swap
    pop
    swap
    pop
    zjmp _strcmp_loop
    _strcmp_end:
    ret

stringlen:
    push '\0'
    mov r0 top
    pop
    push 0
    swap
    _strlen_loop:
        deref
        push r0
        cmpe
        swap
        pop
        swap
        pop
        ;native 90 ;strcmp
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
`

var StdLibFiles map[string]string

func GetStdLibFiles() map[string]string {
	files := make(map[string]string)
	for k, v := range StdDefs {
		files[k] = v
	}
	return files
}

var stdLibTests = []ProgramTestCase{
	{
		name: "std_stringcmp_equal",
		program: `
				entrypoint main
				@imp "std.rmm"
				main:
				push_str "Hello"
				push_str "Hello"
				get_str 0  ; ptr1
				get_str 1  ; ptr2
				mov r0 top ; ptr1 -> r0
				mov r1 top ; ptr2 -> r1
				call stringcmp
				print ; Should print 1 (equal) or 0 (not equal)
				halt`,
		input:           "",
		expected:        []string{"INT 1"},
		additionalFiles: StdDefs,
	},
	{
		name: "std_stringcmp_not_equal",
		program: `
				entrypoint main
				@imp "std.rmm"
				main:
				push_str "Hello"
				push_str "World"
				get_str 0
				get_str 1
				mov r0 top
				mov r1 top
				call stringcmp
				print
				halt`,
		expected:        []string{"INT 0"},
		additionalFiles: StdDefs,
	},
	{
		name: "std_stringlen",
		program: `
				entrypoint main
				@imp "std.rmm"
				main:
				push_str "Hello"
				get_str 0
				call stringlen
				print
				halt`,
		expected:        []string{"INT 5"},
		additionalFiles: StdDefs,
	},
}
