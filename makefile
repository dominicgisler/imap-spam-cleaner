.PHONY: default lint publish --tag
.DEFAULT_GOAL := default

default: --tag
	@docker build --platform linux/amd64,linux/arm64 -f Dockerfile -t dominicgisler/imap-spam-cleaner:$(TAG) .

lint:
	@golangci-lint run

publish: --tag
	@docker push dominicgisler/imap-spam-cleaner:$(TAG)

--tag:
	@if [ "$(TAG)" = "" ]; then echo "TAG not set" && exit 1; else exit 0; fi
