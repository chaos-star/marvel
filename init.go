package marvel

import (
	"github.com/khaos-star/marvel/Config"
	"github.com/khaos-star/marvel/Log"
	"github.com/khaos-star/marvel/Mysql/Gorm"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"time"
)

func Initialize() error {
	var err error
	err, Conf = Config.Initialize()
	if err != nil {
		return err
	}
	//如果配置存在则初始化日志
	if Conf.IsSet("log") {
		LogConf := Conf.GetStringMap("log")
		path := "./log"
		name := "%Y%m%d.log"
		if cPath, ok := LogConf["path"]; ok {
			path = cPath.(string)
		}
		if cName, ok := LogConf["name"]; ok {
			name = cName.(string)
		}
		var options []rotatelogs.Option

		if cLinkName, ok := LogConf["link_name"]; ok {
			if cLinkName != "" {
				options = append(options, rotatelogs.WithLinkName(cLinkName.(string)))
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
			return err
		}
	}

	//如果配置存在则初始化Mysql
	if Conf.IsSet("mysql"){
		MysqlConf := Conf.GetConfig("mysql")
		Mysql,err = Gorm.Initialize(MysqlConf)
		if err != nil{
			return err
		}
	}

	return nil
}
