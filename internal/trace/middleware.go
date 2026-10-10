package trace 

import(
	"context"
	"github.com/google/uuid"
	"github.com/gin-gonic/gin"
)
const ContextTraceIDKey = "traceID"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceID := c.GetHeader("X-Request-ID")
		if traceID == "" {
			traceID = uuid.New().String()
		}
		
	ctx := context.WithValue(c.Request.Context(), ContextTraceIDKey, traceID)
	c.Request = c.Request.WithContext(ctx)
		c.Set(ContextTraceIDKey, traceID)
		c.Header("X-Request-ID", traceID)
		c.Next()
	}
}

func IDFromContext(ctx context.Context) string{
	id, ok := ctx.Value(ContextTraceIDKey).(string)
	if !ok {
		return ""
	}
	return id
}
