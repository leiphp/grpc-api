package cmd

import (
	"fmt"

	"github.com/kr/pretty"
	"grpc-gateway/config"
)

func Describe() {
	fmt.Println("let me describe")
	_, _ = pretty.Println(config.GetConfig())
	fmt.Println("trace config: ")
	//_, _ = pretty.Println(config.LoadTraceConfig())
}
