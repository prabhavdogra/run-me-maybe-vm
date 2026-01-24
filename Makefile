program:
	go run . "program.wm"

program2:
	go run . "test2.wm"

program3:
	go run . "t_fib.wm"

debug:
	go run . "test.wm" --debug

test:
	go test -v ./... -timeout=30s | grep -i fail || true

run_test:
	go test -v ./... -timeout=30s