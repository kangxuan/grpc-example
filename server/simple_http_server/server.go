package main

import (
	"context"
	"google.golang.org/grpc"
	pb "gprc-example/proto"
	"log"
	"net/http"
	"strings"
)

type SearchService struct {
	pb.UnimplementedSearchServiceServer
}

func (s *SearchService) mustEmbedUnimplementedSearchServiceServer() {
	//TODO implement me
	panic("implement me")
}

const PORT = "9001"

func (s *SearchService) Search(cxt context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	return &pb.SearchResponse{Response: r.GetRequest() + " HTTP Server"}, nil
}

func main() {
	// 获取一个ServerMux锁
	mux := GetHTTPServeMux()

	server := grpc.NewServer()
	pb.RegisterSearchServiceServer(server, &SearchService{})

	err := http.ListenAndServe(":"+PORT,
		http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
				log.Println("sdf")
				server.ServeHTTP(w, r)
			} else {
				log.Println("sdf2")
				mux.ServeHTTP(w, r)

			}

			return
		}),
	)
	if err != nil {
		log.Fatalf("http.ListenAndServe err: %v", err)
	}
}

func GetHTTPServeMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("shanla: grpc-example"))
	})

	return mux
}
