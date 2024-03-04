package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/wvalencia19/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	httpListenAddr := os.Getenv("AGG_HTTP_ENDPOINT")
	grpcListenAddr := os.Getenv("AGG_GRPC_ENDPOINT")
	flag.Parse()

	store := makeStore()
	svc := NewInvoiceAggregator(store)
	svc = NewMetricsMiddleWare(svc)
	svc = NewLogMiddleWare(svc)

	go func() {
		log.Fatal(makeGRPCTransport(grpcListenAddr, svc))
	}()

	log.Fatal(makeHTPTransport(httpListenAddr, svc))
}

func makeHTPTransport(listenAddr string, svc Aggregator) error {
	aggregateMetricHandler := NewHTTPMetricHandler("aggregate")
	invoiceMetricHandler := NewHTTPMetricHandler("invoice")
	aggregatorHandler := makeHTTPHandler(aggregateMetricHandler.instrument(handleAggregate(svc)))
	invoiceHandler := makeHTTPHandler(invoiceMetricHandler.instrument(handleGetInvoice(svc)))

	http.HandleFunc("/aggregate", aggregatorHandler)
	http.HandleFunc("/invoice", invoiceHandler)
	http.Handle("/metrics", promhttp.Handler())

	fmt.Println("http transport running on port", listenAddr)

	return http.ListenAndServe(listenAddr, nil)
}

func makeGRPCTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("gRPC transport running on port", listenAddr)
	// make a TCP listener
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err
	}

	defer ln.Close()

	// make a new GRPC native server with(options)
	server := grpc.NewServer([]grpc.ServerOption{}...)
	// Register our GRPC server implementation to the GRPC package
	types.RegisterAggregatorServer(server, NewGRPCAggregatorServer(svc))
	return server.Serve(ln)
}

func makeStore() Storer {
	storeType := os.Getenv("APP_STORE_TYPE")
	switch storeType {
	case "memory":
		return NewMemoryStore()
	default:
		log.Fatalf("invalid store type, given %s", storeType)
		return nil
	}
}
