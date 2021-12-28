package http_tracing

import (
	"fmt"
	"net/http"
	"github.com/opentracing/opentracing-go"
)

func HTTPTracingHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		opName := r.Method + " " + fmt.Sprintf("%s", r.URL)
		var sp opentracing.Span
		parentContext, err := opentracing.GlobalTracer().Extract(
			opentracing.TextMap,
			opentracing.HTTPHeadersCarrier(r.Header))
		if err != nil {
			opName := r.Method + " " + fmt.Sprintf("%s", r.URL)
			sp = opentracing.StartSpan(opName)
		} else {
			sp = opentracing.StartSpan(opName, opentracing.ChildOf(parentContext))
		}
		sp.Tracer().Inject(sp.Context(),
			opentracing.HTTPHeaders,
			opentracing.HTTPHeadersCarrier(w.Header()),
		)
		r = r.WithContext(opentracing.ContextWithSpan(r.Context(), sp))
		defer sp.Finish()
		next.ServeHTTP(w, r)
	})
}
