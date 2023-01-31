package main

import (
	"context"
	"log"
	"time"
	"google.golang.org/grpc"
	"xiamei.guo/grpc-example/pkg/gtls"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"
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

	ctx, cancel := context.WithDeadline(context.Background(),time.Now().Add(time.Duration(5*time.Second)))  //返回最终上下文截止时间。
	defer cancel()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(ctx, &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		statusErr, ok:= status.FromError(err)  //返回 GRPCStatus 的具体错误码
		if ok{
			if statusErr.Code() == codes.DeadlineExceeded{
				log.Fatalln("client.Search err: deadline")
			}
		}
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse())
}