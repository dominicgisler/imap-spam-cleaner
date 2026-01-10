FROM golang:1.24.11 AS build
WORKDIR /app
COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -a -installsuffix cgo -o dist/imap-spam-cleaner

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/dist/ .
ENTRYPOINT [ "/app/imap-spam-cleaner" ]
