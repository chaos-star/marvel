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

func Initialize(port int64, log Log.ILogger, env string, trusted []string) *Web {
	gin.DisableConsoleColor()
	useEnv := gin.DebugMode
	gin.SetMode(gin.DebugMode)
	if env == "prod" || env == "production" {
		useEnv = gin.ReleaseMode
	}
	if env == "test" {
		useEnv = gin.TestMode
	}

	gin.SetMode(useEnv)
	switch useEnv {
	case gin.ReleaseMode:
		gin.DefaultWriter = log.GetOutput()
	case gin.TestMode:
		gin.DefaultWriter = log.GetOutput()
	case gin.DebugMode:

	}
	router := gin.Default()
	if len(trusted) > 0 {
		router.SetTrustedProxies(trusted)
	}

	return &Web{
		router,
		port,
	}
}

func (w *Web) RunServer() {
	addr := fmt.Sprintf(":%d", w.port)
	fmt.Println(fmt.Sprintf("%s Running...", addr))
	w.Run(addr)
}
