start:
	go run ./cmd/api/

test:
	go test ./... --cover

test-integration:
	go test -tags=integration ./... --cover

test-with-race:
	go test -race -tags=integration ./... --cover

build:
	docker build -t lordrahl/subscriptions:latest .

docker-start: build
	docker-compose up

.PHONY: seed
seed: 
	go run ./cmd/seed

twr: test-with-race
ti: test-integration
ds: docker-start