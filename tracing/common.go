package hex_tracing

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
)

type TraceSampleType string

const (
	CONST         TraceSampleType = "const"
	PROBABILISTIC                 = "probabilistic"
	RATELIMITING                  = "rateLimiting"
)

type HexTracerConfig struct {
	// 服务名称，请设置为你的项目名称。
	ServiceName string

	// Agent的地址。Agent使用来的接受Tracing数据的程序。Agent由运维搭建。
	AgentHost string
	AgentPort string

	// Tracing并不是对每次的调用都进行采样，你可以用以下两个参数设置的采样策略。
	// 采样策略：const(采样所有), probabilistic(按概率), rateLimiting(按速度限制）
	SampleType TraceSampleType

	// 采样参数，用于细化采样策略
	// 当采样策略为const时：1表示采样所有， 0表示不采样
	// 当采样策略为probabilistic时，指定一个在0～1之间的小数，表示采样概率。
	// 当采样策略为rateLimiting时，表示每秒发送的Span数。（Span时tracing的一个底层概念，在这里你可以就把他当成tracing)
	SampleParams float64
	// 是否打开日志
	LogEnable bool
	// 是否关闭tracing
	Disable bool
}

func (c *HexTracerConfig) Validate() {
	if c.Disable {
		return
	}
	if c.ServiceName == "" {
		log.Fatal("HexTracerConfig ServiceName can not be empty!")
	}

	if c.SampleType == "" {
		c.SampleType = PROBABILISTIC
	} else {
		var validSameType bool
		for _, v := range []TraceSampleType{CONST, PROBABILISTIC, RATELIMITING} {
			if c.SampleType == v {
				validSameType = true
				break
			}
		}
		if !validSameType {
			msg := "Fatal Error: WRONG HexTracer.SampleType, valid value should be 'cons' or 'probabilistic' or 'rateLimiting'"
			log.Fatal(msg)
		}
	}

	msg := "Fatal Error: HexTracer.SampleParams can not be %v when HexTracer.SampleType is '%s'"
	msg = fmt.Sprintf(msg, c.SampleParams, c.SampleType)

	if c.SampleType == PROBABILISTIC && (!(c.SampleParams >= 0.0 && c.SampleParams <= 1)) {
		log.Fatal(msg)
	}
	addr := c.AgentHost + ":" + c.AgentPort
	log.Println("trace log will send to " + addr + " use UDP protocol")
}

func jaegerConfigFromHexTraceConfig(hexConf *HexTracerConfig) *jaegercfg.Configuration {
	hexConf.Validate()
	var jaegerConf = &jaegercfg.Configuration{
		ServiceName: hexConf.ServiceName,
		Disabled:    hexConf.Disable,
		Sampler: &jaegercfg.SamplerConfig{
			Type:  string(hexConf.SampleType),
			Param: hexConf.SampleParams,
		},
		Headers: &jaeger.HeadersConfig{
			TraceContextHeaderName: "hex-trace-id",
			JaegerDebugHeader: "hex-debug-id",
		},
		Reporter: &jaegercfg.ReporterConfig{
			LogSpans:            hexConf.LogEnable,
			BufferFlushInterval: 1 * time.Second,
			LocalAgentHostPort:  hexConf.AgentHost + ":" + hexConf.AgentPort,
		},
	}
	return jaegerConf
}

func InitHexTracer(hexTracerConfig *HexTracerConfig) (tracer opentracing.Tracer, closer io.Closer, err error) {
	jConf := jaegerConfigFromHexTraceConfig(hexTracerConfig)
	tracer, closer, err = jConf.NewTracer(jaegercfg.Logger(jaeger.StdLogger))
	opentracing.SetGlobalTracer(tracer)
	return
}
