package tests

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

type ProgramTestCase struct {
	name     string
	program  string
	expected []string
}

func TestPrograms(t *testing.T) {
	cases := []ProgramTestCase{
		fib,
		label,
		label2,
		floatPush,
	}

	literalLineRE := regexp.MustCompile(`^\s*(INT -?\d+|FLOAT -?\d+\.?\d*)\s*$`)

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// create a temp file for the program so tests don't rely on test.wm
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
