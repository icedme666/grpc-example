package main

import (
	"context"
	"github.com/EDDYCJY/go-grpc-example/proto"
	"google.golang.org/grpc"
	"io"
	"log"
)

const (
	PORT = "9002"
)

func main() {
	conn, err := grpc.Dial(":"+PORT, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc.Dail err: %v", err)
	}
	defer conn.Close()

	client := proto.NewStreamServiceClient(conn)
	err = printLists(client, &proto.StreamRequest{Pt: &proto.StreamPoint{Name: "gRPC Stream Client: List", Value: 2022}})
	if err != nil {
		log.Fatalf("printLists.err: %v", err)
	}
	err = printRecord(client, &proto.StreamRequest{Pt: &proto.StreamPoint{Name: "gRPC Stream Client: Record", Value: 2022}})
	if err != nil {
		log.Fatalf("printRecord.err: %v", err)
	}
	err = printRoute(client, &proto.StreamRequest{Pt: &proto.StreamPoint{Name: "gRPC Stream Client: Route", Value: 2022}})
	if err != nil {
		log.Fatalf("printRoute.err: %v", err)
	}
}

func printLists(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.List(context.Background(), r)
	if err != nil {
		return err
	}
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}

	return nil
}

func printRecord(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.Record(context.Background())
	if err != nil {
		return err
	}

	for n := 0; n < 6; n++ {
		err := stream.Send(r)
		if err != nil {
			return err
		}
	}

	resp, err := stream.CloseAndRecv() //stream.CloseAndRecv 和 stream.SendAndClose 是配套使用的流方法
	if err != nil {
		return err
	}
	log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)

	return nil
}

func printRoute(client proto.StreamServiceClient, r *proto.StreamRequest) error {
	stream, err := client.Route(context.Background())
	if err != nil {
		return err
	}
	for n := 0; n <= 6; n++ {
		err = stream.Send(r)
		if err != nil {
			return err
		}
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		log.Printf("resp: pj.name: %s, pt.value: %d", resp.Pt.Name, resp.Pt.Value)
	}
	//stream.CloseSend()

	return nil
}
