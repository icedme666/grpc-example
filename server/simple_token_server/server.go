package main

import (
	"context"
	"log"
	"net/http"
	"strings"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"xiamei.guo/grpc-example/pkg/gtls"
	pb "xiamei.guo/grpc-example/proto"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
	auth *Auth
}
func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest)(*pb.SearchResponse, error){
	if err := s.auth.Check(ctx); err != nil{
		return nil, err
	}
	return &pb.SearchResponse{Response: r.GetRequest() + "Token Server"}, nil
}

const PORT = "9004"

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

type Auth struct {
	appKey string
	AppSecret string
}

func (a *Auth) Check(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 失败")
	}

	var (
		appKey string
		appSecret string
	)

	if value, ok := md["app_key"]; ok {
		appKey = value[0]
	}
	if value, ok := md["app_secret"]; ok{
		appSecret = value[0]
	}
	if appKey != a.GetAppKey() || appSecret != a.GetAppSecret(){
		return status.Errorf(codes.Unauthenticated, "自定义认证 Token 无效")
	}
	return nil
}

func (a *Auth) GetAppKey() string {
	return "gxm"
}

func (a *Auth) GetAppSecret() string {
	return "123456"
}