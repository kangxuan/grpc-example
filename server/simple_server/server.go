package main

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
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
	// 根据服务端输入的证书文件和密钥构造 TLS 凭证
	c, err := credentials.NewServerTLSFromFile("../../conf/server.pem", "../../conf/server.key")
	if err != nil {
		log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	}

	// 创建一个grpc服务器
	// grpc.Creds()：返回一个 ServerOption，用于设置服务器连接的凭据。用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项
	server := grpc.NewServer(grpc.Creds(c))
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
