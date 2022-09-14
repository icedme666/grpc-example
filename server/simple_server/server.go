package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"net"
	"xiamei.guo/grpc-example/proto"
)

type SearchService struct {
	proto.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *proto.SearchRequest) (*proto.SearchResponse, error) {
	return &proto.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	server := grpc.NewServer()                                  //创建 gRPC Server 对象，可以理解为它是 Server 端的抽象对象
	proto.RegisterSearchServiceServer(server, &SearchService{}) //将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理

	lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis) //gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
}
