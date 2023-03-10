package main

import (
	"google.golang.org/grpc"
	pb "gprc-example/proto"
	"io"
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
	for {
		r, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.StreamResponse{
				Pt: &pb.StreamPoint{
					Name:  "grpc stream server: Record",
					Value: 1,
				},
			})
		}
		if err != nil {
			return err
		}
		log.Printf("stream.Recv pt.name: %v, pt.value: %d", r.Pt.Name, r.Pt.Value)
	}
}

func (s *StreamService) Route(stream pb.StreamService_RouteServer) error {
	n := 0
	for {
		err := stream.Send(&pb.StreamResponse{
			Pt: &pb.StreamPoint{
				Name:  "grpc stream client: Route",
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
