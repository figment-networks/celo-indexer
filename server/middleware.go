package server

import (
	"time"

	"github.com/gin-gonic/gin"

	"github.com/figment-networks/celo-indexer/metrics"
	"github.com/figment-networks/celo-indexer/utils/reporting"
)

// setupMiddleware sets up middleware for gin application
func (s *Server) setupMiddleware() {
	s.engine.Use(gin.Recovery())
	s.engine.Use(MetricsMiddleware())
	s.engine.Use(ErrorReportingMiddleware())
}

// MetricsMiddleware is a middleware responsible for logging query execution time
func MetricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()
		c.Next()
		elapsed := time.Since(t)

		metrics.ServerRequestDuration.
			WithLabels(c.Request.URL.Path).
			Observe(elapsed.Seconds())
	}
}

func ErrorReportingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer reporting.RecoverError()
		c.Next()
	}
}
