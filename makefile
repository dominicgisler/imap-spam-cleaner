.PHONY: default lint docs publish --tag
.DEFAULT_GOAL := default

default: --tag
	@docker build --platform linux/amd64,linux/arm64 -f Dockerfile -t dominicgisler/imap-spam-cleaner:$(TAG) .

lint:
	@golangci-lint run

docs:
	@docker run --rm -it -p 8000:8000 -v $(PWD):/docs squidfunk/mkdocs-material

publish: --tag
	@docker push dominicgisler/imap-spam-cleaner:$(TAG)

--tag:
	@if [ "$(TAG)" = "" ]; then echo "TAG not set" && exit 1; else exit 0; fi
