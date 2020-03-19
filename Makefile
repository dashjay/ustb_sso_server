.PHONY: build test clean proto

all: build

build:
	go build main.go

test:
	go test test/*.go -v

clean:
	+rm -r ./main

proto:
	protoc -Iprotos protos/BaseMessage.proto --go_out=plugins=grpc:protos
