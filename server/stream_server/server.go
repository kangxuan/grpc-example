package main

import (
	"google.golang.org/grpc"
	pb "gprc-example/proto"
	"log"
	"net"
	"strconv"
)

type StreamService struct {
	pb.UnimplementedStreamServiceServer
}

func (s *StreamService) mustEmbedUnimplementedStreamServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *StreamService) List(r *pb.StreamRequest, stream pb.StreamService_ListServer) error {
	for n := 0; n <= 6; n++ {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
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

func (s *StreamService) Record(stream pb.StreamService_RecordServer) error {
	return nil
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	return nil
}

const PORT = 9002

func main() {
	server := grpc.NewServer()
	pb.RegisterStreamServiceServer(server, &StreamService{})

	lis, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve err: %v", err)
	}
}
