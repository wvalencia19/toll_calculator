obu:
	@go build -o bin/obu ./cmd/obu/
	@./bin/obu
receiver:
	@go build -o bin/data_receiver ./cmd/data_receiver/
	@./bin/data_receiver	

.PHONY: obu