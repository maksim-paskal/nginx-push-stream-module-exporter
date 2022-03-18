tag=dev

test:
	./scripts/validate-license.sh
	go fmt ./cmd/
	go mod tidy
	go test -race ./cmd
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run -v
build:
	docker build --pull . -t paskalmaksim/nginx-push-stream-module-exporter:$(tag)
push:
	docker push -t paskalmaksim/nginx-push-stream-module-exporter:$(tag)
start:
	docker-compose down && docker-compose up
clean:
	docker-compose down
run:
	go run --race -v ./cmd/ -log.level=DEBUG -log.pretty -nginx.address=http://127.0.0.1:18102 $(args)
heap:
	go tool pprof -http=127.0.0.1:8080 http://localhost:8102/debug/pprof/heap
allocs:
	go tool pprof -http=127.0.0.1:8080 http://localhost:8102/debug/pprof/heap
git-prune-gc:
	curl -sSL https://get.paskal-dev.com/git-prune-gc | sh
build-binnary:
	make test
	@./scripts/build-all.sh
	ls -lah _dist
	go mod tidy
scan:
	trivy image \
	-ignore-unfixed --no-progress --severity HIGH,CRITICAL \
	paskalmaksim/nginx-push-stream-module-exporter:$(tag)