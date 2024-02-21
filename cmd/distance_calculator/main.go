package main

import (
	"log"

	"github.com/wvalencia19/tolling/cmd/aggregator/client"
)

const (
	kafkaConsumer      = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:3000/aggregate"
)

func main() {
	calculatorService := NewCalculatorService()

	calculatorService = NewLogMiddleWare(calculatorService)
	//httpClient := client.NewHTTPClient(aggregatorEndpoint)
	grpcClient, err := client.NewGrpcClient(aggregatorEndpoint)

	if err != nil {
		log.Fatal(err)
	}

	kafkaConsumer, err := NewKafkaConsumer(kafkaConsumer, calculatorService, grpcClient)

	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
