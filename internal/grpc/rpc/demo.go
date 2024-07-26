package rpc

import (
	"context"
	"fmt"
	"net/http"
)

var HelloService helloService

type helloService struct {
}

func (h *helloService) mustEmbedUnimplementedHelloServiceServer() {
}

func (h *helloService) Hello(ctx context.Context, in *HelloReq) (*BaseResponse, error) {
	name := "World"

	if in.Name != "" {
		name = in.Name
	}

	resp := &BaseResponse{
		Code:    http.StatusOK,
		Message: fmt.Sprintf("Hello %s", name),
	}

	return resp, nil
}
