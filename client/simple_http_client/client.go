package main

import (
	"context"
	"log"
	"google.golang.org/grpc"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

const PORT = "9003"

func main() {
	tlsCient := gtls.Client{
		ServerName: "grpc-example",
		CertFile: "conf/ca/server/server.crt",
	}
	c, err := tlsCient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsCient.GetTLSCredentials err %v", err)
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}