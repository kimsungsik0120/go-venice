package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"time"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		reqID := c.Request.Context().Value("request_id")

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		log.Info().
			Str("request_id", fmt.Sprintf("%v", reqID)).
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status", status).
			Dur("latency", latency).
			Msg("request completed")
	}
}
