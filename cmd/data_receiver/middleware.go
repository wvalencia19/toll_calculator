package main

import (
	"time"

	"github.com/sirupsen/logrus"
	"github.com/wvalencia19/tolling/types"
)

type LoggingMiddleWare struct {
	next DataProducer
}

func NewLogMiddleWare(next DataProducer) *LoggingMiddleWare {
	return &LoggingMiddleWare{
		next: next,
	}
}

func (l *LoggingMiddleWare) ProduceData(data types.OBUData) error {
	start := time.Now()
	defer func(start time.Time) {
		logrus.WithFields(logrus.Fields{
			"obuID": data.OBUID,
			"lat":   data.Lat,
			"long":  data.Long,
			"took":  time.Since(start),
		}).Info("producing to kafka")
	}(start)

	return l.next.ProduceData(data)
}
