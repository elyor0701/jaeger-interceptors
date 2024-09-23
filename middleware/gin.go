package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"jaeger-interceptors/config"
	"net/http"
)

func TracingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tracer := otel.Tracer(config.TracerNameHTTPMiddleware)
		ctx, span := tracer.Start(c.Request.Context(), fmt.Sprintf("[%s] %s", c.Request.Method, c.Request.RequestURI))
		defer span.End()

		c.Request = c.Request.WithContext(ctx)

		c.Next()

		statusCode := c.Writer.Status()
		span.SetAttributes(attribute.Int("http.status_code", statusCode))

		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				span.RecordError(ginErr.Err)
				span.SetStatus(codes.Error, ginErr.Error())
			}
		} else if statusCode >= 400 {
			span.SetStatus(codes.Error, http.StatusText(statusCode))
		} else {
			span.SetStatus(codes.Ok, "Request processed successfully")
		}
	}
}
