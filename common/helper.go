package common

import (
	"grpc-gateway/hexerror"
)

func HexErrToHTTPResponse(hexerr hexerror.HexError) *HexResponse {
	resp := HexResponse{
		StatusCode:  1,
		Code:        hexerr.Code(),
		Source:      hexerr.Source(),
		Detail:      hexerr.Detail(),
		Description: hexerr.Description(),
	}
	return &resp
}
