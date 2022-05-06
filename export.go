package marvel

import (
	"PT430/marvel/Config"
	"PT430/marvel/Log"
	"PT430/marvel/Mysql/Gorm"
)

var (
	Conf  *Config.Config
	Logger Log.ILogger
	Mysql *Gorm.Engine
)
