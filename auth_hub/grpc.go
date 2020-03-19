package auth_hub

import (
	"context"
	"errors"

	"github.com/jinzhu/copier"

	"ustb_sso/protos"
)

type GrpcHandler struct{}

// Auth GRPC下使用认证的请求
func (g *GrpcHandler) Auth(ctx context.Context, req *protos.AuthReq) (*protos.AuthResponse, error) {
	select {
	case <-ctx.Done():
		return &protos.AuthResponse{}, errors.New("time out")
	default:
		res := doAuth(req.UnionId)
		resProto := new(protos.AuthResponse)
		err := copier.Copy(resProto, &res)
		if err != nil {
			return &protos.AuthResponse{}, err
		}
		return resProto, nil
	}
}

// Func GRPC下使用各种函数
func (g *GrpcHandler) Func(ctx context.Context, req *protos.FuncReq) (*protos.FuncResponse, error) {
	select {
	case <-ctx.Done():
		return &protos.FuncResponse{}, errors.New("time out")
	default:
		res, err := doFunc(req.FuncName, req.UnionId)
		if err != nil {
			return &protos.FuncResponse{}, err
		}
		return &protos.FuncResponse{Content: res}, nil
	}
}
