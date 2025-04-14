package logs

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
	"io"
	"time"
)

type MiddllewareBuilder struct {
	showRequest  bool
	showResponse atomic.Bool
}

type responseBodyWriter struct {
	gin.ResponseWriter
	al *AccessLog
}

type RequestLog struct {
	Method      string
	Path        string
	RequestBody string
}

type ResponseLog struct {
	StatusCode   int
	ResponseBody string
}

type AccessLog struct {
	RequestLog
	ResponseLog
	Cost string
}

func (r responseBodyWriter) Write(b []byte) (int, error) {
	r.al.ResponseBody = string(b)
	return r.ResponseWriter.Write(b)
}

func (c responseBodyWriter) WriteString(s string) (int, error) {
	c.al.ResponseBody = s
	return c.ResponseWriter.WriteString(s)
}

func (c responseBodyWriter) WriteHeader(code int) {
	c.al.StatusCode = code
	c.ResponseWriter.WriteHeader(code)
}

func (b *MiddllewareBuilder) EnableRequest() *MiddllewareBuilder {
	b.showRequest = true
	return b
}

func (b *MiddllewareBuilder) EnableResponse() *MiddllewareBuilder {
	b.showResponse.Store(true)
	return b
}

func NewMiddlewareBuilder() *MiddllewareBuilder {
	return &MiddllewareBuilder{
		showRequest:  false,
		showResponse: atomic.Bool{},
	}
}

func (b *MiddllewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		now := time.Now()
		url := ctx.Request.URL.String()
		if len(url) > 1024 {
			url = url[:1024]
		}
		method := ctx.Request.Method
		al := AccessLog{
			RequestLog: RequestLog{
				Method: method,
				Path:   url,
			},
		}
		requestBody := ctx.Request.Body
		if requestBody != nil && b.showRequest {
			//body 读出来就没了
			body, _ := ctx.GetRawData()
			if len(body) > 1024 {
				body = body[:1024]
			}
			req := io.NopCloser(bytes.NewBuffer(body))
			ctx.Request.Body = req
			al.RequestBody = string(body)
		}

		if b.showResponse.Load() {
			// 获取响应体
			ctx.Writer = responseBodyWriter{
				ResponseWriter: ctx.Writer,
				al:             &al,
			}
		}
		defer func() {
			al.Cost = time.Since(now).String()
			logrus.WithField("Request: ", al.RequestLog).WithField("Response: ", al.ResponseLog).
				WithField("Duration: ", al.Cost).Info("access log")
		}()
		ctx.Next()
	}
}
