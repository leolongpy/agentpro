.PHONY: all build run stop gotool clean help

BINARY="agentx_linux"

all: gotool build

build:
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o ${BINARY}

run:
	@go run ./main.go start

stop:
	@go run ./main.go stop

gotool:
	go fmt ./
	go vet ./

clean:
	@if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

help:
	@echo "make - 格式化 Go 代码, 并编译生成二进制文件"
	@echo "make build - 编译 Go 代码, 生成二进制文件"
	@echo "make run - 直接运行 Go 代码"
	@echo "make stop - 暂停 Go 代码"
	@echo "make clean - 移除二进制文件和 vim swap files"
	@echo "make gotool - 运行 Go 工具 'fmt' and 'vet'"