package main

import (
	"context"
	"google.golang.org/grpc"
	"log"
	//"os"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	//logreporter "github.com/openzipkin/zipkin-go/reporter/log"
	httpreport "github.com/openzipkin/zipkin-go/reporter/http"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

const PORT = "9005"

func main() {
	// set up a span reporter
	//reporter := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	url := "http://127.0.0.1:9411/api/v2/spans"
	reporter := httpreport.NewReporter(url)
	defer func() {
		_ = reporter.Close()
	}()

	// create our local service endpoint
	endpoint, err := zipkin.NewEndpoint("simple_zipkin_server", "127.0.0.1:9005")
	if err != nil {
		log.Fatalf("unable to create local endpoint: %+v\n", err)
	}

	// initialize our tracer
	tracer, err := zipkin.NewTracer(reporter, zipkin.WithLocalEndpoint(endpoint))
	if err != nil {
		log.Fatalf("zipkin.NewTracer err: %v", err)
	}

	tlsClient := gtls.Client{
		ServerName: "grpc-example",
		CaFile: "conf/ca/ca.crt",
		CertFile: "conf/ca/server/server.crt",
		KeyFile: "conf/ca/server/server.key",
	}
	
	c, err := tlsClient.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("GetCredentialsByCA err: %v", err)
	}

	//创建与给定目标（服务端）的连接交互
	conn, err := grpc.Dial(
		":"+PORT, 
		grpc.WithStatsHandler(zipkingrpc.NewClientHandler(tracer)),
		grpc.WithTransportCredentials(c),
	)
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)                           //创建 SearchService 的客户端对象
	resp, err := client.Search(context.Background(), &pb.SearchRequest{ //发送 RPC 请求，等待同步响应，得到回调后返回响应结果
		Request: "gRPC ZipKin",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}
	log.Printf("resp: %s", resp.GetResponse()) //输出响应结果
}
