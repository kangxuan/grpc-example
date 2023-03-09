package main

import (
	"context"
	"google.golang.org/grpc"
	pb "gprc-example/proto"
	"log"
	"net"
	"strconv"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) mustEmbedUnimplementedSearchServiceServer() {
	//TODO implement me
	panic("implement me")
}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " Server"}, nil
}

const PORT = 9001

func main() {
	// 创建一个grpc服务器
	server := grpc.NewServer()
	// 将SearchService注册到grpc server的内部注册中心，这样可以在接受到请求时，通过内部的服务发现，发现该服务端接口并转接进行逻辑处理
	pb.RegisterSearchServiceServer(server, &SearchService{})
	// 创建TCP端口监听
	lis, err := net.Listen("tcp", ":"+strconv.Itoa(PORT))
	if err != nil {
		log.Fatalf("net.Listen err:%v", err)
	}
	// 启动服务
	err = server.Serve(lis)
	if err != nil {
		log.Fatalf("server.Serve err:%v", err)
	}
}
