

integration:
	docker-compose run --rm app make test

test:
	go test -v

jenkins: test
