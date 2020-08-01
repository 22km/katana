package katana

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
)

// RouterGroup ...
type RouterGroup struct {
	routerGroup *gin.RouterGroup
}

// Use ...
func (g *RouterGroup) Use(middleFunc ...MiddlewareFunc) {
	for _, m := range middleFunc {
		g.routerGroup.Use(m.shift())
	}
}

// Group ...
func (g *RouterGroup) Group(relativePath string) *RouterGroup {
	return &RouterGroup{
		g.routerGroup.Group(relativePath),
	}
}

// BasePath ...
func (g *RouterGroup) BasePath() string {
	return g.routerGroup.BasePath()
}

// POST ...
func (g *RouterGroup) POST(relativePath string, f interface{}) {
	g.handle("POST", relativePath, f)
}

// GET ...
func (g *RouterGroup) GET(relativePath string, f interface{}) {
	g.handle("GET", relativePath, f)
}

// DELETE ...
func (g *RouterGroup) DELETE(relativePath string, f interface{}) {
	g.handle("DELETE", relativePath, f)
}

// PATCH ...
func (g *RouterGroup) PATCH(relativePath string, f interface{}) {
	g.handle("PATCH", relativePath, f)
}

// PUT ...
func (g *RouterGroup) PUT(relativePath string, f interface{}) {
	g.handle("PUT", relativePath, f)
}

// OPTIONS ...
func (g *RouterGroup) OPTIONS(relativePath string, f interface{}) {
	g.handle("OPTIONS", relativePath, f)
}

// HEAD ...
func (g *RouterGroup) HEAD(relativePath string, f interface{}) {
	g.handle("HEAD", relativePath, f)
}

// ANY ...
func (g *RouterGroup) ANY(relativePath string, f interface{}) {
	g.handle("ANY", relativePath, f)
}

func checkRouterFunc(f interface{}) error {
	funcType := reflect.TypeOf(f)

	if funcType.Kind() != reflect.Func {
		return fmt.Errorf("type is not function")
	}

	if funcType.NumIn() != 2 {
		return fmt.Errorf("NumIn != 2")
	}

	if funcType.In(0).Kind() != reflect.Ptr {
		return fmt.Errorf("first param is not pointer")
	}

	if funcType.In(0) != reflect.TypeOf(&Context{}) {
		return fmt.Errorf("first param is not valid context")
	}

	if funcType.In(1).Kind() != reflect.Ptr {
		return fmt.Errorf("second param is not pointer")
	}

	if funcType.NumOut() != 2 {
		return fmt.Errorf("NumOut != 2")
	}

	if !reflect.TypeOf(errSuccess).Implements(funcType.Out(1)) {
		return fmt.Errorf("second response is invalid, please implement the ERROR INTERFACE of katana")
	}

	return nil
}

func (g *RouterGroup) handle(method, path string, f interface{}) {
	if err := checkRouterFunc(f); err != nil {
		panic(fmt.Sprintf("func of %s is invalid: %s", g.BasePath()+path, err.Error()))
	}

	handlerFunc := func(ctx *gin.Context) {
		var (
			req  = reflect.New(reflect.TypeOf(f).In(1).Elem()).Interface()
			resp = reflect.New(reflect.TypeOf(f).Out(0).Elem()).Interface()
			c    = NewContext(ctx)
		)

		c.Set("func", f)

		if c == nil || req == nil {
			c.Reply(http.StatusBadRequest, resp, errBadRequest)
			return
		}
		if err := c.Bind(req); err != nil {
			c.Reply(http.StatusBadRequest, resp, errBadRequest.Concat(" | err: %s", err.Error()))
			return
		}

		bytesReq, _ := json.Marshal(req)
		c.Set("req", string(bytesReq))

		retValues := reflect.ValueOf(f).Call([]reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		})

		resp = retValues[0].Interface()
		bytesResp, _ := json.Marshal(resp)
		c.Set("resp", string(bytesResp))

		if retValues[1].IsNil() || retValues[1].Interface().(Error).No() == errSuccess.No() {
			c.Reply(http.StatusOK, resp, errSuccess)
			return
		}

		c.Reply(http.StatusInternalServerError, resp, retValues[1].Interface().(Error))
	}

	switch method {
	case "ANY":
		g.routerGroup.Any(path, handlerFunc)
	default:
		g.routerGroup.Handle(method, path, handlerFunc)
	}
}
