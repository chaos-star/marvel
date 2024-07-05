package marvel

import (
	"fmt"
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
	srv "github.com/chaos-star/marvel/Server"
	Web2 "github.com/chaos-star/marvel/Web"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"os"
	"time"
)

const EnvName = "ARTEMIS_APP_ENV"

func init() {
	var err error
	defer func() {
		if err != nil {
			fmt.Println(fmt.Sprintf("Init Exception:%#v", err))
		}
	}()
	err, Conf = Config.Initialize()
	if err != nil {
		panic(err)
		return
	}
	var (
		sysConf map[string]interface{}
		OsEnv   string
	)
	OsEnv = os.Getenv(EnvName)
	//检测系统日志初始化环境变量
	if Conf.IsSet("system") {
		sysConf = Conf.GetStringMap("system")
		var (
			confEnv interface{}
			env     string
			ok      bool
		)
		if OsEnv == "" {
			if confEnv, ok = sysConf["env"]; ok {
				if env, ok = confEnv.(string); ok {
					OsEnv = env
				}
			}
		}

		Env = Env2.Initialize(OsEnv)

	} else {
		panic(" system configuration loss")
	}
	//如果配置存在则初始化日志
	if Conf.IsSet("log") {
		LogConf := Conf.GetStringMap("log")
		path := "./logs"
		name := "%Y-%m-%d.log"
		if cPath, ok := LogConf["path"]; ok {
			path = cPath.(string)
		}
		if cName, ok := LogConf["name"]; ok {
			name = cName.(string)
		}
		var options []rotatelogs.Option

		if cLinkName, ok := LogConf["link_name"]; ok {
			if cLinkName != "" {
				options = append(options, rotatelogs.WithLinkName(fmt.Sprintf("./%s", cLinkName.(string))))
			}
		}

		if cMaxAge, ok := LogConf["max_age"]; ok {
			if cMaxAge.(int64) > 0 {
				options = append(options, rotatelogs.WithMaxAge(time.Hour*24*time.Duration(cMaxAge.(int64))))
			}
		}

		if cRotationCount, ok := LogConf["rotation_count"]; ok {
			if cRotationCount.(uint) > 0 {
				options = append(options, rotatelogs.WithRotationCount(cRotationCount.(uint)))
			}
		}

		if cRotationSize, ok := LogConf["rotation_size"]; ok {
			if cRotationSize.(int64) > 0 {
				options = append(options, rotatelogs.WithRotationSize(cRotationSize.(int64)*1024*1024))
			}
		}
		err, Logger = Log.Initialize(OsEnv, path, name, options...)
		if err != nil {
			return
		}
		fmt.Println("Logger Initialize [\033[32mSuccess\033[0m]")
	}

	HttpTool = HttpClient.Initialize(Logger)
	//如果配置存在则初始化Mysql
	if Conf.IsSet("mysql") {
		MysqlConf := Conf.Get("mysql")
		Mysql, err = Gorm.Initialize(OsEnv, MysqlConf, Logger)
		if err != nil {
			return
		}
		fmt.Println("Mysql Initialize [\033[32mSuccess\033[0m]")
	}

	//如果配置存在则初始化Redis
	if Conf.IsSet("redis") {
		RedisConf := Conf.Get("redis")
		Redis, err = Cache.Initialize(RedisConf)
		if err != nil {
			panic(err)
		}
		fmt.Println("Redis Initialize [\033[32mSuccess\033[0m]")
	}

	//如果配置存在则初始化Redis
	if Conf.IsSet("redis_cluster") {
		RedisConf := Conf.Get("redis_cluster")
		RedisCluster, err = CacheCluster.Initialize(RedisConf)
		if err != nil {
			panic(err)
		}
		fmt.Println("Redis Cluster Initialize [\033[32mSuccess\033[0m]")
	}

	//如果配置存在则初始化Rabbitmq
	if Conf.IsSet("rabbitmq") {
		MqConf := Conf.Get("rabbitmq")
		MQ, err = Mq.Initialize(MqConf, Logger)
		if err != nil {
			panic(err)
		}
		fmt.Println("RabbitMQ Initialize [\033[32mSuccess\033[0m]")
	}

	//如果配置存在则初始化Etcd

	if Conf.IsSet("etcd") {
		EtcdConf := Conf.GetStringMap("etcd")
		Etcd, err = etcd.Initialize(EtcdConf)
		if err != nil {
			panic(err)
		}
		fmt.Println("Etcd Initialize [\033[32mSuccess\033[0m]")
	}

	if Etcd != nil {
		if _, ok := sysConf["rpc_port"]; !ok {
			panic("loss rpc port")
		}
		if int(sysConf["rpc_port"].(int64)) > 0 {
			Server = srv.Initialize(Etcd, int(sysConf["rpc_port"].(int64)), sysConf["prefix"].(string))
			fmt.Println("Rpc Server Initialize [\033[32mSuccess\033[0m]")
		}
	}

	if HttpPort, ok := sysConf["http_port"]; ok && int(HttpPort.(int64)) > 0 {
		var trustedProxies = Conf.GetStringSlice("system.http_trusted_proxy")
		_env := Env.GetEnv()
		Web = Web2.Initialize(HttpPort.(int64), Logger, _env, trustedProxies)
		fmt.Println("Web Server Initialize [\033[32mSuccess\033[0m]")
	}

	Cron = CronJob.Initialize()

	return
}
