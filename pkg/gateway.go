package pkg

import (
	"context"
	"encoding/json"
	"grpc-gateway/common"
	"grpc-gateway/hexerror"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc/status"
	//"grpc-gateway/common"
)

type GatewayServer struct {
	Address string
	Mux     *runtime.ServeMux
	opts    []runtime.ServeMuxOption
}

func HexAuthMatcher(headerName string) (metadataName string, ok bool) {
	if headerName == "Partner-Id" {
		return "partner_id", true
	}
	if headerName == "User-Id" {
		return "user_id", true
	}
	return "", false
}

func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(http.StatusOK)
	hexErr := hexerror.HexErrorFromString(status.Convert(err).Message())
	resp := common.HexErrToHTTPResponse(hexErr)
	_ = json.NewEncoder(w).Encode(resp)
}

func NewGateWayServer(opts ...runtime.ServeMuxOption) (*GatewayServer, error) {
	gw := &GatewayServer{
		opts: []runtime.ServeMuxOption{
			runtime.WithIncomingHeaderMatcher(HexAuthMatcher),
			runtime.WithMarshalerOption("*", NewHexMarshaler()),
			//runtime.WithProtoErrorHandler()
		},
	}
	for _, opt := range opts {
		gw.opts = append(gw.opts, opt)
	}
	gw.Mux = runtime.NewServeMux(gw.opts...)
	runtime.HTTPError = CustomHTTPError
	return gw, nil
}

