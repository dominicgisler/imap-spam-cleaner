.PHONY: default lint
.DEFAULT_GOAL := default

default:
	@docker build -f Dockerfile -t dominicgisler/imap-spam-cleaner .

lint:
	@golangci-lint run
