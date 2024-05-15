package api

import (
	_ "xxx-server/docs"
	"xxx-server/infrastructure/config"
	"xxx-server/interface/handler/http_handler"
	"xxx-server/interface/handler/middleware"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func newHTTPServer() *gin.Engine {
	debug := config.C.Server.Debug
	if debug {
		gin.ForceConsoleColor()
		gin.SetMode(gin.DebugMode)
	} else {
		gin.DisableConsoleColor()
		gin.SetMode(gin.ReleaseMode)
	}

	pathPrefix := config.APP_PREFIX + "/v1"

	router := gin.New()
	router.MaxMultipartMemory = 128 << 20 // 128 MiB

	router.Use(gin.LoggerWithConfig(gin.LoggerConfig{
		SkipPaths: []string{
			pathPrefix + "/xxx/:id/mid",
		},
	}), gin.Recovery(), middleware.CorsHandler())

	v1 := router.Group(pathPrefix)
	v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 后台任务模块
	th := http_handler.NewTaskHandler()
	gbh := v1.Group("/task")
	gbh.DELETE("", th.Deletes)
	gbh.GET("/:id", th.QueryDetail)
	gbh.GET("", th.Query)
	gbh.PUT("/:id", th.Restart)
	gbh.POST("/:id/download", th.Download)
	gbh.GET("/:id/zip", th.GetZip)

	// XXX数据模块
	xh := http_handler.NewXxxDataHandler()
	gxh := v1.Group("/xxx")
	gxh.POST("", xh.Create)
	gxh.DELETE("", xh.Deletes)
	gxh.GET("", xh.QueryList)
	gxh.POST("/:id", xh.QueryDetail)

	// 设置模块
	sh := http_handler.NewSettingHandler()
	gsh := v1.Group("/setting")
	gsh.GET("/:id", sh.QueryDetail)
	gsh.GET("/alert/:adcode", sh.QueryAlertForAdcode)
	gsh.GET("", sh.Query)
	gsh.POST("/:section", sh.Create)
	gsh.DELETE("", sh.Deletes)
	gsh.POST("/shapefile", sh.ShpToWkt)

	if debug {
		if config.C.Server.RunLocal {
			pprof.Register(router)
		} else {
			pprof.Register(router, config.APP_PREFIX+pprof.DefaultPrefix)
		}
	}
	return router
}

func Run() error {
	httpServer := newHTTPServer()
	return httpServer.Run(config.C.Server.Addr)
}
