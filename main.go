package main

import (
	//"google.golang.org/grpc"
	"grpc-gateway/cmd"
	//"grpc-gateway/helper"
	//"grpc-gateway/services"
	//"net"
	"os"
	"log"

	"github.com/urfave/cli"
)

//func main(){
//	//grpc服务监听方式
//	rpcServer := grpc.NewServer(grpc.Creds(helper.GetServerCreds()))
//	//注册商品服务
//	services.RegisterGoodsServiceServer(rpcServer,new(services.GoodsService))
//	//注册订单服务
//	services.RegisterOrdersServiceServer(rpcServer,new(services.OrdersService))
//	//用户积分服务
//	services.RegisterUserServiceServer(rpcServer,new(services.UserService))
//
//	//rpc服务监听方式
//	lis, err := net.Listen("tcp",":8081")
//	if err != nil {
//		log.Fatal(err)
//	}else {
//		log.Println("grpc服务开启成功")
//	}
//	rpcServer.Serve(lis)
//
//}

var (
	cliApp *cli.App
)

func init() {
	// Initialise a CLI app
	cliApp = cli.NewApp()
	cliApp.Name = "grpc-gateway"
	cliApp.Usage = "reverse proxy grpc service to http"
	cliApp.Author = "Lei Xiao Tian"
	cliApp.Email = "leixiaotian@100txy.com"
	cliApp.Version = "0.0.1"
}

func main() {
	// Set the CLI app commands
	cliApp.Commands = []cli.Command{
		{
			Name:  "run",
			Usage: "run proxy server",
			Action: func(c *cli.Context) {
				cmd.Run()
			},
		},
		{
			Name:  "describe",
			Usage: "describe the service which will be proxy",
			Action: func(c *cli.Context) {
				cmd.Describe()
			},
		},
	}

	if err := cliApp.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
