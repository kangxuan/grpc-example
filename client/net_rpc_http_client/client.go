package main

import (
	"fmt"
	"log"
	"net/rpc"
)

type Args struct {
	X, Y int
}

var (
	reply  int
	reply2 int
)

func main() {
	// 建立RPC HTTP连接到服务端
	client, err := rpc.DialHTTP("tcp", ":9001")
	if err != nil {
		log.Fatalf("rpc.DialHTTP err: %v", err)
	}

	args := Args{
		X: 10,
		Y: 20,
	}

	// 同步调用
	err = client.Call("ServiceA.Add", &args, &reply)
	if err != nil {
		log.Fatalf("client.Call ServiceA.Add err: %v", err)
	}

	fmt.Printf("ServiceA.Add: %d+%d=%d\n", args.X, args.Y, reply)

	// 异步调用
	divCall := client.Go("ServiceA.Add", &args, &reply2, nil)
	replyCall := <-divCall.Done
	fmt.Println(replyCall.Error)
	fmt.Println(reply2)
}
