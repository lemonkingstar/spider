package pmetric

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lemonkingstar/spider/pkg/iserver"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// a global register which is used to collect metrics we need.
// it will be initialized when process is up for safe usage.
// and then be revised later when service is initialized.
var globalRegister prometheus.Registerer

func init() {
	// set default global register
	globalRegister = prometheus.DefaultRegisterer
}

// Register must only be called after backbone engine is started.
func Register() prometheus.Registerer {
	return globalRegister
}

func MustRegister(cs ...prometheus.Collector) {
	globalRegister.MustRegister(cs...)
}

const (
	Namespace = "px"

	LabelProcessName = "process_name"
	LabelHost        = "host"
	LabelOrigin      = "origin"
	LabelRemote      = "remote"
	LabelHandler     = "handler"
	LabelHTTPStatus  = "status_code"
)

type Service struct {
	httpHandler http.Handler

	registry        prometheus.Registerer
	requestTotal    *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

// NewService returns new metrics service
func NewService() *Service {
	registry := prometheus.NewRegistry()
	register := prometheus.WrapRegistererWith(prometheus.Labels{LabelProcessName: iserver.GetAppName(),
		LabelHost: ""}, registry)

	// set up global register
	globalRegister = register

	srv := Service{registry: register}

	srv.requestTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: Namespace + "_http_request_total",
			Help: "http request total.",
		},
		[]string{LabelHandler, LabelHTTPStatus},
	)
	register.MustRegister(srv.requestTotal)

	srv.requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    Namespace + "_http_request_duration_millisecond",
			Help:    "Histogram of latencies for HTTP requests.",
			Buckets: []float64{10, 30, 50, 70, 100, 200, 300, 400, 500, 1000, 2000, 5000},
		},
		[]string{LabelHandler},
	)
	register.MustRegister(srv.requestDuration)
	register.MustRegister(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	register.MustRegister(prometheus.NewGoCollector())

	srv.httpHandler = promhttp.InstrumentMetricHandler(
		registry, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
	)
	return &srv
}

// Registry returns the prometheus.Registerer
func (s *Service) Registry() prometheus.Registerer {
	return s.registry
}

func (s *Service) MustRegister(cs ...prometheus.Collector) {
	s.registry.MustRegister(cs...)
}

func (s *Service) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	s.httpHandler.ServeHTTP(resp, req)
}

func (s *Service) Handler(c *gin.Context) {
	s.httpHandler.ServeHTTP(c.Writer, c.Request)
}
