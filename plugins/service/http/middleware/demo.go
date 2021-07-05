package middleware

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func Demo() gin.HandlerFunc {
	return func(c *gin.Context) {
		fmt.Println("request:", c.Params)
		fmt.Println("path:", c.FullPath())

		blw := &bodyLogWriter{
			body:           bytes.NewBufferString(""),
			ResponseWriter: c.Writer,
		}

		c.Writer = blw
		c.Next()
		statusCode := c.Writer.Status()
		fmt.Println(statusCode)
		fmt.Println("Response body: " + blw.body.String())
	}
}
