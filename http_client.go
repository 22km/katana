package katana

import (
	"encoding/json"
	"runtime"

	"github.com/22km/katana/log"
	"github.com/go-resty/resty/v2"
)

// Client 封装 resty client
type Client struct {
	*resty.Client
}

// NewHTTPClient ...
// resty: https://github.com/go-resty/resty
func NewHTTPClient(ctx *Context) *Client {
	client := resty.New()
	bean := log.NewBean("")

	// 设置函数信息
	pc, _, _, ok := runtime.Caller(1)
	if ok {
		bean.SetFuncInfo(log.GetFuncInfo(pc))
	} else {
		bean.SetFuncInfo("unknow")
	}

	// 请求前, 注入trace
	client.OnBeforeRequest(func(c *resty.Client, req *resty.Request) error {
		// 获取 trace
		trace := GetTrace(ctx.Request)
		cspanid := genSpanID() // 为下游生成 cspanid

		// 记录 trace 信息
		bean.Add("traceid", trace.TraceID)
		bean.Add("spanid", trace.SpanID)
		bean.Add("cspanid", cspanid)

		// 注入 trace
		headers := trace.ToMap()
		headers[ktnHeaderCspanID] = cspanid // 注入前需要替换 spanid
		req.SetHeaders(headers)

		return nil
	})

	// 返回后, 打印log
	client.OnAfterResponse(func(c *resty.Client, resp *resty.Response) error {
		// log
		bean.Add("proc_time", resp.Time().Seconds())
		bean.Add("errno", resp.StatusCode())
		bean.Add("uri", ctx.Request.RequestURI) // 当前接口
		bean.Add("url", resp.Request.URL)       // 下游接口

		// request body
		_req, _ := json.Marshal(resp.Request.Body)
		bean.Add("req", string(_req))

		// request form
		_form, _ := json.Marshal(resp.Request.FormData)
		bean.Add("form", string(_form))

		// response body
		_resp := resp.Body()
		bean.Add("resp", string(_resp))

		if resp.IsSuccess() {
			bean.SetTitle(dlTagHTTPSuccess)
			bean.Info()
		} else {
			bean.SetTitle(dlTagHTTPFailed)
			bean.Warning()
		}

		return nil
	})

	return &Client{client}
}
