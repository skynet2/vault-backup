.PHONY: test
test:
	@cd ci && docker-compose up --force-recreate -d
	@cd cmd && go .test