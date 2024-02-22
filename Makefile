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

gateway:
	@go build -o bin/gateway ./cmd/gateway/
	@./bin/gateway
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative types/ptypes.proto

.PHONY: obu, aggregator