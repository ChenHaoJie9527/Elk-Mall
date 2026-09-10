package router

import (
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

func SetupPprof(app *gin.Engine, prefix string) {
	group := createGroup(prefix, app)
	group.GET("/", gin.WrapF(pprof.Index))
	group.GET("/cmdline", gin.WrapF(pprof.Cmdline))
	group.GET("/profile", gin.WrapF(pprof.Profile))
	group.GET("/symbol", gin.WrapF(pprof.Symbol))
	group.GET("/trace", gin.WrapF(pprof.Trace))
	group.GET("/heap", pprofHandler("heap"))
	group.GET("/goroutine", pprofHandler("goroutine"))
	group.GET("/block", pprofHandler("block"))
	group.GET("/mutex", pprofHandler("mutex"))
}

// pprof 处理器
func pprofHandler(name string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		pprof.Handler(name).ServeHTTP(ctx.Writer, ctx.Request)
	}
}
