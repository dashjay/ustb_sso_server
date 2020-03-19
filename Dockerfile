FROM golang:alpine AS build

WORKDIR /go/sso_server
COPY go.mod go.sum /go/sso_server/

ENV GO111MODULE on
ENV GOPROXY https://goproxy.cn

RUN go mod download

COPY ./ /go/sso_server

RUN go build -o sso_server ./main.go

FROM alpine:latest

WORKDIR /opt

COPY --from=build /go/sso_server /opt/

EXPOSE 80/tcp 81/tcp

ENTRYPOINT ["./sso_server"]