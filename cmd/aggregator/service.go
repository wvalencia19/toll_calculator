package main

import (
	"fmt"

	"github.com/wvalencia19/tolling/types"
)

type Aggregator interface {
	AggregateDistance(types.Distance) error
}

type Storer interface {
	Insert(types.Distance) error
}

type InvoiceAggregator struct {
	store Storer
}

func NewInvoiceAggregator(s Storer) *InvoiceAggregator {
	return &InvoiceAggregator{
		store: s,
	}
}

func (i *InvoiceAggregator) AggregateDistance(distance types.Distance) error {
	fmt.Println("processing and inserting distance in the storage", distance)

	return i.store.Insert(distance)
}
