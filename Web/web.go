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

func Initialize(port int64, log Log.ILogger) *Web {
	gin.DisableConsoleColor()
	gin.DefaultWriter = log.GetOutput()
	return &Web{
		gin.Default(),
		port,
	}
}

func (w *Web) RunServer() {
	w.Run(fmt.Sprintf(":%d", w.port))
}
