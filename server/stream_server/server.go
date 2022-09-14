package main

import (
	"github.com/EDDYCJY/go-grpc-example/proto"
	"google.golang.org/grpc"
	"io"
	"log"
	"net"
)

type StreamService struct{}

const (
	PORT = "9002"
)

func main() {
	server := grpc.NewServer()
	proto.RegisterStreamServiceServer(server, &StreamService{})
	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}
	server.Serve(lis)
}

func (s *StreamService) List(r *proto.StreamRequest, stream proto.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&proto.StreamResponse{
			Pt: &proto.StreamPoint{
				Name:  r.Pt.Name,
				Value: r.Pt.Value + int32(n),
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *StreamService) Record(stream proto.StreamService_RecordServer) error {
	for {
		//对每一个 Recv 都进行了处理，当发现 io.EOF (流关闭) 后，需要将最终的响应结果发送给客户端，同时关闭正在另外一侧等待的 Recv
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&proto.StreamResponse{Pt: &proto.StreamPoint{Name: "gRPC Stream Server: Record", Value: 1}})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}

func (s *StreamService) Route(stream proto.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&proto.StreamResponse{
			Pt: &proto.StreamPoint{
				Name:  "gPRC Stream Client: Route",
				Value: int32(n),
			},
		})
		if err != nil {
			return err
		}
		r, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		n++
		log.Printf("stream.Recv pt.name: %s, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}

	return nil
}
