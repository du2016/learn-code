package main

import (
	"fmt"
	"io"
	"log"

	opentracing "github.com/opentracing/opentracing-go"
	jaegercfg "github.com/uber/jaeger-client-go/config"
	"os"
)

const (
	numParts  = 3
	partTerms = 100000000
)

func main() {
	tracer, tracerCloser, err := createTracer()
	if err != nil {
		log.Fatal("Could not create Jaeger tracer", err)
	}

	parentSpan := tracer.StartSpan("root")

	parentSpan.SetTag("totalTerms", 1)

	results := make(chan float64)
	for i := 0; i < numParts; i++ {
		part := i
		go func() {
			sp := tracer.StartSpan("calculatePart", opentracing.ChildOf(parentSpan.Context()))
			defer sp.Finish()

			sp.SetTag("part", part)

			termStart := part * partTerms
			results <- calculatePart(termStart, termStart+partTerms)
		}()
	}

	partsRemaining := numParts
	var piEstimate float64
	for partsRemaining > 0 {
		select {
		case part := <-results:
			piEstimate += part * 4.0
			partsRemaining--
		}
	}

	fmt.Printf("pi ~= %#v\n", piEstimate)

	parentSpan.Finish()

	tracerCloser.Close()
}

func calculatePart(termStart int, termEnd int) float64 {
	var partialSum float64
	for k := termStart; k < termEnd; k++ {
		numerator := 1.0
		if k%2 != 0 {
			numerator = -1.0
		}
		partialSum += numerator / float64(((2 * k) + 1))
	}
	return partialSum
}

func createTracer() (opentracing.Tracer, io.Closer, error) {
	cfg, err := jaegercfg.FromEnv()
	if err != nil {
		return nil, nil, err
	}

	if cfg.ServiceName == "" {
		cfg.ServiceName = "jaeger tracer demo"
	}

	cfg.Sampler = &jaegercfg.SamplerConfig{
		Type:  "const",
		Param: 1,
	}

	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}

	opentracing.SetGlobalTracer(tracer)

	return tracer, closer, nil
}
