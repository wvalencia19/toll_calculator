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

	kafkaConsumer, err := NewKafkaConsumer(kafkaConsumer, calculatorService, client.NewClient(aggregatorEndpoint))

	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
