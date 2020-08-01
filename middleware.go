package katana

import (
	"reflect"
	"time"

	"github.com/22km/katana/log"
	"github.com/gin-gonic/gin"
)

// MiddlewareFunc ...
type MiddlewareFunc func(c *Context)

func (m MiddlewareFunc) shift() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := &Context{Context: c}
		reflect.ValueOf(m).Call([]reflect.Value{
			reflect.ValueOf(ctx),
		})
	}
}

// GinRecovery ...
func GinRecovery() MiddlewareFunc {
	return func(c *Context) {
		gin.RecoveryWithWriter(gin.DefaultErrorWriter)
	}
}

// Recorder ...
func Recorder() MiddlewareFunc {
	return func(c *Context) {
		start := time.Now()
		c.Next()
		cost := time.Now().Sub(start)

		// 记录耗时
		bean := log.NewBean("")
		bean.Add("proc_time", cost.Seconds())

		// 记录 handler 信息
		if f, has := c.Get("func"); has {
			bean.SetFuncInfo(log.GetFuncInfo(reflect.ValueOf(f).Pointer()))
		}

		// 记录错误信息
		err, _ := c.Get("err")
		bean.Add("errno", err.(Error).No())
		bean.Add("errmsg", err.(Error).Error())
		bean.Add("msg", err.(Error).Error())

		// 记录 trace 信息
		trace := GetTrace(c.Request)
		bean.Add("traceid", trace.TraceID)
		bean.Add("spanid", trace.SpanID)
		bean.Add("cspanid", "") // access日志没有cspanid

		// 记录调用信息
		bean.Add("uri", c.Request.RequestURI)
		bean.Add("url", "") // access日志不写下游调用

		// 记录 请求 & 返回
		req, _ := c.Get("req")
		resp, _ := c.Get("resp")
		bean.Add("req", req)
		bean.Add("resp", resp)

		if err.(Error).No() == errSuccess.No() {
			bean.SetTitle(dlTagHTTPSuccess)
			bean.Info()
		} else {
			bean.SetTitle(dlTagHTTPFailed)
			bean.Warning()
		}
	}
}

// RawData MiddlewareFunc ...
func CopyRawData() MiddlewareFunc {
	return func(c *Context) {
		body, _ := c.GetRawData()
		c.Set("RowData", body)
		c.Next()
	}
}
