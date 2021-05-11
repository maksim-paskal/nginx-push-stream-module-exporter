test:
	./scripts/validate-license.sh
	go fmt ./cmd/
	go mod tidy
	go test -race ./cmd
	golangci-lint run -v
build:
	docker-compose build
start:
	docker-compose down && docker-compose up
clean:
	docker-compose down
run:
	go run --race -v ./cmd/ -log.level=DEBUG -log.pretty -nginx.address=http://127.0.0.1:18102 $(args)