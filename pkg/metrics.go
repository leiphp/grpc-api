package pkg

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const nameSpace = "gateway"

var APIRequestCounter = promauto.NewCounterVec(prometheus.CounterOpts{
	Namespace: nameSpace,
	Name:      "api_request_counter",
	Help:      "gateway api request counter",
}, []string{"api_name"})
