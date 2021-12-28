package pkg

import (
	"context"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"

	"grpc-gateway/common"
	"grpc-gateway/config"
	"grpc-gateway/hexerror"
	oauthPb "grpc-gateway/protos/oauth"
)

type AuthFunc func(ctx context.Context, token string, apiCode string) (*common.UserInfo, hexerror.HexError)

func NewAuthFunc(oauthServerAddress string) AuthFunc {
	conn, err := grpc.Dial(oauthServerAddress,
		grpc.WithInsecure(),
	)
	if err != nil {
		log.Fatal("cant connect to oauth server: " + oauthServerAddress)
		return nil
	}
	client := oauthPb.NewHexOauthClient(conn)

	return func(ctx context.Context, token, apiCode string) (*common.UserInfo, hexerror.HexError) {
		req := &oauthPb.IntrospectTokenRequest{
			Token:   token,
			ApiCode: apiCode,
		}
		resp, err := client.IntrospectToken(ctx, req)
		if err != nil {
			log.Error(err)
			_hexErr := hexerror.HexErrorFromString(status.Convert(err).Message())
			return nil, _hexErr
		}
		userInfo := &common.UserInfo{
			UserId:    resp.UserId,
			PartnerId: resp.PartnerId,
			ScopeId:   resp.ScopeId,
		}
		return userInfo, nil
	}
}

func defaultAuthFunc() AuthFunc {
	conf, _ := config.GetConfig()
	return NewAuthFunc(conf.OauthGRPCConfig.Address)
}

var Authorization AuthFunc = defaultAuthFunc()
