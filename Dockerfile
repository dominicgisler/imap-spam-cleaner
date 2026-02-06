FROM --platform=$BUILDPLATFORM golang:1.25 AS build
WORKDIR /app
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
COPY . .
RUN GOOS=$TARGETOS GOARCH=$TARGETARCH CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-X main.version=${VERSION}" -o dist/imap-spam-cleaner

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/dist/ .
ENTRYPOINT [ "/app/imap-spam-cleaner" ]
