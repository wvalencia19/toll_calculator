package main

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wvalencia19/tolling/types"
)

type LogMiddleWare struct {
	next Calculator
}

func NewLogMiddleWare(next Calculator) Calculator {
	return &LogMiddleWare{
		next: next,
	}
}

func (m *LogMiddleWare) CalculateDistance(data types.OBUData) (dist float64, err error) {
	defer func(start time.Time) {
		logrus.WithFields((logrus.Fields{
			"took": time.Since(start),
			"err":  err,
			"dist": fmt.Sprintf("%.2f", dist),
		})).Info("calculate distance")
	}(time.Now())

	dist, err = m.next.CalculateDistance(data)
	return
}
