package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"

	"github.com/wvalencia19/tolling/types"
	"google.golang.org/grpc"
)

func main() {
	httpListenAddr := flag.String("httplistenaddr", ":3000", "the listed address of the HTTP server")
	grpcListenAddr := flag.String("grpclistenaddr", ":3001", "the listed address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)
	svc = NewLogMiddleWare(svc)

	go func() {
		log.Fatal(makeGRPCTransport(*grpcListenAddr, svc))
	}()

	log.Fatal(makeHTPTransport(*httpListenAddr, svc))
}

func makeHTPTransport(listenAddr string, svc Aggregator) error {
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.HandleFunc("/invoice", handleGetInvoice(svc))
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

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
			return
		}
		if err := svc.AggregateDistance(distance); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}
}

func handleGetInvoice(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		values, ok := r.URL.Query()["obu"]

		if !ok {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "missing OBU ID"})
			return
		}

		obuID, err := strconv.Atoi(values[0])
		if err != nil {
			writeJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid OBU ID"})
			return
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}

		writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
