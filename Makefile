test:
	./scripts/validate-license.sh
	go fmt ./cmd/
	go mod tidy
	golangci-lint run --allow-parallel-runners -v --enable-all --disable gochecknoglobals,funlen,gosec --fix
build:
	docker-compose build
start:
	docker-compose down && docker-compose up
clean:
	docker-compose down
run:
	./scripts/test.sh