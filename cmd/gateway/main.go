package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/wvalencia19/tolling/cmd/aggregator/client"
)

type apiFunc func(w http.ResponseWriter, r *http.Request) error

func main() {
	listenAddr := flag.String("listenAddr", ":6000", "port for http server")
	aggregatorServiceAddr := flag.String("aggregatorServiceAddr", "http://localhost:3000", "address of the aggregator service")
	flag.Parse()

	client := client.NewHTTPClient(*aggregatorServiceAddr)
	invoiceHandler := NewInvoiceHandler(client)

	http.HandleFunc("/invoice", makeAPIFunc(invoiceHandler.HandleGetInvoice))

	logrus.Info("gateway running on port", *listenAddr)

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

type InvoiceHandler struct {
	client client.Client
}

func NewInvoiceHandler(c client.Client) *InvoiceHandler {
	return &InvoiceHandler{
		client: c,
	}
}

func (h *InvoiceHandler) HandleGetInvoice(w http.ResponseWriter, r *http.Request) error {
	queryValues := r.URL.Query()

	obuValue := queryValues.Get("obu")

	if obuValue == "" {
		return errors.Errorf("missing obu id")
	}

	obuID, err := strconv.Atoi(obuValue)
	if err != nil {
		return errors.Errorf("invalid obu id")
	}
	inv, err := h.client.GetInvoice(context.Background(), obuID)
	if err != nil {
		return err
	}
	return writeJSON(w, http.StatusOK, inv)
}

func writeJSON(w http.ResponseWriter, code int, v any) error {
	w.Header().Add("Content-Type", "application/json")

	w.WriteHeader(code)
	return json.NewEncoder(w).Encode(v)
}

func makeAPIFunc(fn apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(start time.Time) {
			logrus.WithFields(logrus.Fields{
				"took": time.Since(start),
				"uri":  r.RequestURI,
			}).Info("REQ :: ")
		}(time.Now())
		if err := fn(w, r); err != nil {
			writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}
}
