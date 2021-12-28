package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	hexTracing "grpc-gateway/tracing"

	"strings"
)

type HTTPConfig struct {
	Address string
}

type OauthGRPCConfig struct {
	Address string
}

type GRPCService struct {
	Name    string
	Address string
	Enable  bool
}

var v *viper.Viper

type Config struct {
	Debug           bool
	OauthGRPCConfig OauthGRPCConfig `mapstructure:"oauth"`
	HTTPConfig      HTTPConfig      `mapstructure:"http"`
	Services        map[string]*GRPCService
}

var OpenRancherTest bool
var ClientSleepTime int64

func init() {
	v = viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.SetEnvPrefix("GATEWAY")
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	v.AutomaticEnv()
	_ = v.ReadInConfig()

	traceDefault := map[string]string{
		"trace.service_name":  "oauth",
		"trace.agent_host":    "localhost",
		"trace.agent_port":    "6831",
		"trace.sample_type":   "const",
		"trace.sample_params": "1",
		"trace.disable":       "false",
	}
	setDefault(v, traceDefault)
	OpenRancherTest = v.GetBool("open_rancher_test")
	ClientSleepTime = v.GetInt64("client_sleep_time")
}

func GetConfig() (*Config, error) {
	var config = new(Config)
	err := v.Unmarshal(config)
	if err != nil {
		logrus.Fatal("read config fail: ", err.Error())
	}
	enabledServices := make(map[string]*GRPCService, 0)
	for name, service := range config.Services {
		if service.Enable {
			service.Name = name
			enabledServices[name] = service
		}
	}
	config.Services = enabledServices
	return config, err
}

func LoadTraceConfig() *hexTracing.HexTracerConfig {
	return &hexTracing.HexTracerConfig{
		ServiceName:  v.GetString("trace.service_name"),
		AgentHost:    v.GetString("trace.agent_host"),
		AgentPort:    v.GetString("trace.agent_port"),
		SampleType:   hexTracing.TraceSampleType(v.GetString("trace.sample_type")),
		SampleParams: v.GetFloat64("trace.sample_params"),
		Disable:      v.GetBool("trace.disable"),
	}
}
