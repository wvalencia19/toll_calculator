package aggservice

import (
	"context"

	"github.com/wvalencia19/tolling/types"
)

type MiddleWare func(Service) Service

type loggingMiddleWare struct {
	next Service
}

func newLoggingMiddleWare() MiddleWare {
	return func(next Service) Service {
		return loggingMiddleWare{
			next: next,
		}
	}
}

func (mw loggingMiddleWare) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (mw loggingMiddleWare) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}

type instrumentationMiddleWare struct {
	next Service
}

func newInstrumentationMiddleWare() MiddleWare {
	return func(next Service) Service {
		return instrumentationMiddleWare{
			next: next,
		}
	}
}

func (imw instrumentationMiddleWare) Aggregate(_ context.Context, dist types.Distance) error {
	return nil
}

func (imw instrumentationMiddleWare) Calculate(_ context.Context, dist int) (*types.Invoice, error) {
	return nil, nil
}
