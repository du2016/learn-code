package main

import (
	"github.com/opentracing/opentracing-go"
	tracelog "github.com/opentracing/opentracing-go/log"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"log"
)

func main() {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		log.Println(err)
	}
	cfg.Sampler = &jaegercfg.SamplerConfig{
		Type:  "const",
		Param: 1,
	}
	cfg.ServiceName = "jaeger tracer demo"
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		log.Println(err)
	}
	defer closer.Close()
	//opentracing.SetGlobalTracer(tracer)

	parentSpan := tracer.StartSpan("root")
	defer parentSpan.Finish()
	parentSpan.LogFields(
		tracelog.String("hello", "world"),
	)
	parentSpan.LogKV("foo", "bar")

	childspan := tracer.StartSpan("child span", opentracing.ChildOf(parentSpan.Context()))
	defer childspan.Finish()
}
