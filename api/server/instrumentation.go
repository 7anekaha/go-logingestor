package server

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type Instrumentation struct {
	requestLatency    prometheus.Histogram
	requestCounterAdd prometheus.Counter
	requestCounterGet prometheus.Counter
	Service
}

func NewInstrumentation(
	requestLatency prometheus.Summary, 
	requestCounterAdd prometheus.Counter, 
	requestCounterGet prometheus.Counter,
	) ServiceMiddleware {
		return func(next Service) Service {
			return Instrumentation{
				requestLatency:    requestLatency,
				requestCounterAdd: requestCounterAdd,
				requestCounterGet: requestCounterGet,
				Service:           next,
			}
		}
}

func (i Instrumentation) Add(ctx context.Context, log Log) error {
	defer func(begin time.Time) {
		i.requestCounterAdd.Add(1)
		i.requestLatency.Observe(float64(time.Since(begin).Milliseconds()))
	}(time.Now())

	return i.Service.Add(ctx, log)
}

func (i Instrumentation) GetAll(ctx context.Context) ([]Log, error) {
	defer func() {
		i.requestCounterGet.Add(1)
	}()

	return i.Service.GetAll(ctx)
}
