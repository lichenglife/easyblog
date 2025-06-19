package monitor

import (
	"context"
	"fmt"
	"net/http/pprof"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

// Monitor 监控结构体
type Monitor struct {
	// 服务名称
	serviceName string
	// 服务地址
	serviceAddr string
	// 是否启用pprof
	enablePprof bool
	// 是否启用trace
	enableTrace bool
	// 是否启用metrics
	enableMetrics bool
	// metrics注册表
	registry *prometheus.Registry
	// 自定义metrics
	httpRequestsTotal   *prometheus.CounterVec
	httpRequestDuration *prometheus.HistogramVec
}

// NewMonitor 创建监控实例
func NewMonitor(serviceName, serviceAddr string, enablePprof, enableTrace, enableMetrics bool) *Monitor {
	return &Monitor{
		serviceName:   serviceName,
		serviceAddr:   serviceAddr,
		enablePprof:   enablePprof,
		enableTrace:   enableTrace,
		enableMetrics: enableMetrics,
		registry:      prometheus.NewRegistry(),
	}
}

// Init 初始化监控
func (m *Monitor) Init() error {
	// 初始化metrics
	if m.enableMetrics {
		if err := m.initMetrics(); err != nil {
			return fmt.Errorf("初始化metrics失败: %v", err)
		}
	}

	// 初始化trace
	if m.enableTrace {
		if err := m.initTrace(); err != nil {
			return fmt.Errorf("初始化trace失败: %v", err)
		}
	}

	return nil
}

// RegisterRoutes 注册监控路由
func (m *Monitor) RegisterRoutes(r *gin.Engine) {
	// 注册pprof路由
	if m.enablePprof {
		// 创建pprof路由组
		pprofGroup := r.Group("/debug/pprof")
		{
			pprofGroup.GET("/", gin.WrapF(pprof.Index))
			pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
			pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
			pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
			pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
			pprofGroup.GET("/allocs", gin.WrapF(pprof.Handler("allocs").ServeHTTP))
			pprofGroup.GET("/block", gin.WrapF(pprof.Handler("block").ServeHTTP))
			pprofGroup.GET("/goroutine", gin.WrapF(pprof.Handler("goroutine").ServeHTTP))
			pprofGroup.GET("/heap", gin.WrapF(pprof.Handler("heap").ServeHTTP))
			pprofGroup.GET("/mutex", gin.WrapF(pprof.Handler("mutex").ServeHTTP))
			pprofGroup.GET("/threadcreate", gin.WrapF(pprof.Handler("threadcreate").ServeHTTP))
		}
	}

	// 注册metrics路由
	if m.enableMetrics {
		r.GET("/metrics", gin.WrapH(promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{})))
	}
}

// initMetrics 初始化metrics
func (m *Monitor) initMetrics() error {
	// 注册默认的metrics
	m.registry.MustRegister(
		prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}),
		prometheus.NewGoCollector(),
	)

	// 创建自定义metrics
	m.httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	m.httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// 注册自定义metrics
	m.registry.MustRegister(
		m.httpRequestsTotal,
		m.httpRequestDuration,
	)

	return nil
}

// initTrace 初始化trace
func (m *Monitor) initTrace() error {
	// 创建Jaeger导出器
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
	if err != nil {
		return err
	}

	// 创建资源
	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(m.serviceName),
		),
	)
	if err != nil {
		return err
	}

	// 创建TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)

	// 设置全局TracerProvider
	otel.SetTracerProvider(tp)

	return nil
}

// Middleware 监控中间件
func (m *Monitor) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 处理请求
		c.Next()

		// 记录metrics
		if m.enableMetrics {
			duration := time.Since(start).Seconds()
			status := fmt.Sprintf("%d", c.Writer.Status())
			method := c.Request.Method
			path := c.FullPath()

			// 记录请求计数
			m.httpRequestsTotal.WithLabelValues(method, path, status).Inc()

			// 记录请求延迟
			m.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
		}
	}
}
