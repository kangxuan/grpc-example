package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
)

type Args struct {
	X, Y int
}

var (
	reply  int
	reply2 int
)

func main() {
	// 建立TCP连接
	conn, err := net.Dial("tcp", ":9001")
	if err != nil {
		log.Printf("rpc.Dial err: %v", err)
	}

	// 使用JSON协议
	client := rpc.NewClientWithCodec(jsonrpc.NewClientCodec(conn))

	args := &Args{
		X: 20,
		Y: 100,
	}

	// 同步调用
	err = client.Call("ServiceB.Reduce", args, &reply)
	if err != nil {
		log.Printf("client.Call err: %v", err)
	}
	fmt.Printf("ServiceB.Reduce: %d-%d=%d\n", args.X, args.Y, reply)

	// 异步调用
	divCall := client.Go("ServiceB.Reduce", args, &reply2, nil)
	replyCall := <-divCall.Done // 接收调用结果
	fmt.Println(replyCall.Error)
	fmt.Println(reply2)
}
