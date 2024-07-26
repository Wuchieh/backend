package grpc

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"string_backend_0001/internal/grpc/rpc"
	"testing"
)

func TestRpc(t *testing.T) {
	client, err := grpc.NewClient("localhost:8082", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}

	defer client.Close()

	serviceClient := rpc.NewHelloServiceClient(client)
	req := &rpc.HelloReq{Name: ""}

	hello, err := serviceClient.Hello(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(hello.String())
}
