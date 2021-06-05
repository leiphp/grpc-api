package main

import (
	"crypto/tls"
	"crypto/x509"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"grpc-gateway/services"
	"io/ioutil"
	"log"
	"net"
)

func main(){
	//服务端加载证书
	//creds, err := credentials.NewServerTLSFromFile("keys/server.crt", "keys/server_no_passwd.key")
	//if err != nil {
	//	log.Fatal(err)
	//}

	//服务端加载证书,双向验证
	cert, _ := tls.LoadX509KeyPair("cert/server.pem", "cert/server.key")
	certPool := x509.NewCertPool()
	ca, _ := ioutil.ReadFile("cert/ca.pem")
	certPool.AppendCertsFromPEM(ca)

	creds := credentials.NewTLS(&tls.Config{
		Certificates: 				 []tls.Certificate{cert}, //服务端证书
		ClientAuth:                  tls.RequireAndVerifyClientCert, //双向验证
		ClientCAs:                   certPool,
	})

	rpcServer := grpc.NewServer(grpc.Creds(creds))
	services.RegisterGoodsServiceServer(rpcServer,new(services.GoodsService))

	//rpc服务监听方式
	lis, err := net.Listen("tcp",":8081")
	if err != nil {
		log.Fatal(err)
	}else {
		log.Println("grpc服务开启成功")
	}
	rpcServer.Serve(lis)

	//http服务监听方式
}
