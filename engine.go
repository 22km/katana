package katana

import (
	"github.com/gin-gonic/gin"
)

// Engine ...
type Engine struct {
	engine *gin.Engine
}

// New build engine
func New(mode string) *Engine {
	gin.SetMode(mode)
	e := &Engine{
		engine: gin.New(),
	}

	return e
}

// Run ...
func (e *Engine) Run(port string) error {
	return e.engine.Run(":" + port)
}

// Group ...
func (e *Engine) Group(relativePath string) *RouterGroup {
	return &RouterGroup{
		e.engine.Group(relativePath),
	}
}

// Use ...
func (e *Engine) Use(middleFunc ...MiddlewareFunc) {
	for _, m := range middleFunc {
		e.engine.Use(m.shift())
	}
}
