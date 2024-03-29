# toll-calculator


This system consists of a suite of microservices designed to simulate the reception of GPS coordinates transmitted by vehicles. It processes these coordinates to calculate the distance each vehicle travels. Based on these calculations, it generates invoices that reflect the total distance traveled.

The core objective is to employ a variety of communication protocols between microservices, closely simulating a realistic operational environment. We aim to primarily utilize the native packages provided by Golang to minimize complexity and avoid the overhead introduced by external libraries, especially for HTTP client-server interactions. Furthermore, the project will incorporate advanced logging and metrics instrumentation techniques. We plan to use tools such as Logrus for logging, Prometheus for metrics collection, and Grafana for dashboard creation. 


![](diagram.png) 

## Help documentation
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

### Grafana

data source URL: `http://host.docker.internal:9090c`