package router

import (
	"net/http"

	"github.com/ChenHaoJie9527/Elk-Mall/adaptor"
	"github.com/ChenHaoJie9527/Elk-Mall/config"
	"github.com/gin-gonic/gin"
)

type IRouter interface {
	Register(app *gin.Engine)
	SpanFilter(r *gin.Context) bool
	AccessRecordFilter(r *gin.Context) bool
}

type Router struct {
	FullPPROF bool           // 是否开启 PPROF
	rootPath  string         // 根路径
	conf      *config.Config // 配置
	checkFunc func() error   // 检查函数
	// admin     *admin.Ctrl
	// customer  *customer.Ctrl
}

func NewRouter(conf *config.Config, adaptor *adaptor.Adaptor, checkFunc func() error) *Router {
	return &Router{
		FullPPROF: conf.Server.EnablePprof,
		rootPath:  "/api/mall",
		conf:      conf,
		checkFunc: checkFunc,
	}
}

// 注册业务路由
func (r *Router) Register(app *gin.Engine) {
	// TODO: pprof 路由
	if r.conf.Server.EnablePprof {
		SetupPprof(app, "/debug/pprof")
	}

	app.Any("/ping", r.checkServer())

	// 创建根路由组
	root := createGroup(r.rootPath, app)

	// 注册路由
	r.route(root)
}

// 跨度过滤器
func (r *Router) SpanFilter(ctx *gin.Context) bool {
	return true
}

// 检查服务
func (r *Router) checkServer() func(*gin.Context) {
	return func(ctx *gin.Context) {
		err := r.checkFunc()
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{})
	}
}

// 创建路由组
func createGroup(rootPath string, app *gin.Engine) *gin.RouterGroup {
	return app.Group(rootPath)
}

// 注册路由
func (r *Router) route(root *gin.RouterGroup) {
	r.customerRoute(root)
	r.adminRoute(root)
}

// 注册客户路由
func (r *Router) customerRoute(root *gin.RouterGroup) {}

// 注册管理员路由
func (r *Router) adminRoute(root *gin.RouterGroup) {}
