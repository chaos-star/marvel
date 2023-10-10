package marvel

import (
	"fmt"
	"github.com/chaos-star/marvel/Config"
	etcd "github.com/chaos-star/marvel/Etcd"
	"github.com/chaos-star/marvel/Log"
	"github.com/chaos-star/marvel/Mysql/Gorm"
	srv "github.com/chaos-star/marvel/Server"
	Web2 "github.com/chaos-star/marvel/Web"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"time"
)

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
	var sysConf map[string]interface{}
	//检测系统日志初始化环境变量
	if Conf.IsSet("system") {
		sysConf = Conf.GetStringMap("system")
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
		err, Logger = Log.Initialize(path, name, options...)
		if err != nil {
			return
		}
	}

	//如果配置存在则初始化Mysql
	if Conf.IsSet("mysql") {
		MysqlConf := Conf.GetConfig("mysql")
		Mysql, err = Gorm.Initialize(MysqlConf)
		if err != nil {
			return
		}
	}

	//如果配置存在则初始化Etcd

	if Conf.IsSet("etcd") {
		EtcdConf := Conf.GetStringMap("etcd")
		Etcd, err = etcd.Initialize(EtcdConf)
		if err != nil {
			return
		}
	}

	if Etcd != nil {
		if _, ok := sysConf["port"]; !ok {
			panic("loss rpc port")
		}
		if int(sysConf["port"].(int64)) > 0 {
			Server = srv.Initialize(Etcd, int(sysConf["port"].(int64)), sysConf["prefix"].(string))
		}
	}

	if HttpPort, ok := sysConf["http_port"]; ok && int(HttpPort.(int64)) > 0 {
		var trustedProxies = Conf.GetStringSlice("system.http_trusted_proxy")
		Web = Web2.Initialize(HttpPort.(int64), Logger, sysConf["env"].(string), trustedProxies)
	}

	return
}
