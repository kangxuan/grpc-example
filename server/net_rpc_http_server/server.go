package main

import (
	"net"
	"net/http"
	"net/rpc"
)

type Args struct {
	X, Y int
}

type ServiceA struct{}

func (s *ServiceA) Add(a *Args, reply *int) error {
	*reply = a.X + a.Y
	return nil
}

func main() {
	service := new(ServiceA)
	// 注册RPC服务
	err := rpc.Register(service)
	if err != nil {
		panic(err)
	}

	// 基于HTTP协议
	rpc.HandleHTTP()
	lis, err := net.Listen("tcp", ":9001")
	if err != nil {
		panic(err)
	}

	err = http.Serve(lis, nil)
	if err != nil {
		panic(err)
	}
}
