package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"short-code/global"
	"short-code/initialize"
	"short-code/middleware"
	"short-code/service/choreography"
	"strconv"
	"time"
)

type Server interface {
	StartUp()
}

type DefaultServer struct {
	router *gin.Engine
}

func (e *DefaultServer) StartUp() {
	e.Init()
	e.Run()
}
func (e *DefaultServer) Init() {
	initialize.ConfByEnv(os.Args[1])
	initialize.Logrus()

	initialize.RedisClient()
	initialize.LocalCache()
	initialize.GRomClickHouse()
	initialize.Limiter()
	initialize.Tasks()
	//initialize end

	gin.SetMode(global.CONF.Server.GinMode)
	gin.DefaultWriter = io.Discard
	router := gin.Default()

	router.Use(middleware.LoggerMiddleware())
	initialize.PromHandler(router)
	e.router = router
	e.pprofStart()
}

func (e *DefaultServer) pprofStart() {
	e.router.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			log.Fatalf("pprof server failed: %v", err)
		}
	}()
}

func (e *DefaultServer) Run() {
	port := strconv.FormatUint(global.CONF.Server.Port, 10)
	srv := &http.Server{Addr: ":" + port, Handler: e.router}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			global.LOG.Fatalf("listen: %s\n", err)
		}
	}()
	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	global.LOG.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		global.LOG.Fatal("Server Shutdown:", err)
	}
	e.destroy()
	global.LOG.Info("Server exiting")
}

func (e *DefaultServer) destroy() {
	var bindDataService = choreography.BindDataService{}
	bindDataService.Flush()
}
