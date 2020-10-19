package router

import (
	"github.com/1024casts/fastim/handler/v1/im"
	"github.com/1024casts/fastim/handler/v1/user"
	"github.com/1024casts/snake/router/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

// Load loads the middlewares, routes, handlers.
func Load(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	// 使用中间件
	g.Use(middleware.NoCache)
	g.Use(middleware.Options)
	g.Use(middleware.Secure)
	g.Use(middleware.Logging())
	g.Use(middleware.RequestID())
	g.Use(mw...)

	// 404 Handler.
	//g.NoRoute(handler.RouteNotFound)
	//g.NoMethod(handler.RouteNotFound)

	// swagger api docs
	g.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// pprof router 性能分析路由
	// 默认关闭，开发环境下可以打开
	// 访问方式: HOST/debug/pprof
	// 通过 HOST/debug/pprof/profile 生成profile
	// 查看分析图 go tool pprof -http=:5000 profile
	// see: https://github.com/gin-contrib/pprof
	pprof.Register(g)

	// user
	u := g.Group("/v1")
	u.Use(middleware.AuthMiddleware())
	{
		u.POST("/users/login", user.PhoneLogin)
	}

	// im
	i := g.Group("/v1")
	i.Use(middleware.AuthMiddleware())
	{
		i.POST("/im/send", im.Send)
		i.GET("/im/chat/list", im.ChatList)
		i.POST("/im/msg/list", im.MsgList)
	}

	return g
}
