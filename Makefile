program:
	go run . "program.cmm"

program2:
	go run . "program2.cmm"

program3:
	go run . "t_fib.cmm"

debug:
	go run . "program.cmm" --debug

test:
	go test -v ./... -timeout=30s | grep -i fail || true

run_test:
	go test -v ./... -timeout=30s

tools:
	cd vscode-extension && vsce package

install_tools:
	cd vscode-extension && code --install-extension compilememaybeasm-0.0.2.vsix

uninstall_tools:
	code --uninstall-extension undefined_publisher.compilememaybeasm