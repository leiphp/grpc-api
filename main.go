package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-gateway/services"
	"log"
	"net"
)

func main(){
	//服务端加载证书
	creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server_no_passwd.key")
	if err != nil {
		log.Fatal(err)
	}
	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterGoodsServiceServer(rpcServer,new(services.GoodsService))
	lis, err := net.Listen("tcp",":8081")
	if err != nil {
		log.Fatal(err)
	}else {
		log.Println("grpc服务开启成功")
	}
	rpcServer.Serve(lis)
}
