package main

import "log"

const kafkaConsumer = "obudata"

func main() {
	calculatorService := NewCalculatorService()
	calculatorService = NewLogMiddleWare(calculatorService)

	kafkaConsumer, err := NewKafkaConsumer(kafkaConsumer, calculatorService)

	if err != nil {
		log.Fatal(err)
	}
	kafkaConsumer.Start()
}
