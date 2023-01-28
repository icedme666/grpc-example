package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	pb "xiamei.guo/grpc-example/proto"
)

const PORT = "9001"

func main() {
	//根据客户端输入的证书文件和密钥构造 TLS 凭证。
	cert, err := tls.LoadX509KeyPair("conf/client/client.pem", "conf/client/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile("conf/ca.pem")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}
	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "grpc-example",
		RootCAs:      certPool,
	})

	////创建与给定目标（服务端）的连接交互
	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c)) //WithTransportCredentials()返回一个配置连接的 DialOption 选项
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)                           //创建 SearchService 的客户端对象
	resp, err := client.Search(context.Background(), &pb.SearchRequest{ //发送 RPC 请求，等待同步响应，得到回调后返回响应结果
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse()) //输出响应结果
}
