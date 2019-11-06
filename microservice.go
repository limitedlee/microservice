package main

import (
	"fmt"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"net/http"
	"strings"
)

type MicService struct {
	GrpcServer *grpc.Server
}

func (m *MicService) NewServer() {
	m.GrpcServer = grpc.NewServer()
}

func (m *MicService) Start() {
	mux := http.NewServeMux()
	mux.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})
	http.ListenAndServe("127.0.0.1:8888", grpcHandleFunc(m.GrpcServer, mux))
}

func grpcHandleFunc(grpcServer *grpc.Server, otherHander http.Handler) http.Handler {
	return h2c.NewHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ProtoMajor == 2 && strings.Contains(r.Header.Get("Content-Type"), "application/grpc") {
			grpcServer.ServeHTTP(w, r)
		} else {
			otherHander.ServeHTTP(w, r)
		}
	}), &http2.Server{})
}
