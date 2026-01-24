package tests

// Common macro definitions for tests
var StdDefs = map[string]string{
	"stddefs.wm": `@def STDOUT 1
	@def STDIN 0
	@def open native 0
	@def write native 1
	@def read native 2
	@def close native 3
	@def free native 4
	@def malloc native 5
	@def exit native 6

	@def RONLY 0
	@def WONLY 1
	@def RDWR 2
	@def CREAT 64
	@def EXCL 128
	`,
}
