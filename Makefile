.PHONY: all build run gotool clean help

BINARY="server"

all: check run

linux-build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build

windows-build:
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ${BINARY}.exe

mac-build:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ${BINARY}

run:
	@go run main.go serve

rpc:
	@go run main.go grpc

check:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 启动 http server项目"
	@echo "make linux-build - 编译 Go 代码, 生成Linux系统的二进制文件"
	@echo "make windows-build - 编译 Go 代码, 生成Windows系统的exe文件"
	@echo "make mac-build - 编译 Go 代码, 生成Mac系统的二进制文件"
	@echo "make run - 启动 http server"
	@echo "make rpc - 启动grpc server"
	@echo "make clean - 移除二进制文件"
	@echo "make check - 运行 Go 工具 'fmt' and 'vet'"