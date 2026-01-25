package tests

// 1. Strcmp Test (Explicit equality check)
var StrcmpTest = ProgramTestCase{
	name: "native_strcmp",
	program: `@imp "stddefs.wm"
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
	program: `@imp "stddefs.wm"
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
	program: `@imp "stddefs.wm"
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
	program: `@imp "stddefs.wm"
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

var NativeStringTest = []ProgramTestCase{
	StrcmpTest,
	StrcpyTest,
	MemcpyTest,
	NativeIntToStrTest,
}
