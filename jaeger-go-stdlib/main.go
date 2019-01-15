package main

import (
	"flag"
	"log"

	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"time"
)

var (
	zipkinURL = flag.String("url",
		"http://localhost:9411/api/v1/spans", "Zipkin server URL")
	serverPort = flag.String("port", "8000", "server port")
	actorKind  = flag.String("actor", "server", "server or client")
)

const (
	server = "server"
	client = "client"
)

func main() {
	flag.Parse()
	cfg := config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
	}
	if *actorKind != server && *actorKind != client {
		log.Fatal("Please specify '-actor server' or '-actor client'")
	}

	transport,err:=jaeger.NewUDPTransport("127.0.0.1:6831",65000)
	if err!=nil {
		log.Println(err)
	}

	tracer, closer, err := cfg.New(
		"aaa",
		config.Reporter(jaeger.NewRemoteReporter(
			transport,
			jaeger.ReporterOptions.BufferFlushInterval(1*time.Second),
		)),
	)

	if *actorKind == server {
		runServer(tracer)
		return
	}

	runClient(tracer)

	closer.Close()
}
