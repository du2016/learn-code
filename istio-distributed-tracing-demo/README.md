# istio distributed treacing demo

## requirement

- istio on k8s with auto injection

## architecture

ingress --> democlient-->demoserver

## build

GOOS=linux go build demo-client.go

GOOS=linux go build demo-server.go

## deploy

```
mv demo-client demo-server /opt/
kubectl create -f demo.yaml
```