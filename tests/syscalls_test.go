package tests

import "os"

// Common macro definitions are now in common.go (StdDefs)

// 1. Basic Write (Native 1)
// Verifies that 'write' correctly outputs a string from the stack to Stdout (ID 1).
var writeTest = ProgramTestCase{
	name: "syscall_write",
	program: `@imp "stddefs.rmm"
		push_str "hello"
		get_str 0     ; Pushes ptr to "hello"
		push STDOUT   ; Push File Descriptor (fd) 1 (Stdout)
		write         ; Execute write(fd, ptr)`,
	expected:        []string{"hello"},
	additionalFiles: StdDefs,
}

// 2. Write to Stderr (Native 1)
// Verifies that 'write' can target Stderr (ID 2), which is captured separately in tests.
var writeStderrTest = ProgramTestCase{
	name: "syscall_write_stderr",
	program: `@imp "stddefs.rmm"
		push_str "error"
		get_str 0
		push 2        ; Push File Descriptor (fd) 2 (Stderr)
		write
		halt`,
	expectedStderr:  []string{"error"},
	additionalFiles: StdDefs,
}

// 3. Buffer Lifecycle (Malloc -> Read -> Free)
// Verifies the complete lifecycle of dynamic memory:
// 1. Allocation (malloc) of a buffer on the Heap.
// 2. Writing to that buffer via 'read' syscall (using injected Input).
// 3. Deallocation (free) of the buffer.
var echoLifecycleTest = ProgramTestCase{
	name: "syscall_echo_lifecycle",
	program: `@imp "stddefs.rmm"
		; Allocate 5 bytes on Heap
		push 5
		malloc      ; Pops size 5, Pushes Heap Pointer (ptr)
		dup         ; Duplicate ptr so we can use it for read AND free later
		
		; Read 5 bytes from Stdin into the allocated buffer
		push 5      ; Buffer size to read
		push STDIN  ; File Descriptor (fd) 0 (Stdin)
		read        ; read(fd, len, ptr). Pops FD, Len, Ptr. Reads 5 bytes from input.
		
		; Free the pointer (original copy on stack)
		free        ; free(ptr). Pops ptr. Deletes from heap.
		halt`,
	input:           "hello",    // Injected Stdout for 'read'
	expected:        []string{}, // Logic verification (no crash) is sufficient
	additionalFiles: StdDefs,
}

// 4. Double Free Error (Native 4)
// Verifies that attempting to free the same pointer twice causes a crash (safety check).
var doubleFreeTest = ProgramTestCase{
	name: "syscall_error_double_free",
	program: `@imp "stddefs.rmm"
		push 10
		malloc      ; Alloc ptr
		dup         ; Dup ptr
		free        ; Free once (Success)
		free        ; Free again (Expect Error: double free)
		halt`,
	expectedError:   "double free",
	additionalFiles: StdDefs,
}

// 5. Read Overflow Error (Native 2)
// Verifies that 'read' prevents writing more bytes than explicitly allocated (Heap Safety).
var readOverflowTest = ProgramTestCase{
	name: "syscall_error_read_overflow",
	program: `@imp "stddefs.rmm"
		push 2
		malloc      ; Allocate size 2
		
		push 5      ; Attempt to read 5 bytes
		push STDIN
		read        ; Expect Error: buffer overflow (5 > 2)
		halt`,
	input:           "hello",
	expectedError:   "buffer overflow",
	additionalFiles: StdDefs,
}

// 6. Invalid Pointer (Native 4)
// Verifies that 'free' prevents freeing arbitrary integers that are not valid heap pointers.
var invalidFreeTest = ProgramTestCase{
	name: "syscall_error_invalid_free",
	program: `@imp "stddefs.rmm"
		push 99999  ; Random integer
		free        ; Expect Error: invalid heap pointer
		halt`,
	expectedError:   "free pointer must be pointer",
	additionalFiles: StdDefs,
}

// 7. Explicit Exit (Native 5)
// Verifies that 'exit(code)' terminates the program with the specified status code.
var exitTest = ProgramTestCase{
	name: "syscall_exit_explicit",
	program: `@imp "stddefs.rmm"
		push 69     ; Exit Code
		exit
		`,
	expectedError:   "exit status 69",
	additionalFiles: StdDefs,
}

// 8. File Operations (Open -> Read -> Close)
// Verifies that we can open a file, read from it, and close it.
var fileOpsTest = ProgramTestCase{
	name: "syscall_file_ops",
	program: `@imp "stddefs.rmm"
		; 1. Read filename "A" from Stdin into Heap
		push 1          ; Size 1
		malloc          ; -> ptr
		dup             ; Keep ptr for open
		dup             ; Keep ptr for read
		push 1          ; len
		push STDIN
		                ; My inswap analysis: [ptr, 1, 0]. Swap 0(Top) and ptr(2). -> [0, 1, ptr].
		                ; Top: ptr. Next: 1. Next: 0. Correct.
		read            ; read(0, 1, ptr) -> reads "A"
		
		; 2. Open file "A" (Create)
		; Stack: ptr
		dup             ; Keep ptr for 2nd open (Read)
		push 1          ; len of filename "A"
		push CREAT      ; 64
		push RDWR       ; 2
		add             ; flags (66)
		open            ; open(flags, 1, ptr) -> pushes FD (should be 3)
		pop             ; Discard FD (assume 3 for simplicity)

		; 3. Write "Hello" to FD 3
		push_str "Hello"
		get_str 0       ; ptr to "Hello" (Index 0)
		push 3          ; FD
		write           ; write(3, ptr)
		pop

		; 4. Close FD 3
		push 3
		close

		; 5. Open file "A" (Read)
		; Stack: ptr
		push 1          ; len
		push RONLY      ; flags (0)
		open            ; -> FD (should be 3 again)
		pop             ; Discard FD

		; 6. Read from FD 3
		push 5          ; Size to read "Hello"
		malloc          ; -> bufPtr
		push 5          ; len
		push 3          ; FD
		read            ; read(3, 5, bufPtr)
		
		; 7. Close FD 3
		push 3
		close
		halt`,
	input:           "A",
	additionalFiles: StdDefs,
	expected:        []string{},
	cleanup: func() {
		os.Remove("../A")
	},
}

var syscallTests = []ProgramTestCase{
	writeTest,
	writeStderrTest,
	echoLifecycleTest,
	doubleFreeTest,
	readOverflowTest,
	invalidFreeTest,
	exitTest,
	fileOpsTest,
}
