package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strconv"
	"time"
)

const DefaultMetricPath = "/metrics"

var httpRequestTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
	Name: "http_request_total",
	Help: "counter: 统计请求数量",
}, []string{"code", "method", "handler", "host", "url"})

var httpRequestDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_request_duration",
	Help: "histogram：统计响应时间",
}, []string{"code", "method", "handler", "url"})

func init() {
	// 注册收集器
	prometheus.MustRegister(httpRequestTotal)
	prometheus.MustRegister(httpRequestDuration)
}

type monitor struct{}

func NewMonitor(e *gin.Engine) *monitor {
	m := &monitor{}
	// 注册metrics路由
	e.GET(DefaultMetricPath, prometheusHandler())
	// 注册中间件
	e.Use(m.HandleFunc())
	return m
}

func (m *monitor) HandleFunc() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// metrics 不统计
		if ctx.Request.URL.Path == DefaultMetricPath {
			ctx.Next()
			return
		}
		start := time.Now()
		ctx.Next()
		status := strconv.Itoa(ctx.Writer.Status())
		httpRequestTotal.WithLabelValues(status, ctx.Request.Method, ctx.HandlerName(), ctx.Request.Host, ctx.Request.RequestURI).Inc()
		httpRequestDuration.WithLabelValues(status, ctx.Request.Method, ctx.HandlerName(), ctx.Request.URL.Path).Observe(time.Since(start).Seconds())
	}
}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()
	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}
