package pkg

import (
	"context"
	"encoding/json"
	"fmt"

	//"github.com/mitchellh/go-server-timing"
	"reflect"
	"runtime"
	"time"

	"grpc-gateway/common"
	"grpc-gateway/config"
	"grpc-gateway/pkg/logger"
	oauthpb "grpc-gateway/protos/oauth"

	"net/http"
	"strconv"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
	"grpc-gateway/hexerror"
	"grpc-gateway/utils"
)

func WriteResponse(w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(response)
	return err
}

type statusWriter struct {
	http.ResponseWriter
	status int
	length int
}

func (w *statusWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *statusWriter) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = 200
	}
	n, err := w.ResponseWriter.Write(b)
	w.length += n
	return n, err
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func DurationToStr(duration time.Duration) string {
	return fmt.Sprintf("%dms", uint64(duration/time.Millisecond))
}

func HTTPLog(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		st := time.Now()
		r = r.WithContext(context.WithValue(context.Background(), "st", st))
		start := time.Now()
		sw := statusWriter{ResponseWriter: w}
		next.ServeHTTP(&sw, r)
		duration := time.Since(start)
		logger.Pre().Infof("HTTP Request Info: duration=%s, status=%d method=%s url=%s",
			duration, sw.status, r.Method, r.URL)
	})
}

func HTTPAuthorization(next http.Handler) http.Handler {
	conf, _ := config.GetConfig()
	conn, err := grpc.Dial(conf.OauthGRPCConfig.Address, grpc.WithInsecure())
	utils.CheckError(err)
	client := oauthpb.NewHexOauthClient(conn)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if value, ok := r.Header["Authorization"]; ok && len(value) >= 1 {
			token := strings.TrimLeft(value[0], "Bearer ")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if token == "" {
				err := hexerror.Unauthorized("no token provide")
				errResp := common.HexErrToHTTPResponse(err)
				_ = WriteResponse(w, errResp)
				return
			}
			request := &oauthpb.IntrospectTokenRequest{Token: token}
			resp, err := client.IntrospectToken(ctx, request)

			if err != nil {
				hexErr := hexerror.HexErrorFromString(status.Convert(err).Message())
				_ = WriteResponse(w, common.HexErrToHTTPResponse(hexErr))
				return
			}
			r.Header["Partner-Id"] = []string{strconv.FormatUint(resp.PartnerId, 10)}
			r.Header["User-Id"] = []string{strconv.FormatUint(resp.UserId, 10)}
			next.ServeHTTP(w, r)
			return
		} else {
			err := hexerror.Unauthorized("no Authorization")
			errResp := common.HexErrToHTTPResponse(err)
			_ = WriteResponse(w, errResp)
			return
		}
	})
}
