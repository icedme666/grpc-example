package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io/ioutil"
	"log"
	"net"
	pb "xiamei.guo/grpc-example/proto"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = "9001"

func main() {
	//LoadX509KeyPair: 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	cert, err := tls.LoadX509KeyPair("conf/ca/server/server.crt", "conf/ca/server/server.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
		return
	}
	certPool := x509.NewCertPool() //创建一个新的、空的 CertPool
	ca, err := ioutil.ReadFile("conf/ca/ca.crt")
	if err != nil {
		log.Fatalf("ioutil.ReadFile err: %v", err)
	}
	if ok := certPool.AppendCertsFromPEM(ca); !ok { //尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
		log.Fatalf("certPool.AppendCertsFromPEM err")
		return
	}
	//NewTLS: 构建基于 TLS 的 TransportCredentials 选项
	c := credentials.NewTLS(&tls.Config{ //Config 结构用于配置 TLS 客户端或服务器
		Certificates: []tls.Certificate{cert},        //设置证书链，允许包含一个或多个
		ClientAuth:   tls.RequireAndVerifyClientCert, //要求必须校验客户端的证书
		ClientCAs:    certPool,                       //设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	})

	//创建 gRPC Server 对象，可以理解为它是 Server 端的抽象对象
	server := grpc.NewServer(grpc.Creds(c))                  //Creds()返回一个 ServerOption，用于设置服务器连接的凭据。
	pb.RegisterSearchServiceServer(server, &SearchService{}) //将 SearchService（其包含需要被调用的服务端接口）注册到 gRPC Server 的内部注册中心。这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理

	lis, err := net.Listen("tcp", ":"+PORT) //创建 Listen，监听 TCP 端口
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis) //gRPC Server 开始 lis.Accept，直到 Stop 或 GracefulStop
}
