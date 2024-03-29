package marvel

import (
	"github.com/chaos-star/marvel/Cache"
	"github.com/chaos-star/marvel/CacheCluster"
	"github.com/chaos-star/marvel/Config"
	"github.com/chaos-star/marvel/CronJob"
	Env2 "github.com/chaos-star/marvel/Env"
	etcd "github.com/chaos-star/marvel/Etcd"
	"github.com/chaos-star/marvel/HttpClient"
	"github.com/chaos-star/marvel/Log"
	"github.com/chaos-star/marvel/Mq"
	"github.com/chaos-star/marvel/Mysql/Gorm"
	Srv "github.com/chaos-star/marvel/Server"
	Web2 "github.com/chaos-star/marvel/Web"
)

var (
	Conf         *Config.Config
	Logger       Log.ILogger
	Mysql        *Gorm.Engine
	Etcd         *etcd.Engine
	Env          *Env2.Env
	Server       *Srv.RpcServer
	Web          *Web2.Web
	Cron         *CronJob.CronJob
	MQ           *Mq.Mq
	HttpTool     *HttpClient.HttpClient
	Redis        *Cache.Cache
	RedisCluster *CacheCluster.Cluster
)
