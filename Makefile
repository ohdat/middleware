.PHONY: all build run gotool clean help

BINARY="server"

all: check

check:
	go fmt ./
	go vet ./


help:
	@echo "make check - 运行 Go 工具 'fmt' and 'vet'"