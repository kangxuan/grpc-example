package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	pb "gprc-example/proto"
	"log"
	"os"
	"strconv"
)

const PORT = 9001

func main() {
	//// 根据客户端输入的证书文件和密钥构造 TLS 凭证
	//c, err := credentials.NewClientTLSFromFile("../../conf/server.pem", "grpc-example")
	//if err != nil {
	//	log.Fatalf("credentials.NewClientTLSFromFile err: %v", err)
	//}

	cert, err := tls.LoadX509KeyPair("../../conf/client/client.pem", "../../conf/client/client.key")
	if err != nil {
		log.Fatalf("tls.LoadX509KeyPair err: %v", err)
	}

	certPool := x509.NewCertPool()
	ca, err := os.ReadFile("../../ca.pem")
	if err != nil {
		log.Fatalf("os.ReadFile err: %v", err)
	}

	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		log.Fatalf("certPool.AppendCertsFromPEM err")
	}

	c := credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cert},
		ServerName:   "grpc-example",
		RootCAs:      certPool,
	})
	// 连接server
	conn, err := grpc.Dial(":"+strconv.Itoa(PORT), grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err:%v", err)
	}
	defer func(conn *grpc.ClientConn) {
		err := conn.Close()
		if err != nil {
			panic(err)
		}
	}(conn)

	// 创建 SearchService 的客户端对象
	client := pb.NewSearchServiceClient(conn)
	// 发送 RPC 请求，等待同步响应，得到回调后返回响应结果
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err:%v", err)
	}
	log.Printf("resp:%s", resp.GetResponse())
}
