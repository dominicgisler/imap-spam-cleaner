.PHONY: default lint publish --tag
.DEFAULT_GOAL := default

default:
	@docker build -f Dockerfile -t dominicgisler/imap-spam-cleaner .

lint:
	@golangci-lint run

publish: --tag
	@docker push dominicgisler/imap-spam-cleaner:$(TAG)

--tag:
	@if [ "$(TAG)" = "" ]; then echo "TAG not set" && exit 1; else exit 0; fi
