package http

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/trace"
)

func LogMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		logger := gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			var statusColor, methodColor, resetColor string
			if param.IsOutputColor() {
				statusColor = param.StatusCodeColor()
				methodColor = param.MethodColor()
				resetColor = param.ResetColor()
			}

			if param.Latency > time.Minute {
				param.Latency = param.Latency.Truncate(time.Second)
			}
			span := trace.SpanFromContext(c.Request.Context())
			spanContext := span.SpanContext()
			traceId := spanContext.TraceID()
			spanId := spanContext.SpanID()

			log := fmt.Sprintf("[GIN] %v |%s %3d %s | %13v | %15s |%s %-7s %s %#v",
				param.TimeStamp.Format("2006/01/02 - 15:04:05"),
				statusColor, param.StatusCode, resetColor,
				param.Latency,
				param.ClientIP,
				methodColor, param.Method, resetColor,
				param.Path,
			)
			if spanContext.IsValid() {
				log = fmt.Sprintf("%s | trace_id=%s, span_id=%s", log, traceId.String(), spanId.String())
			}

			log = fmt.Sprintf("%s\n%s", log, param.ErrorMessage)

			return log
		})
		logger(c)
	}
}
