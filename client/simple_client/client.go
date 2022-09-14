package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	"xiamei.guo/grpc-example/proto"
)

const PORT = "9001"

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure()) //创建与给定目标（服务端）的连接交互
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	defer conn.Close()

	client := proto.NewSearchServiceClient(conn)                           //创建 SearchService 的客户端对象
	resp, err := client.Search(context.Background(), &proto.SearchRequest{ //发送 RPC 请求，等待同步响应，得到回调后返回响应结果
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse()) //输出响应结果
}
