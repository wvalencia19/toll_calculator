# toll-calculator


### Kafka

https://github.com/confluentinc/confluent-kafka-go

go get github.com/confluentinc/confluent-kafka-go/kafka

https://developer.confluent.io/quickstart/kafka-local/

### Prometheus

https://github.com/prometheus/prometheus

Run with docker

```
docker run --name prometheus -d -p 127.0.0.1:9090:9090 prom/prometheus
```

installing prometheus golang client

```
go get github.com/prometheus/client_golang/prometheus
```


### Install protobuf compiler

Linux

```
sudo apt install -y protobuf-compiler
```

Mac

```
brew install protobuf
```

### Install GRPC and Protobuffer plugins for Golang.

Protobuffers

```
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
```

GRPC

```
go install  google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```

Set /go/bin directory in your path

```
PATH="${PATH}:${HOME}/go/bin"
```

Install package dependencies

```
go get google.golang.org/protobuf
go get google.golang.org/grpc
```