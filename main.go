package main

import (
	"google.golang.org/grpc"
	"grpc-gateway/services"
	"net"
)

func main(){
	rpcServer := grpc.NewServer()
	services.RegisterGoodsServiceServer(rpcServer,new(services.GoodsService))
	lis,_ := net.Listen("tcp",":8081")
	rpcServer.Serve(lis)
}
