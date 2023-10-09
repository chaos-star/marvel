package marvel

import (
	"github.com/chaos-star/marvel/Config"
	Env2 "github.com/chaos-star/marvel/Env"
	etcd "github.com/chaos-star/marvel/Etcd"
	"github.com/chaos-star/marvel/Log"
	"github.com/chaos-star/marvel/Mysql/Gorm"
	Srv "github.com/chaos-star/marvel/Server"
)

var (
	Conf   *Config.Config
	Logger Log.ILogger
	Mysql  *Gorm.Engine
	Etcd   *etcd.Engine
	Env    *Env2.Env
	Server *Srv.RpcServer
)
