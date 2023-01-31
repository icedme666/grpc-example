package main

import (
	"context"
	"log"
	"google.golang.org/grpc"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

const PORT = "9004"

type Auth struct {
	AppKey string
	AppSecret string
}

func (a *Auth) GetRequestMetadata(ctx context.Context, uri ...string)(map[string]string, error){
	return map[string]string{"app_key": a.AppKey, "app_secret": a.AppSecret}, nil
}

func (a *Auth) RequireTransportSecurity() bool {
	return true
}

func main() {
	tlsCient := gtls.Client{
		ServerName: "grpc-example",
		CertFile: "conf/ca/server/server.crt",
	}
	c, err := tlsCient.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsCient.GetTLSCredentials err %v", err)
	}
	
	auth := Auth{
		AppKey: "gxm",
		AppSecret: "123456",
	}
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c), grpc.WithPerRPCCredentials(&auth))
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