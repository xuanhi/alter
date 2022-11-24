.PHONY: alpine build
alpine:
	CGO_ENABLED=0 go build -o ./example/alterwebhook/feishu ./example/alterwebhook/main.go
build:
	go build -o feishu ./example/alterwebhook/main.go