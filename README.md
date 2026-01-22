# Go Stack Virtual Machine

A lightweight, stack-based virtual machine written in Go, featuring a custom assembly-like language, a lexer with preprocessor support, and a parser.

## Features

- **Stack-based Execution**: All operations occur on a central stack.
- **Preprocessor**:
  - `@imp "file.wm"`: Resolve and include external files.
  - `@def NAME VALUE`: Simple macro substitution.
- **Support for Multiple Literals**: Integers, Floats, and Characters.
- **Flexible Instruction Set**: Arithmetic, stack manipulation, control flow, and comparison.
- **Label Support**: Use labels for jump targets instead of hardcoded instruction pointers.

> **Note on Strings**: String literals (e.g., `"hello"`) are pushed onto the stack as a sequence of individual characters. Use `write` to output them.

## Quick Start

### Prerequisites
- Go 1.25 or higher

### Running a Program
```bash
go run . path/to/source.wm
```

### Running in Debug Mode
Debug mode prints the lexed tokens, parsed instruction list, and the final state of the stack.
```bash
go run . path/to/source.wm --debug
```

### Running Tests
```bash
go test -v ./...
```

## Instruction Set

| Instruction | Description |
| :--- | :--- |
| `push <val>` | Push a literal value (int, float, char) onto the stack. |
| `pop` | Remove the top value from the stack. |
| `write <fd> <len>` | Pop `len` characters and write to file descriptor (1=stdout, 2=stderr). |
| `dup` | Duplicate the top stack value. |
| `indup <idx>` | Duplicate the value at the given stack index to the top. |
| `swap` | Swap the top two stack values. |
| `inswap <idx>` | Swap the top value with the value at the given index. |
| `add`, `sub`, `mul`, `div`, `mod` | Arithmetic operations (pops 2, pushes result). |
| `cmpe`, `cmpne`, `cmpg`, `cmpl`, `cmpge`, `cmle` | Comparison operations (pops 2, pushes bool int 1/0). |
| `jmp <label/ptr>` | Unconditional jump. |
| `zjmp <label/ptr>` | Jump if top of stack is 0 (pops condition). |
| `nzjmp <label/ptr>` | Jump if top of stack is NOT 0 (pops condition). |
| `print` | Pop and print the top value (int, char, float). |
| `halt` | Stop execution. |
| `noop` | No operation. |

## Preprocessor Directives

- **Imports**: Resolve paths relative to the current file.
  ```assembly
  @imp "stdlib.wm"
  ```
- **Macros**: Simple text substitution for constants or shorthand.
  ```assembly
  @def MAX_SIZE 100
  push MAX_SIZE
  ```

## Project Structure

- `main.go`: Entry point, orchestrates lexing, parsing, and execution.
- `internal/lexer/`: Lexical analysis and macro expansion.
- `internal/parser/`: Token processing and label resolution.
- `internal/token/`: Token types and helper functions.
- `tests/`: Integration and unit tests.
- `instruction.go`: VM execution loop and logic.
- `literal.go`: Data type definitions and literal arithmetic.
