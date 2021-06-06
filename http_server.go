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
	grpcEndPoint := "localhost:8081"
	opts := []grpc.DialOption{grpc.WithTransportCredentials(helper.GetClientCreds())}
	//商品服务
	err := services.RegisterGoodsServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndPoint, opts)
	if err != nil {
		log.Fatal(err)
	}
	//商订单服务
	err = services.RegisterOrdersServiceHandlerFromEndpoint(context.Background(), gwmux, grpcEndPoint, opts)
	if err != nil {
		log.Fatal(err)
	}

	httpServer := &http.Server{
		Addr: ":8080",
		Handler: gwmux,
	}
	log.Println("grpc代理http服务开启成功")
	httpServer.ListenAndServe()
}
