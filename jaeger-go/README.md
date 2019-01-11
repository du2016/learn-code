# jaeger tracer demo 

reference [tracing-examples](https://github.com/signalfx/tracing-examples)


# start jaeger

```
docker run -d --name jaeger \
  -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 \
  -p 5775:5775/udp \
  -p 6831:6831/udp \
  -p 6832:6832/udp \
  -p 5778:5778 \
  -p 16686:16686 \
  -p 14268:14268 \
  -p 9411:9411 \
  jaegertracing/all-in-one:1.8
```

# set env

14268 is the port to accept jaeger.thrift directly from clients

```
export JAEGER_ENDPOINT="http://127.0.0.1:14268/api/traces"
```
