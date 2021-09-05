package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"hotel-engine/core/common/logtags"
	"hotel-engine/infrastructure/logger"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func GinBodyLogMiddleware(c *gin.Context) {
	blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
	statusCode := c.Writer.Status()
	if statusCode >= 400 {
		content := blw.body.String()
		logger.WithData(map[string]interface{}{
			"body":       content,
			"statusCode": statusCode,
			"url":        c.Request.RequestURI,
			"method":     c.Request.Method,
		}).WithStatusCode(statusCode).WithUrl(c.Request.RequestURI).
			WithName(logtags.GinRequestFailed).Error("Http request failed")
	}
}
