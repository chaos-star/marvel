package Web

import (
	"fmt"
	"github.com/chaos-star/marvel/Log"
	"github.com/gin-gonic/gin"
)

type Web struct {
	*gin.Engine
	port int64
}

func Initialize(port int64, log Log.ILogger, env string) *Web {
	gin.DisableConsoleColor()
	gin.DefaultWriter = log.GetOutput()
	gin.SetMode(gin.DebugMode)
	if env == "prod" || env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}
	if env == "test" {
		gin.SetMode(gin.TestMode)
	}

	return &Web{
		gin.Default(),
		port,
	}
}

func (w *Web) RunServer() {
	w.Run(fmt.Sprintf(":%d", w.port))
}
