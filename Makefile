.PHONY: test
test:
	@cd ci && docker-compose up --force-recreate -d
	@go test -json -coverprofile=./coverage.txt -covermode=atomic ./...