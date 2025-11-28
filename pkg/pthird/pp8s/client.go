package pp8s

import (
	"context"
	"fmt"
	"time"

	"github.com/VictoriaMetrics/metricsql"
	"github.com/prometheus/client_golang/api"
	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/prometheus/common/model"
)

// Client is the interface for prometheus client
type Client interface {
	QueryVector(query string, top int64) ([]*model.Sample, int64, error)
}

type Option struct {
	// The address of the Prometheus to connect to.
	Address string
}

type client struct {
	api v1.API
}

func NewClient(opt Option) (Client, error) {
	cli, err := api.NewClient(api.Config{
		Address: opt.Address,
	})
	if err != nil {
		return nil, err
	}
	return &client{
		api: v1.NewAPI(cli),
	}, nil
}

func NewDefault(addr string) (Client, error) {
	return NewClient(Option{Address: addr})
}

func (c *client) checkPromql(query string) error {
	_, err := metricsql.Parse(query)
	return err
}

func (c *client) QueryVector(query string, top int64) ([]*model.Sample, int64, error) {
	var count int64
	samples := make([]*model.Sample, 0)

	err := c.checkPromql(query)
	if err != nil {
		return samples, count, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	queryCount := fmt.Sprintf("count(%s)", query)
	result, _, err := c.api.Query(ctx, queryCount, time.Now())
	if err != nil {
		return samples, count, err
	}
	if result.Type() == model.ValVector {
		for _, v := range result.(model.Vector) {
			count = int64(v.Value)
			break
		}
	}
	if top > 0 {
		query = fmt.Sprintf("topk(%d, %s)", top, query)
	}
	result, _, err = c.api.Query(ctx, query, time.Now())
	if err != nil {
		return samples, count, err
	}
	if result.Type() == model.ValVector {
		for _, v := range result.(model.Vector) {
			samples = append(samples, v)
		}
	}
	return samples, count, nil
}
