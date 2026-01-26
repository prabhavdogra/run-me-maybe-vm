package tests

var DuplicateEntrypointTest = ProgramTestCase{
	name: "duplicate_entrypoint",
	program: `
		entrypoint start
		entrypoint start
		start:
			push 1
			halt
	`,
	expectedError: "cannot define entrypoint more than once",
}

var entrypointTests = []ProgramTestCase{
	DuplicateEntrypointTest,
}
