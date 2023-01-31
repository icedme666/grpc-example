package main

import (
	"context"
	"log"
	"net"
	//"os"
	"google.golang.org/grpc"
	"github.com/openzipkin/zipkin-go"
	zipkingrpc "github.com/openzipkin/zipkin-go/middleware/grpc"
	//logreporter "github.com/openzipkin/zipkin-go/reporter/log"
	httpreport "github.com/openzipkin/zipkin-go/reporter/http"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}
func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest)(*pb.SearchResponse, error){
	return &pb.SearchResponse{Response: r.GetRequest() + "HTTP Server"}, nil
}

const PORT = "9005"

func main(){
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

	tlsServer := gtls.Server{
		CaFile: "conf/ca/ca.crt",
		CertFile: "conf/ca/server/server.crt",
		KeyFile: "conf/ca/server/server.key",
	}

	c, err := tlsServer.GetCredentialsByCA()
	if err != nil {
		log.Fatalf("tlsServer.GetCredentialsByCA err: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc.StatsHandler(zipkingrpc.NewServerHandler(tracer)),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}