package marvel

import (
	"marvel/Config"
	"marvel/Log"
	"marvel/Mysql/Gorm"
)

var (
	Conf  *Config.Config
	Logger Log.ILogger
	Mysql *Gorm.Engine
)
