program:
	go run . "test.wm"

program2:
	go run . "test2.wm"

program3:
	go run . "t_fib.wm"

debug:
	go run . "test.wm" --debug

test:
	go test -v ./... | grep -i fail || true

run_test:
	go test -v ./...