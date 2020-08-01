package katana

import (
	"net/http/httptest"

	"github.com/gin-gonic/gin"
)

// Context ...
type Context struct {
	*gin.Context
}

type reply struct {
	ErrNo int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

type myTempReader struct{}

func (r myTempReader) Read(p []byte) (n int, err error) {
	return 0, nil
}

// NewContext ...
func NewContext(c *gin.Context) *Context {
	return &Context{
		Context: c,
	}
}

// NewTestContext 生成一个用于单测的context
func NewTestContext() *Context {
	gin.SetMode(gin.TestMode)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Params = []gin.Param{gin.Param{Key: "k", Value: "v"}}
	c.Header("Content-Type", "application/json")
	c.Request = httptest.NewRequest("POST", "/mock/request", myTempReader{})
	return NewContext(c)
}

// Reply ...
func (c *Context) Reply(code int, obj interface{}, err Error) {
	r := &reply{
		ErrNo: err.No(),
		Msg:   err.Error(),
		Data:  obj,
	}
	c.Set("err", err)
	c.JSON(code, r)
}
