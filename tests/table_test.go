package tests

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

type ProgramTestCase struct {
	name            string
	program         string
	expected        []string
	additionalFiles map[string]string // Optional: for @imp tests, filename -> content
	input           string            // Optional: stdin for the program
	expectedError   string            // Optional: for error case tests
	expectedStderr  []string          // Optional: valid stderr output
}

func TestPrograms(t *testing.T) {
	cases := operatorsTest
	cases = append(cases, []ProgramTestCase{
		fib,
		label,
		label2,
		floatPush,
		isPrime,
	}...)
	cases = append(cases, PreprocessorTests...)
	cases = append(cases, stringTests...)
	cases = append(cases, syscallTests...)

	literalLineRE := regexp.MustCompile(`.+`)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// Catch panics if an error is expected
			defer func() {
				if r := recover(); r != nil {
					if tc.expectedError == "" {
						t.Fatalf("unexpected panic in test %s: %v", tc.name, r)
					}
					errStr := fmt.Sprint(r)
					if !strings.Contains(errStr, tc.expectedError) {
						t.Fatalf("expected error containing %q, got %q", tc.expectedError, errStr)
					}
				}
			}()

			// If test has additional files (for @imp), create a temp directory
			// Otherwise, use a single temp file (legacy behavior)
			if tc.additionalFiles != nil {
				tmpDir, err := os.MkdirTemp("", "vm-test-*")
				if err != nil {
					t.Fatalf("failed to create temp dir: %v", err)
				}
				defer os.RemoveAll(tmpDir)

				// Write all additional files
				for filename, content := range tc.additionalFiles {
					filePath := filepath.Join(tmpDir, filename)
					if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
						t.Fatalf("failed to write file %s: %v", filename, err)
					}
				}

				// Write main program with absolute path
				mainFilePath := filepath.Join(tmpDir, "main.wm")
				if err := os.WriteFile(mainFilePath, []byte(tc.program), 0644); err != nil {
					t.Fatalf("failed to write main file: %v", err)
				}

				// Get absolute path to VM module root (parent of tests directory)
				vmModuleRoot, err := filepath.Abs("..")
				if err != nil {
					t.Fatalf("failed to get VM module root: %v", err)
				}

				var stdout, stderr strings.Builder
				cmd := exec.Command("go", "run", ".", mainFilePath)
				cmd.Dir = vmModuleRoot
				cmd.Stdout = &stdout
				cmd.Stderr = &stderr
				if tc.input != "" {
					cmd.Stdin = strings.NewReader(tc.input)
				}
				err = cmd.Run()

				outStr := stdout.String()
				errStr := stderr.String()

				if tc.expectedError != "" {
					// We expect an error (non-zero exit or panic usually prints to stderr)
					// Check if the error message is in stderr or the error object
					combinedErr := fmt.Sprintf("%v\n%s", err, errStr)
					if !strings.Contains(combinedErr, tc.expectedError) {
						t.Fatalf("expected error containing %q, got err=%v, stderr=%q", tc.expectedError, err, errStr)
					}
					return
				} else if err != nil {
					// Unexpected error
					t.Fatalf("program failed unexpectedly: %v\nstderr: %s", err, errStr)
				}

				// Validate Stdout
				if len(tc.expected) > 0 {
					lines := strings.Split(strings.TrimSpace(outStr), "\n")
					var outputs []string
					for _, ln := range lines {
						if literalLineRE.MatchString(ln) {
							outputs = append(outputs, strings.TrimSpace(ln))
						}
					}
					if len(outputs) < len(tc.expected) {
						t.Fatalf("expected at least %d stdout lines, got %d; full stdout:\n%s", len(tc.expected), len(outputs), outStr)
					}
					for i := range tc.expected {
						if outputs[i] != tc.expected[i] {
							t.Fatalf("stdout mismatch at index %d: expected %s, got %s\nfull stdout:\n%s", i, tc.expected[i], outputs[i], outStr)
						}
					}
				} else if len(tc.expectedStderr) == 0 && outStr != "" {
					// Strict verification: if no expected stdout is specified, enforce that no stdout is produced
					t.Fatalf("unexpected stdout produced: %q", outStr)
				}

				// Validate Stderr (if expected)
				if len(tc.expectedStderr) > 0 {
					lines := strings.Split(strings.TrimSpace(errStr), "\n")
					var outputs []string
					for _, ln := range lines {
						if literalLineRE.MatchString(ln) {
							outputs = append(outputs, strings.TrimSpace(ln))
						}
					}
					if len(outputs) < len(tc.expectedStderr) {
						t.Fatalf("expected at least %d stderr lines, got %d; full stderr:\n%s", len(tc.expectedStderr), len(outputs), errStr)
					}
					for i := range tc.expectedStderr {
						if outputs[i] != tc.expectedStderr[i] {
							t.Fatalf("stderr mismatch at index %d: expected %s, got %s\nfull stderr:\n%s", i, tc.expectedStderr[i], outputs[i], errStr)
						}
					}
				}
				return
			}

			// Legacy single-file mode
			tmp, err := os.CreateTemp("", "prog-*.wm")
			if err != nil {
				t.Fatalf("failed to create temp file: %v", err)
			}
			defer os.Remove(tmp.Name())

			if _, err := tmp.WriteString(tc.program); err != nil {
				t.Fatalf("failed to write program to temp file: %v", err)
			}
			if err := tmp.Close(); err != nil {
				t.Fatalf("failed to close temp file: %v", err)
			}

			cmd := exec.Command("go", "run", "..", tmp.Name())
			outBytes, err := cmd.CombinedOutput()

			if err != nil {
				// If we expect an error, this is fine, the recover() block will handle it
				panic(string(outBytes))
			}

			out := strings.TrimSpace(string(outBytes))
			if out == "" {
				t.Fatal("program produced no output")
			}

			// collect literal lines (PRINT outputs)
			lines := strings.Split(out, "\n")
			var outputs []string
			for _, ln := range lines {
				if literalLineRE.MatchString(ln) {
					outputs = append(outputs, strings.TrimSpace(ln))
				}
			}

			if len(outputs) < len(tc.expected) {
				t.Fatalf("expected at least %d literal outputs, got %d; full output:\n%s", len(tc.expected), len(outputs), out)
			}
			for i := range tc.expected {
				if outputs[i] != tc.expected[i] {
					t.Fatalf("mismatch at index %d: expected %s, got %s\nfull output:\n%s", i, tc.expected[i], outputs[i], out)
				}
			}
		})
	}
}
