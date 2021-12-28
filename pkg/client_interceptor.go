package pkg

import (
	"context"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"grpc-gateway/pkg/logger"
	"grpc-gateway/hexerror"
)

func APICodeFromMethodName(method string) string {
	apiCode := strings.Trim(method, "/")
	apiCode = strings.Replace(apiCode, "/", ".", -1)
	return apiCode
}

func GetTokenFromOutgoingContext(ctx context.Context) string {
	md, exist := metadata.FromOutgoingContext(ctx)
	if !exist {
		return ""
	}
	tokenList := md.Get("Authorization")
	if len(tokenList) == 0 {
		return ""
	}
	// value example: "Bearer token"
	tokenStr := tokenList[0]
	tokens := strings.Split(tokenStr, " ")
	if len(tokens) == 0 {
		return ""
	} else {
		return tokens[1]
	}
}

func AuthClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	start := time.Now()
	apiCode := APICodeFromMethodName(method)
	if strings.HasSuffix(method, "Ping") {
		err := invoker(ctx, method, req, reply, cc, opts...)
		logger.Pre().Infof("Invoked1 RPC method=%s; Duration=%s; Error=%v", method, time.Since(start), err)
		return err
	}
	APIRequestCounter.WithLabelValues(apiCode).Inc()
	token := GetTokenFromOutgoingContext(ctx)
	if token == "" {
		return hexerror.Unauthorized("no token")
	}
	userInfo, unAuthorizationErr := Authorization(ctx, token, apiCode)
	if unAuthorizationErr != nil {
		return unAuthorizationErr
	}
	authDuration := time.Since(start)
	kVPair := []string{
		"partner_id", strconv.FormatUint(userInfo.PartnerId, 10),
		"user_id", strconv.FormatUint(userInfo.UserId, 10),
		"scope_id", strconv.FormatUint(userInfo.ScopeId, 10),
	}
	ctx = metadata.AppendToOutgoingContext(ctx, kVPair...)
	err := invoker(ctx, method, req, reply, cc, opts...)
	logger.Pre().Infof("Invoked2 RPC method=%s; OauthDuration=%s; Duration=%s; Error=%v", method, authDuration, time.Since(start), err)
	return err
}

func DebugLogClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	err := invoker(ctx, method, req, reply, cc, opts...)
	logger.Pre().Infof("test log Invoked RPC method=%s;err = %v", method, err)
	return err
}
