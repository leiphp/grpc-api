package main

import (
	"google.golang.org/grpc"
	"grpc-gateway/helper"
	"grpc-gateway/services"
	"log"
	"net"
)

func main(){
	//grpc服务监听方式
	rpcServer := grpc.NewServer(grpc.Creds(helper.GetServerCreds()))
	//注册商品服务
	services.RegisterGoodsServiceServer(rpcServer,new(services.GoodsService))
	//注册订单服务
	services.RegisterOrdersServiceServer(rpcServer,new(services.OrdersService))

	//rpc服务监听方式
	lis, err := net.Listen("tcp",":8081")
	if err != nil {
		log.Fatal(err)
	}else {
		log.Println("grpc服务开启成功")
	}
	rpcServer.Serve(lis)

}
