package main

import (
	"net"
	"net/rpc"
)

type Args struct {
	X, Y int
}

type ServiceB struct{}

func (s *ServiceB) Reduce(a *Args, reply *int) error {
	*reply = a.X - a.Y
	return nil
}

func main() {
	service := new(ServiceB)

	// 注册RPC服务
	err := rpc.Register(service)
	if err != nil {
		panic(err)
	}

	// 基于TCP协议
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}

	for {
		conn, _ := lis.Accept()
		rpc.ServeConn(conn)
	}
}
