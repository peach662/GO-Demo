package trace

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ginContextTraceIDKey is used with gin.Context Set/Get (string keys).
const ginContextTraceIDKey = "traceID"

// contextKey avoids colliding with other packages' context values (SA1029).
type contextKey string

const traceIDKey contextKey = "traceID"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}

		ctx := context.WithValue(c.Request.Context(), traceIDKey, traceID)
		c.Request = c.Request.WithContext(ctx)
		c.Set(ginContextTraceIDKey, traceID)
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}

func IDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(traceIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
