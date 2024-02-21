obu:
	@go build -o bin/obu ./cmd/obu/
	@./bin/obu
receiver:
	@go build -o bin/data_receiver ./cmd/data_receiver/
	@./bin/data_receiver
calculator:
	@go build -o bin/distance_calculator ./cmd/distance_calculator/
	@./bin/distance_calculator

aggregator:
	@go build -o bin/aggregator ./cmd/aggregator/
	@./bin/aggregator	

.PHONY: obu, aggregator