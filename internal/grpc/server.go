package grpc

import (
	"google.golang.org/grpc"
	"net"
	"string_backend_0001/internal/grpc/rpc"
)

var (
	rpcServer *grpc.Server
	ser       net.Listener
)

func Run() error {
	//creds, err := credentials.NewServerTLSFromFile("./server.crt", "./server.key")
	//if err != nil {
	//	return err
	//}

	rpcServer = grpc.NewServer(
	// grpc.Creds(creds),
	)

	rpc.RegisterHelloServiceServer(rpcServer, &rpc.HelloService)

	Ser, err := net.Listen("tcp", ":8082")

	if err != nil {
		return err
	}

	ser = Ser

	return rpcServer.Serve(ser)
}

func Close() error {
	if rpcServer != nil {
		rpcServer.Stop()
	}

	if ser == nil {
		return nil
	}

	return ser.Close()
}
