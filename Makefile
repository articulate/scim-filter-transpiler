tools:
	@curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s v1.18.0

lint:
	@GOGC=30 ./bin/golangci-lint run ./

test:
	docker-compose run --rm app go test -v

build:
	@go build
