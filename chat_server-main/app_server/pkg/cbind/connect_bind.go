package cbind

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Binder struct {
	g *gin.RouterGroup
}

func NewBinder(g *gin.RouterGroup) *Binder {
	return &Binder{g: g}
}

func (b *Binder) Bind(path string, handler http.Handler) {
	b.g.POST(path+"*action", func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	})
}
