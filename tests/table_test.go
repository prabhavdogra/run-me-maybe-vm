package tests

import (
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

	literalLineRE := regexp.MustCompile(`^\s*(INT -?\d+|FLOAT -?\d+\.?\d*|CHAR .)\s*$`)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
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

				// Run from module root, passing absolute path to main file
				cmd := exec.Command("go", "run", ".", mainFilePath)
				cmd.Dir = vmModuleRoot
				outBytes, err := cmd.CombinedOutput()
				if err != nil {
					t.Fatalf("program failed: %v\noutput:\n%s", err, string(outBytes))
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
				t.Fatalf("program failed: %v\noutput:\n%s", err, string(outBytes))
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
