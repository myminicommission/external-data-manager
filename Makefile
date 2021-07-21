start:
	go run . -- "@every 1m"

start-docker: build-docker
	@docker run mmc-external-datamanager:SNAPSHOT

build-docker:
	@docker build -t mmc-external-datamanager:SNAPSHOT .

remove-docker:
	@docker rmi mmc-external-datamanager:SNAPSHOT

test:
	go test -v ./...

test-docker-build: build-docker remove-docker