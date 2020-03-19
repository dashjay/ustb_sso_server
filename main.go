package main

import (
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"

	"github.com/dashjay/logging"
	"google.golang.org/grpc"

	"ustb_sso/auth_hub"
	"ustb_sso/env"
	"ustb_sso/protos"
)

func main() {
	// http接口
	http.HandleFunc("/auth", auth_hub.DoAuthHTTP)
	http.HandleFunc("/func", auth_hub.DoFuncHTTP)
	logging.Info("http-sso server starting ")
	go func() {
		err := http.ListenAndServe(fmt.Sprintf(":%s", env.HTTPPort), nil)
		if err != nil {
			panic(err)
		}
	}()

	// grpc server
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", env.GRPCPort))
	if err != nil {
		panic(err)
	}
	s := grpc.NewServer()
	protos.RegisterAuthServerServer(s, &auth_hub.GrpcHandler{})
	logging.Info("grpc-sso server starting")

	err = s.Serve(lis)
	if err != nil {
		panic(err)
	}
}
