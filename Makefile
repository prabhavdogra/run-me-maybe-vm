program:
	go run . "program.rmm"

program2:
	go run . "program2.rmm"

program3:
	go run . "t_fib.rmm"

debug:
	go run . "program.rmm" --debug

test:
	go test -v ./... -timeout=30s | grep -i fail || true

run_test:
	go test -v ./... -timeout=30s

tools:
	cd vscode-extension && vsce package

install_tools:
	cd vscode-extension && code --install-extension runmemaybeasm-0.0.1.vsix

uninstall_tools:
	code --uninstall-extension undefined_publisher.runmemaybeasm