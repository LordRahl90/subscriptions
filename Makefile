start:
	go run ./cmd/

test:
	go test ./... --cover

test-integration:
	go test --tags=integration ./... --cover

test-with-race:
	go test -race ./... --cover

twr: test-with-race
ti: test-integration