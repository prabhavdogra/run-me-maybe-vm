package tests

// 1. Strcmp Test (Explicit equality check)
var StrcmpTest = ProgramTestCase{
	name: "native_strcmp",
	program: `@imp "stddefs.rmm"
		push_str "hello"
		push_str "hello"
		get_str 0
		get_str 1
		strcmp           ; Expect 1 (Equal)
		print
		
		push_str "world"
		get_str 0        ; "hello"
		get_str 2        ; "world"
		strcmp           ; Expect 0 (Not Equal)
		print
		halt`,
	expected: []string{
		"INT 1",
		"INT 0",
	},
	additionalFiles: StdDefs,
}

// 2. Strcpy Test (Copy string)
var StrcpyTest = ProgramTestCase{
	name: "native_strcpy",
	program: `@imp "stddefs.rmm"
		push 10
		malloc           ; Dst buffer
		push_str "copy"
		get_str 0        ; Src string
		strcpy           ; strcpy(dest, src) -> Pushes dest
		
		push STDOUT      ; Stack: [dest, STDOUT] -> Correct for write(fd, ptr) where fd is top
		
		write            ; Prints "copy"
		halt`,
	expected: []string{
		"copy",
	},
	additionalFiles: StdDefs,
}

// 3. Memcpy Test (Copy memory block)
var MemcpyTest = ProgramTestCase{
	name: "native_memcpy",
	program: `@imp "stddefs.rmm"
		push 10
		malloc           ; Dest buffer
		push_str "data"
		get_str 0        ; Src buffer
		push 4           ; Size
		memcpy           ; memcpy(dest, src, size) -> Pushes dest
		
		push STDOUT
		write            ; Prints "data" (since it copies 4 bytes)
		halt`,
	expected: []string{
		"data",
	},
	additionalFiles: StdDefs,
}

// 4. IntToStr Test (Convert int to string)
// Note: int_to_str native function (ID 99) behaves same as instruction?
// Native 99: Stack inputs [int]. Output [ptr].
var NativeIntToStrTest = ProgramTestCase{
	name: "native_int_to_str",
	program: `@imp "stddefs.rmm"
		push 12345
		native 99		 ; Converts to string pointer
		push STDOUT
		write            ; Prints "12345"
		halt`,
	expected: []string{
		"12345",
	},
	additionalFiles: StdDefs,
}

// 5. Strcat Test
var StrcatTest = ProgramTestCase{
	name: "native_strcat",
	program: `@imp "stddefs.rmm"
		push 20
		malloc           ; dest
		dup
		push_str "Hello "
		get_str 0        ; dest, dest, src "Hello "
		strcpy           ; dest, dest (strcpy returns dest)
		pop              ; dest
		
		dup              ; dest, dest
		push_str "World"
		get_str 1        ; dest, dest, src "World"
		strcat           ; dest, dest (strcat returns dest)
		pop              ; dest
		
		push STDOUT
		write            ; Prints "Hello World"
		halt`,
	expected: []string{
		"Hello World",
	},
	additionalFiles: StdDefs,
}

// 6. Native Strlen Test
var NativeStrlenTest = ProgramTestCase{
	name: "native_strlen",
	program: `@imp "stddefs.rmm"
		push_str "12345"
		get_str 0
		strlen           ; Native 94. Returns length (5)
		int_to_str       ; Back to string implementation check
		push STDOUT
		write
		halt`,
	expected: []string{
		"5",
	},
	additionalFiles: StdDefs,
}

// 7. Time Test
var TimeTest = ProgramTestCase{
	name: "native_time",
	program: `@imp "stddefs.rmm"
		push 0
		time             ; [0, T]
		cmpl             ; Checks T > 0? (Because a=T, b=0). Stack: [0, T, 1] (if T>0)
		
		; Branching based on result
		nzjmp success    ; Pops 1. Stack: [0, T].
		push_str "Time is zero or negative"
		get_str 0
		push STDOUT
		write
		halt

		success:
		push_str "Time ok"
		get_str 1
		push STDOUT
		write
		halt`,
	expected: []string{
		"Time ok",
	},
	additionalFiles: StdDefs,
}

// 8. Realloc Test
var ReallocTest = ProgramTestCase{
	name: "native_realloc",
	program: `@imp "stddefs.rmm"
		push 5
		malloc           ; ptr (size 5)
		dup              ; ptr, ptr
		push_str "Hi"
		get_str 0        ; ptr, ptr, src
		strcpy
		pop              ; ptr
		
		; Now realloc to 10
		push 10
		realloc          ; [ptr, 10] -> [new_ptr]
		
		dup              ; new_ptr, new_ptr
		push_str " there"
		get_str 1        ; new_ptr, new_ptr, src
		strcat           ; Append
		pop              ; new_ptr
		
		push STDOUT
		write
		halt`,
	expected: []string{
		"Hi there",
	},
	additionalFiles: StdDefs,
}

// 9. Assert Test
var AssertTest = ProgramTestCase{
	name: "native_assert",
	program: `@imp "stddefs.rmm"
		push 1
		assert           ; Should pass
		
		push 0
		assert           ; Should panic "assertion failed"
		halt`,
	expectedError:   "assertion failed",
	additionalFiles: StdDefs,
}

// 10. NULL Test
var NullTest = ProgramTestCase{
	name: "keyword_null",
	program: `@imp "stddefs.rmm"
		push NULL        ; Pushes LiteralNull
		print            ; Should print "NULL"
		halt`,
	expected: []string{
		"NULL",
	},
	additionalFiles: StdDefs,
}

var NativeStringTest = []ProgramTestCase{
	StrcmpTest,
	StrcpyTest,
	MemcpyTest,
	NativeIntToStrTest,
	StrcatTest,
	NativeStrlenTest,
	TimeTest,
	ReallocTest,
	AssertTest,
	NullTest,
}
