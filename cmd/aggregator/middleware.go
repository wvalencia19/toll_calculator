package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wvalencia19/tolling/types"
)

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
