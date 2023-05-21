start:
	go run ./cmd/api/

test:
	go test ./... --cover

test-integration:
	go test -tags=integration ./... --cover

test-with-race:
	go test -race -tags=integration ./... --cover

twr: test-with-race
ti: test-integration