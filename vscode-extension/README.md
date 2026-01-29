# rmm Assembly VSCode Extension

Syntax highlighting support for the rmm Assembly language (`.rmm` files).

## Features

- **Syntax Highlighting**:
  - Keywords (`push`, `pop`, `add`, `mov`, etc.)
  - Control Flow (`jmp`, `call`, `ret`)
  - Registers (`r0` - `r31`)
  - Directives (`@imp`, `@def`)
  - Literals (Integers, Floats, Strings, Characters)
  - Comments (`;`)

## Installation

1. Copy the `vscode-extension` folder to your local machine.
2. Run `code --install-extension vscode-extension` directory, or:
   - Package it using `vsce package`.
   - Install the resulting `.vsix` file using `code --install-extension rmm-assembly-0.0.1.vsix`.

## Development

- Open this folder in VSCode.
- Press `F5` to launch a new Extension Development Host window with the extension loaded.
