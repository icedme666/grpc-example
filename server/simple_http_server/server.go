package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"google.golang.org/grpc"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}
func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest)(*pb.SearchResponse, error){
	return &pb.SearchResponse{Response: r.GetRequest() + "HTTP Server"}, nil
}

const PORT = "9003"

func main(){
	certFile := "conf/ca/server/server.crt"
	keyFile := "conf/ca/server/server.key"
	tlsServer := gtls.Server{
		CertFile: certFile,
		KeyFile: keyFile,
	}

	c, err := tlsServer.GetTLSCredentials()
	if err != nil {
		log.Fatalf("tlsServer.GetTLSCredentials err: %v", err)
	}

	mux := GetHTTPServeMux()

	server := grpc.NewServer(grpc.Creds(c))
	pb.RegisterSearchServiceServer(server, &SearchService{})

	http.ListenAndServeTLS(
		":"+PORT, certFile, 
		keyFile, 
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc"){
				server.ServeHTTP(w, r)	
			} else {
				mux.ServeHTTP(w, r)
			}
			return
		}),
    )
}

func GetHTTPServeMux() *http.ServeMux{
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("gxm: go-grpc-example"))
	})
	return mux
}