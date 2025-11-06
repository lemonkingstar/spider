package pmetrics

import "github.com/prometheus/client_golang/prometheus"

type Custom struct {
	ds []*prometheus.Desc
	fs []func(chan<- prometheus.Metric)
}

func NewCustom() *Custom {
	c := &Custom{}
	MustRegister(c)
	return c
}

func (c *Custom) Describe(ch chan<- *prometheus.Desc) {
	for _, v := range c.ds {
		ch <- v
	}
}

func (c *Custom) Collect(ch chan<- prometheus.Metric) {
	for _, v := range c.fs {
		v(ch)
	}
}

func (c *Custom) Add(d *prometheus.Desc, f func(chan<- prometheus.Metric)) {
	c.ds = append(c.ds, d)
	c.fs = append(c.fs, f)
}
