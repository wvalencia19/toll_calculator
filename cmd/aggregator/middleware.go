package main

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/sirupsen/logrus"
	"github.com/wvalencia19/tolling/types"
)

type MatrixMiddleWare struct {
	reqCounterAgg  prometheus.Counter
	reqLatencyAgg  prometheus.Histogram
	reqCounterCalc prometheus.Counter
	reqLatencyCalc prometheus.Histogram
	errCounterAgg  prometheus.Counter
	errCounterCalc prometheus.Counter
	next           Aggregator
}

func NewMetricsMiddleWare(next Aggregator) *MatrixMiddleWare {
	errCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "error_aggregator_request_counter",
		Name:      "aggregate",
	})
	errCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "error_aggregator_request_counter",
		Name:      "aggregate",
	})
	reqCounterAgg := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "aggregate",
	})

	reqLatencyAgg := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "aggregate",
		Buckets:   []float64{0.1, 0.5, 1},
	})
	reqCounterCalc := promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "aggregator_request_counter",
		Name:      "calculate",
	})

	reqLatencyCalc := promauto.NewHistogram(prometheus.HistogramOpts{
		Namespace: "aggregator_request_latency",
		Name:      "calculate",
		Buckets:   []float64{0.1, 0.5, 1},
	})

	return &MatrixMiddleWare{
		reqCounterAgg:  reqCounterAgg,
		reqLatencyAgg:  reqLatencyAgg,
		reqCounterCalc: reqCounterCalc,
		reqLatencyCalc: reqLatencyCalc,
		errCounterAgg:  errCounterAgg,
		errCounterCalc: errCounterCalc,
		next:           next,
	}
}

func (m *MatrixMiddleWare) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		m.reqLatencyAgg.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterAgg.Inc()
		if err != nil {
			m.errCounterAgg.Inc()
		}
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *MatrixMiddleWare) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {
	defer func(start time.Time) {
		m.reqLatencyCalc.Observe(float64(time.Since(start).Seconds()))
		m.reqCounterCalc.Inc()
		if err != nil {
			m.errCounterCalc.Inc()
		}
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}

type LogMiddleWare struct {
	next Aggregator
}

func NewLogMiddleWare(next Aggregator) Aggregator {
	return &LogMiddleWare{
		next: next,
	}
}

func (m *LogMiddleWare) AggregateDistance(distance types.Distance) (err error) {
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"took":  time.Since(start),
			"err":   err,
			"obuID": distance.OBUID,
		}).Info("aggregate distance")
	}(time.Now())
	err = m.next.AggregateDistance(distance)
	return
}

func (m *LogMiddleWare) CalculateInvoice(obuID int) (inv *types.Invoice, err error) {

	defer func(start time.Time) {
		var (
			distance float64
			amount   float64
		)
		if inv != nil {
			distance = inv.TotalDistance
			amount = inv.TotalDistance
		}
		logrus.WithFields(logrus.Fields{
			"took":     time.Since(start),
			"err":      err,
			"obuID":    obuID,
			"distance": distance,
			"amount":   amount,
		}).Info("calculate invoice")
	}(time.Now())
	inv, err = m.next.CalculateInvoice(obuID)
	return
}
