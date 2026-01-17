package main

import (
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestPrograms(t *testing.T) {
	cases := []struct {
		name     string
		program  string
		expected []string
	}{
		{
			name: "fibonacci",
			program: `
				push 10
				push 1
				push 1
				push 0

				indup 2
				inswap 1
				pop

				dup
				inswap 2
				pop

				indup 1
				indup 2
				add
				swap
				print

				push 1
				indup 0
				sub
				inswap 0
				nzjmp 4
			`,
			expected: []string{"0", "1", "1", "2", "3", "5", "8", "13", "21", "34", "55"},
		},
	}

	intLineRE := regexp.MustCompile(`^\s*-?\d+\s*$`)

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

			// collect integer-only lines (PRINT outputs)
			lines := strings.Split(out, "\n")
			var numbers []string
			for _, ln := range lines {
				if intLineRE.MatchString(ln) {
					numbers = append(numbers, strings.TrimSpace(ln))
				}
			}

			if len(numbers) < len(tc.expected) {
				t.Fatalf("expected at least %d integer outputs, got %d; full output:\n%s", len(tc.expected), len(numbers), out)
			}
			for i := range tc.expected {
				if numbers[i] != tc.expected[i] {
					t.Fatalf("mismatch at index %d: expected %s, got %s\nfull output:\n%s", i, tc.expected[i], numbers[i], out)
				}
			}
		})
	}
}
