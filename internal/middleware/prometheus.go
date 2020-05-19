package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var defaultMetricPath = "/metrics"

var requestTotal = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "request_total",
	Help: "How many HTTP requests processed, partitioned by status code and HTTP method",
})

type monitor struct{}

func NewMonitor(ctx *gin.Engine) *monitor {
	m := &monitor{}
	m.Use(ctx)
	ctx.Use(m.HandleFunc())
}

func (m *monitor) Use(ctx *gin.Engine) {

}

func (m *monitor) HandleFunc() gin.HandlerFunc {
	return func(context *gin.Context) {
		if context.Request.URL.Path == defaultMetricPath {
			context.Next()
			return
		}
		requestTotal.Inc()
		context.Next()
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
