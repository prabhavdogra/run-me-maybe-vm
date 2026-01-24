package tests

// Common macro definitions
// Defines macros for Syscall IDs and standard file descriptors for readability within tests.
var stddefs = map[string]string{
	"stddefs.tash": `@def STDOUT 1
	@def STDIN 0
	@def write native 1
	@def read native 2
	@def malloc native 3
	@def free native 4`,
}

// 1. Basic Write (Native 1)
// Verifies that 'write' correctly outputs a string from the stack to Stdout (ID 1).
var writeTest = ProgramTestCase{
	name: "syscall_write",
	program: `@imp "stddefs.tash"
		push "hello"  ; Push characters 'h', 'e', 'l', 'l', 'o'
		push STDOUT   ; Push File Descriptor (fd) 1 (Stdout)
		push 5        ; Push Length of string
		write         ; Execute write(fd, len), pops len, fd, then 5 chars
		halt`,
	expected:        []string{"hello"},
	additionalFiles: stddefs,
}

// 2. Write to Stderr (Native 1)
// Verifies that 'write' can target Stderr (ID 2), which is captured separately in tests.
var writeStderrTest = ProgramTestCase{
	name: "syscall_write_stderr",
	program: `@imp "stddefs.tash"
		push "error"
		push 2        ; Push File Descriptor (fd) 2 (Stderr)
		push 5
		write
		halt`,
	expectedStderr:  []string{"error"},
	additionalFiles: stddefs,
}

// 3. Buffer Lifecycle (Malloc -> Read -> Free)
// Verifies the complete lifecycle of dynamic memory:
// 1. Allocation (malloc) of a buffer on the Heap.
// 2. Writing to that buffer via 'read' syscall (using injected Input).
// 3. Deallocation (free) of the buffer.
var echoLifecycleTest = ProgramTestCase{
	name: "syscall_echo_lifecycle",
	program: `@imp "stddefs.tash"
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
	additionalFiles: stddefs,
}

// 4. Double Free Error (Native 4)
// Verifies that attempting to free the same pointer twice causes a crash (safety check).
var doubleFreeTest = ProgramTestCase{
	name: "syscall_error_double_free",
	program: `@imp "stddefs.tash"
		push 10
		malloc      ; Alloc ptr
		dup         ; Dup ptr
		free        ; Free once (Success)
		free        ; Free again (Expect Error: double free)
		halt`,
	expectedError:   "double free",
	additionalFiles: stddefs,
}

// 5. Read Overflow Error (Native 2)
// Verifies that 'read' prevents writing more bytes than explicitly allocated (Heap Safety).
var readOverflowTest = ProgramTestCase{
	name: "syscall_error_read_overflow",
	program: `@imp "stddefs.tash"
		push 2
		malloc      ; Allocate size 2
		
		push 5      ; Attempt to read 5 bytes
		push STDIN
		read        ; Expect Error: buffer overflow (5 > 2)
		halt`,
	input:           "hello",
	expectedError:   "buffer overflow",
	additionalFiles: stddefs,
}

// 6. Invalid Pointer (Native 4)
// Verifies that 'free' prevents freeing arbitrary integers that are not valid heap pointers.
var invalidFreeTest = ProgramTestCase{
	name: "syscall_error_invalid_free",
	program: `@imp "stddefs.tash"
		push 99999  ; Random integer
		free        ; Expect Error: invalid heap pointer
		halt`,
	expectedError:   "invalid heap pointer",
	additionalFiles: stddefs,
}

var syscallTests = []ProgramTestCase{
	writeTest,
	writeStderrTest,
	echoLifecycleTest,
	doubleFreeTest,
	readOverflowTest,
	invalidFreeTest,
}
