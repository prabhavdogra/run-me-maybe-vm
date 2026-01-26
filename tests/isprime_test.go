package tests

var isPrime = ProgramTestCase{
	name: "is_prime",
	program: `push 2          ; n = 2
		jmp main_loop

		main_loop:
			; Check if n > 100
			indup 0         ; [n, n]
			push 100        ; [n, n, 100]
			indup 0			; [n, n, 100, n]
			push 100		; [n, n, 100, n, 100]
			cmpg            ; [n, n, 100, res]
			nzjmp end       ; [n, n, 100]
			pop             ; [n, n]
			pop             ; [n]

			; Prepare for primality check
			push 2          ; d = 2
			; [n, d]
			jmp check_loop

		check_loop:
			; Check if d == n (is prime)
			indup 0         ; [n, d, n]
			indup 1         ; [n, d, n, d]
			cmpe            ; [n, d, n, d, res]
			nzjmp is_prime  ; [n, d, n, d]
			pop             ; [n, d, n]
			pop             ; [n, d]

			; Check if n "%" d == 0 (not prime)
			indup 0         ; [n, d, n]
			indup 1         ; [n, d, n, d]
			mod             ; [n, d, rem]
			push 0          ; [n, d, rem, 0]
			cmpe            ; [n, d, rem, 0, res]
			nzjmp not_prime ; [n, d, rem, 0]
			pop             ; [n, d, rem]
			pop             ; [n, d]

			; d++
			push 1
			add             ; [n, d+1]
			jmp check_loop

		is_prime:
			; Entered with [n, d, n, d]
			pop
			pop             ; [n, d]
			pop             ; [n]
			indup 0         ; [n, n]
			print           ; print n. stack: [n]
			jmp next_iter

		not_prime:
			; Entered with [n, d, rem, 0]
			pop
			pop             ; [n, d]
			pop             ; [n]
			jmp next_iter

		next_iter:
			; Stack: [n]
			push 1
			add             ; n++
			jmp main_loop

		end:
			halt`,
	expected: []string{"INT 2", "INT 3", "INT 5", "INT 7", "INT 11", "INT 13", "INT 17", "INT 19", "INT 23", "INT 29", "INT 31", "INT 37", "INT 41", "INT 43", "INT 47", "INT 53", "INT 59", "INT 61", "INT 67", "INT 71", "INT 73", "INT 79", "INT 83", "INT 89", "INT 97"},
}
