program:
	go run . "test.wm"

debug:
	go run . "test.wm" --debug

test:
	go test -v ./...