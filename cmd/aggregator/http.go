package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/wvalencia19/tolling/types"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
)

type HTTPFunc func(http.ResponseWriter, *http.Request) error

type APIError struct {
	Code int
	Err  error
}

func (e APIError) Error() string {
	return e.Err.Error()
}

func makeHTTPHandler(fn HTTPFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			if apiErr, ok := err.(APIError); ok {

				writeJSON(w, apiErr.Code, map[string]string{"error": apiErr.Err.Error()})
			}
		}
	}
}

type HTTPMetricHandler struct {
	reqCounter    prometheus.Counter
	errReqCounter prometheus.Counter
	reqLatency    prometheus.Histogram
}

func NewHTTPMetricHandler(reqName string) *HTTPMetricHandler {
	return &HTTPMetricHandler{
		reqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_counter"),
			Name:      "aggregator",
		}),
		reqLatency: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "request_latency"),
			Name:      "aggregator",
			Buckets:   []float64{0.1, 0.5, 1},
		}),
		errReqCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: fmt.Sprintf("http_%s_%s", reqName, "err_counter"),
			Name:      "aggregator",
		}),
	}
}

func (h *HTTPMetricHandler) instrument(next HTTPFunc) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var err error
		defer func(start time.Time) {
			latency := time.Since(start).Seconds()
			logrus.WithFields(logrus.Fields{
				"latency": latency,
				"request": r.RequestURI,
				"err":     err,
			}).Info()
			h.reqLatency.Observe(latency)
			h.reqCounter.Inc()
			if err != nil {
				h.errReqCounter.Inc()
			}
		}(time.Now())

		err = next(w, r)
		return err
	}
}

func handleAggregate(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		var distance types.Distance
		if r.Method != "POST" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid http method %s", r.Method),
			}
		}

		if err := json.NewDecoder(r.Body).Decode(&distance); err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("failed to decode the response body %s", err),
			}

		}
		if err := svc.AggregateDistance(distance); err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}
		return writeJSON(w, http.StatusOK, map[string]string{"msg": "ok"})
	}
}

func handleGetInvoice(svc Aggregator) HTTPFunc {
	return func(w http.ResponseWriter, r *http.Request) error {
		if r.Method != "GET" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid http method %s", r.Method),
			}
		}
		queryValues := r.URL.Query()

		obuValue := queryValues.Get("obu")

		if obuValue == "" {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("missing OBU id"),
			}
		}

		obuID, err := strconv.Atoi(obuValue)
		if err != nil {
			return APIError{
				Code: http.StatusBadRequest,
				Err:  fmt.Errorf("invalid OBU id %s", obuValue),
			}
		}

		invoice, err := svc.CalculateInvoice(obuID)
		if err != nil {
			return APIError{
				Code: http.StatusInternalServerError,
				Err:  err,
			}
		}

		return writeJSON(w, http.StatusOK, invoice)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Add("Content-Type", "application/json")

	return json.NewEncoder(w).Encode(v)
}
