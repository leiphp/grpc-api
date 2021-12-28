package utils

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"

	log "github.com/sirupsen/logrus"
)

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func GetFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

type Middleware = func(http.Handler) http.Handler

func ChainHTTPMiddleware(wares ...Middleware) Middleware {
	if len(wares) == 0 {
		log.Fatal("no middleware!")
	}
	return func(handler http.Handler) http.Handler {
		var finalHandler = wares[len(wares)-1](handler)
		for i := len(wares) - 2; i >= 0; i-- {
			ware := wares[i]
			finalHandler = ware(finalHandler)
			fmt.Printf("wrapper %v middleware\n", GetFunctionName(ware))
		}
		return finalHandler
	}
}
