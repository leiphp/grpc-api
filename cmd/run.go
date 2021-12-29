package cmd

import (
	"context"
	"fmt"
	"grpc-gateway/pkg"
	"grpc-gateway/utils"
	"math"
	"net/http"
	_ "net/http/pprof"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"grpc-gateway/config"
	//"grpc-gateway/pkg"
	commentPb "grpc-gateway/protos/comment"
	goodsPb "grpc-gateway/protos/goods"
	omsPb "grpc-gateway/protos/oms"

	hexTracing "grpc-gateway/tracing"
	grpc_tracing "grpc-gateway/tracing/grpc-tracing"
	http_tracing "grpc-gateway/tracing/http-tracing"
)

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type RegisterFunc func(context.Context, *runtime.ServeMux, string, []grpc.DialOption) error

func DebugClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()

	err := invoker(ctx, method, req, reply, cc, opts...)
	// Logic after invoking the invoker
	log.Infof("debug Invoked RPC method=%s; Duration=%s; Error=%v", method,
		time.Since(start), err)
	return err
}

func RunProf() {
	addr := "localhost:6060"
	log.Info("profile listen on " + addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Error("profile server error")
	}
	log.Info("profile server quit")
}

func Run() {
	conf, _ := config.GetConfig()
	if conf.Debug {
		go RunProf()
	}

	if len(conf.Services) == 0 {
		log.Fatal("no config services!")
	}
	tracer, closer, err := hexTracing.InitHexTracer(config.LoadTraceConfig())
	if err != nil {
		log.Fatal("init hex tracer fail: ", err)
	}
	defer func() {
		if err := closer.Close(); err != nil {
			log.Error("close tracer fail: ", err)
		}
	}()

	ctx := context.Background()
	server, err := pkg.NewGateWayServer()
	checkError(err)
	supportServiceMap := map[string]RegisterFunc{
		"comment":           commentPb.RegisterCommentServiceHandlerFromEndpoint,
		"goods":             goodsPb.RegisterGoodsServiceHandlerFromEndpoint,
		//"payflow":           payflowPb.RegisterPayflowHandlerFromEndpoint,
		"oms":               omsPb.RegisterOmsHandlerFromEndpoint,
		//"coupon":            couponPb.RegisterCouponServiceHandlerFromEndpoint,
		//"payment":           paymentPb.RegisterPaymentIntegrationHandlerFromEndpoint,
		//"product":           productPb.RegisterProductServiceHandlerFromEndpoint,
	}

	opts := []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(math.MaxInt32), grpc.MaxCallSendMsgSize(math.MaxInt32)),
		grpc.WithUnaryInterceptor(grpc_middleware.ChainUnaryClient(
			grpc_tracing.HexTracingClientInterceptor(tracer),
			pkg.AuthClientInterceptor,
		)),
	}
	for _, serviceConfig := range conf.Services {
		if !serviceConfig.Enable {
			continue
		}
		registerFunc, ok := supportServiceMap[serviceConfig.Name]
		if ok {
			err := registerFunc(ctx, server.Mux, serviceConfig.Address, opts)
			checkError(err)
			fmt.Printf("register %s complete (address: %s)\n", serviceConfig.Name, serviceConfig.Address)
		} else {
			log.Fatalf("do not support service %+v\n", serviceConfig.Name)
		}
	}
	go func() {
		mux := http.NewServeMux()
		mux.Handle("/metrics", promhttp.Handler())
		fmt.Println("metric service listen on 0.0.0.0:7070")
		checkError(http.ListenAndServe("0.0.0.0:7070", mux))
	}()

	log.Println("HTTP service listen on " + conf.HTTPConfig.Address)
	middleware := utils.ChainHTTPMiddleware(pkg.HTTPLog, http_tracing.HTTPTracingHandler)
	err = http.ListenAndServe(conf.HTTPConfig.Address, middleware(server.Mux))
	checkError(err)
}
