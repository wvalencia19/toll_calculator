package main

import (
	"log"

	"github.com/wvalencia19/tolling/cmd/aggregator/client"
)

const (
	kafkaConsumer      = "obudata"
	aggregatorEndpoint = "http://127.0.0.1:4000"
)

func main() {
	calculatorService := NewCalculatorService()

	calculatorService = NewLogMiddleWare(calculatorService)
	httpClient := client.NewHTTPClient(aggregatorEndpoint)
	// grpcClient, err := client.NewGrpcClient(aggregatorEndpoint)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	kafkaConsumer, err := NewKafkaConsumer(kafkaConsumer, calculatorService, httpClient)

	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
