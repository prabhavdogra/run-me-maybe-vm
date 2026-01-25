program:
	go run . "program.wm"

program2:
	go run . "program2.wm"

program3:
	go run . "t_fib.wm"

debug:
	go run . "program.wm" --debug

test:
	go test -v ./... -timeout=30s | grep -i fail || true

run_test:
	go test -v ./... -timeout=30s