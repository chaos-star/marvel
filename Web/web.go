package Web

import (
	"fmt"
	"github.com/chaos-star/marvel/Log"
	"github.com/gin-gonic/gin"
)

type Web struct {
	*gin.Engine
	port     int64
	safeCidr []string
}

func Initialize(port int64, log Log.ILogger, env string, trusted []string) *Web {
	gin.DisableConsoleColor()
	useEnv := gin.DebugMode
	if env == "prod" || env == "production" {
		useEnv = gin.ReleaseMode
	}
	if env == "test" {
		useEnv = gin.TestMode
	}
	gin.SetMode(useEnv)
	router := gin.Default()
	if len(trusted) > 0 {
		router.SetTrustedProxies(trusted)
	}

	return &Web{
		router,
		port,
		trusted,
	}
}

type IGroupRouter interface {
	Router(group *gin.RouterGroup)
}

func (w *Web) LoadRouterGroup(prefix string, routers []IGroupRouter, handlers ...gin.HandlerFunc) {
	group := w.Group(prefix, handlers...)
	for _, router := range routers {
		router.Router(group)
	}
}

func (w *Web) RunServer() {
	addr := fmt.Sprintf(":%d", w.port)
	fmt.Println(fmt.Sprintf("Trusted CIDR:%#v Port: %s Running...", w.safeCidr, addr))
	w.Run(addr)
}
