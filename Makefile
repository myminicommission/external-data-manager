start:
	go run . -- "@every 1m"

test:
	go test -v ./...

test-docker-build:
	@docker build -t mmc-external-datamanager:SNAPSHOT .
	@docker rmi mmc-external-datamanager:SNAPSHOT