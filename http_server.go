package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"grpc-gateway/helper"
	"grpc-gateway/services"
	"log"
	"net/http"
)

func main(){
	//http方式监听grpc服务，防止其他非grpc项目调用失败
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(helper.GetClientCreds())}
	err := services.RegisterGoodsServiceHandlerFromEndpoint(context.Background(), gwmux,"localhost:8081", opts)
	if err != nil {
		log.Fatal(err)
	}
	httpServer := &http.Server{
		Addr: ":8080",
		Handler: gwmux,
	}
	httpServer.ListenAndServe()
}
