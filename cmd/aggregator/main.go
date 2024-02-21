package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"

	"github.com/wvalencia19/tolling/types"
)

func main() {
	listenAddr := flag.String("listenaddr", ":3000", "the listed address of the HTTP server")
	flag.Parse()

	store := NewMemoryStore()
	svc := NewInvoiceAggregator(store)

	makeHTPTransport(*listenAddr, svc)
}

func makeHTPTransport(listenAddr string, svc Aggregator) {
	fmt.Println("HTTP transport running on port", listenAddr)

	http.HandleFunc("/aggregate", handleAggregate(svc))
	http.ListenAndServe(listenAddr, nil)
}

func handleAggregate(svc Aggregator) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var distance types.Distance
		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
}
