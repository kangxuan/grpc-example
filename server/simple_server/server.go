package main

import (
	"context"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	pb "gprc-example/proto"
	"log"
	"net"
	"runtime/debug"
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
	//c, err := credentials.NewServerTLSFromFile("../../conf/server.pem", "../../conf/server.key")
	//if err != nil {
	//	log.Fatalf("credentials.NewServerTLSFromFile err: %v", err)
	//}

	//// 基于CA的TSL证书认证
	//// 从证书相关文件中读取和解析信息，得到证书公钥、密钥对
	//cert, err := tls.LoadX509KeyPair("../../conf/server/server.pem", "../../conf/server/server.key")
	//if err != nil {
	//	log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	//}
	//
	//// 创建一个新的、空的 CertPool
	//certPool := x509.NewCertPool()
	//ca, err := os.ReadFile("../../conf/ca.pem")
	//if err != nil {
	//	log.Fatalf("os.ReadFile err: %v", err)
	//}
	//// 尝试解析所传入的 PEM 编码的证书。如果解析成功会将其加到 CertPool 中，便于后面的使用
	//if ok := certPool.AppendCertsFromPEM(ca); !ok {
	//	log.Fatalf("certPool.AppendCertsFromPEM err")
	//}

	//// 构建基于 TLS 的 TransportCredentials 选项
	//c := credentials.NewTLS(&tls.Config{
	//	Certificates: []tls.Certificate{cert},        // 设置证书链，允许包含一个或多个
	//	ClientAuth:   tls.RequireAndVerifyClientCert, // 要求必须校验客户端的证书
	//	ClientCAs:    certPool,                       // 设置根证书的集合，校验方式使用 ClientAuth 中设定的模式
	//})

	// 创建一个grpc服务器
	// grpc.Creds()：返回一个 ServerOption，用于设置服务器连接的凭据。用于 grpc.NewServer(opt ...ServerOption) 为 gRPC Server 设置连接选项
	opts := []grpc.ServerOption{
		//grpc.Creds(c), // TSL认证
		grpc_middleware.WithUnaryServerChain( // 因为原grpc只支持一个拦截器，使用了go-grpc-middleware
			RecoveryInterceptor, // RPC 方法的异常保护和日志输出
			LoggingInterceptor,  // RPC 方法的入参出参的日志输出
		),
	}

	server := grpc.NewServer(opts...)
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

// RecoveryInterceptor 异常处理拦截器
func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}

// LoggingInterceptor 日志处理拦截器
func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("grpc method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("grpc method: %s, %v", info.FullMethod, resp)
	return resp, err
}
